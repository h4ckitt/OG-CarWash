package config

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
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
		ImageConfig: ImageConfig{
			Template: func() string {
				if prefix := os.Getenv("IMAGE_PREFIX"); prefix != "" {
					return prefix
				}
				return "license-capture"
			}(),
			Location: os.Getenv("IMAGE_STORAGE_LOCATION"),
		},
		RunConfig: RunConfig{
			Port: func() string {
				if port := os.Getenv("PORT"); port != "" {
					return port
				}
				return "8080"
			}(),
		},
	}

	if conf.ImageConfig.Location == "" {
		return errors.New("No Image Storage Location Set")
	}

	return nil
}

func GetConfig() Config {
	return conf
}
