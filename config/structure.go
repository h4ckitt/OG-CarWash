package config

type Config struct {
	PGConfig       PostgresConfig
	MongoConfig    MongoConfig
	FirebaseConfig FirebaseConfig
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

type FirebaseConfig struct {
	ServiceFileName string
}
