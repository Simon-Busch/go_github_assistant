package utils

import (
	"github.com/joho/godotenv"
	"os"
)

func InitEnv() (string, string, error) {
	err := godotenv.Load()
	if err != nil {
			return "", "", err
	}

	ghUserName := os.Getenv("GITHUB_USERNAME")
	ghToken := os.Getenv("GITHUB_TOKEN");
	return ghUserName, ghToken, nil
}
