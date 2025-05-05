package models

type Photo struct {
    User_id   int `json:"user_id"`
    Photos []string `json:"photos"`
}
