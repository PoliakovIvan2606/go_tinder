package background

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"tinder/internal/app/store"

	"github.com/redis/go-redis/v9"
)

type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc
	st     *store.Store
	redis  *redis.Client
	done  chan struct{}
}

func NewWorker(st *store.Store, redis  *redis.Client) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		ctx:    ctx,
		cancel: cancel,
		st:     st,
		redis: redis,
	}
}

func (w *Worker)LoadAllUserIDsToQueue() error {
	rows, err := w.st.User().IdFromUsers()
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	ctx := context.Background()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}

		// Добавляем ID в очередь
		if err := w.redis.RPush(ctx, "reco:queue", id).Err(); err != nil {
			return fmt.Errorf("redis push error: %w", err)
		}
	}

	return nil
}

// func (w *Worker)StartRecoWorker() {
// 	go func() {
// 		ctx := context.Background()
// 		for {
// 			res, err := w.redis.BLPop(ctx, 0, "reco:queue").Result()
// 			if err != nil {
// 				log.Println("BLPOP error:", err)
// 				continue
// 			}
// 			if len(res) < 2 {
// 				continue
// 			}

// 			idStr := res[1]
// 			userID, err := strconv.Atoi(idStr)
// 			if err != nil {
// 				log.Println("Invalid user ID:", idStr)
// 				continue
// 			}

// 			log.Println("Processing recommendations for user:", userID)
// 			w.LoadPreferencesUser(userID)
// 		}
// 	}()
// }


func (w *Worker) StartRecoWorker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Reco worker stopped")
				return
			default:
				res, err := w.redis.BLPop(ctx, 0, "reco:queue").Result()
				if err != nil {
					if ctx.Err() != nil {
						// Завершение по контексту
						return
					}
					log.Println("BLPOP error:", err)
					continue
				}
				if len(res) < 2 {
					continue
				}

				idStr := res[1]
				userID, err := strconv.Atoi(idStr)
				if err != nil {
					log.Println("Invalid user ID:", idStr)
					continue
				}

				log.Println("Processing recommendations for user:", userID)
				w.LoadPreferencesUser(userID)
			}
		}
	}()
}

func (w *Worker) LoadPreferencesUser(userID int) error {
	rows, err := w.st.User().IdPreferencesUser(userID)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	ctx := context.Background()

	for rows.Next() {
		var candidateID int
		if err := rows.Scan(&candidateID); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}

		// Добавляем предпочтение — ID другого пользователя
		if err := w.redis.RPush(ctx, fmt.Sprintf("user:%d:preferences", userID), candidateID).Err(); err != nil {
			return fmt.Errorf("redis push error: %w", err)
		}
	}

	return nil
}

func (w *Worker) Stop() {
	w.cancel()
}
