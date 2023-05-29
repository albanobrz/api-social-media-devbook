package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ConnectionDBString = ""
	Port               = 0
	SecretKey          []byte
)

func Load() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 9000
	}

	ConnectionDBString = os.Getenv("DB_MYSQL_URI")

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
