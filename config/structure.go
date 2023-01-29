package config

type Config struct {
	PGConfig       PostgresConfig
	MongoConfig    MongoConfig
	FirebaseConfig FirebaseConfig
	ImageConfig    ImageConfig
	RunConfig      RunConfig
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

type ImageConfig struct {
	Template string
	Location string
}

type RunConfig struct {
	Port string
}
