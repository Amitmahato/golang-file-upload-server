package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	fmt.Println("Invoked load env")
	err := godotenv.Load(".env")
	fmt.Println("Success loading .env file")
	if err != nil {
		log.Fatal("Failure loading .env file")
	}
}

func GetEnvWithKey(key string) string {
	value := os.Getenv(key)
	return value
}