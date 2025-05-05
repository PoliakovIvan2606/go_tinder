package main

import (
	"fmt"
	"log"
	"tinder/internal/app/apiserver"

	"github.com/BurntSushi/toml"
)

func main() {
	config := apiserver.NewConfig()
	
	_, err := toml.DecodeFile("configs/apiserver.toml", config)
	if err != nil {
		log.Fatal(err)
	}

	server := apiserver.NewApp(config)

	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	fmt.Println("run server")
}