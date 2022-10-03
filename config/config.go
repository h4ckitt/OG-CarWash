package config

import (
	"github.com/joho/godotenv"
	"os"
)

var conf Config

func Load() error {
	err := godotenv.Load(".env")

	if err != nil {
		return err
	}

	conf = Config{
		PGConfig: PostgresConfig{
			DatabaseHost:     os.Getenv("PG_HOST"),
			DatabasePort:     os.Getenv("PG_PORT"),
			DatabaseUser:     os.Getenv("PG_USER"),
			DatabasePassword: os.Getenv("PG_PASS"),
			DatabaseName:     os.Getenv("PG_NAME"),
		},
		MongoConfig: MongoConfig{
			DatabaseUser:     os.Getenv("MG_USER"),
			DatabasePassword: os.Getenv("MG_PASS"),
			DatabasePort:     os.Getenv("MG_PORT"),
			DatabaseName:     os.Getenv("MG_NAME"),
			DatabaseHost:     os.Getenv("MG_HOST"),
		},
		FirebaseConfig: FirebaseConfig{
			ServiceFileName: os.Getenv("FIREBASE_SERVICE_JSON"),
		},
	}

	return nil
}

func GetConfig() Config {
	return conf
}
