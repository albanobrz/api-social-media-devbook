package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ConnectionDBString      = ""
	ConnectionDBStringMongo = ""
	Port                    = 0
	SecretKey               []byte
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
	ConnectionDBString = fmt.Sprintf("%s:%s@tcp(172.19.0.2:3306)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	ConnectionDBStringMongo := fmt.Sprintf("mongodb://%s:%s@172.19.0.3:27017/",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	// fmt.Println(ConnectionDBStringMongo)

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
