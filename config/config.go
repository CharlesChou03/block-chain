package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	Version       = "0.0.1"
	DBUser        = ""
	DBPassword    = ""
	DBHost        = ""
	DBPort        = ""
	DBName        = ""
	RedisHost     = ""
	RedisPort     = ""
	RedisPassword = ""
	RPCEndpoint   = ""
)

func Setup() {
	// DB env config
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RPCEndpoint = os.Getenv("RPC_ENDPOINT")
}
