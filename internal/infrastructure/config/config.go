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

	// local mysql connection:
	// ConnectionDBString = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
	// container db connection:
	ConnectionDBString = os.Getenv("DB_MYSQL_URI")

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
