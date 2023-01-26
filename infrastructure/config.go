package infrastructure

import (
	"fmt"

	"github.com/joho/godotenv"
)

func NewConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}
