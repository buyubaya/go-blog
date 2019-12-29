package config


type Config struct {
	DB *DBConfig
	Firebase *FirebaseConfig
	Redis *RedisConfig
}


type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Charset  string
}


// FIREBASE
type FirebaseConfig struct {
	ServiceAccountKey string
}


// REDIS
type RedisConfig struct {
	Address string
	Password string
}


func GetConfig() *Config {
	return &Config {
		DB: &DBConfig {
			Dialect:  "mysql",
			Username: "XIX96fMr4k",
			Password: "cqJJGWzIHn",
			Host:     "remotemysql.com",
			Port:     3306,
			Name:     "XIX96fMr4k",
			Charset:  "utf8",
		},

		Firebase: &FirebaseConfig {
			ServiceAccountKey: "./firebase/serviceAccountKey.json",
		},

		Redis: &RedisConfig {
			Address: "localhost:6379",
			Password: "",
		},

	}
}