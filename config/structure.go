package config

type Config struct {
	PGConfig    PostgresConfig
	MongoConfig MongoConfig
}

type PostgresConfig struct {
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
}

type MongoConfig struct {
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
}
