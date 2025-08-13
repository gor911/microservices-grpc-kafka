package config

import (
	"errors"
	"log"
	"os"

	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/grpcclient"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/kafkaproducer"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/postgres"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/httpserver"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/logger"
	"github.com/joho/godotenv"
)

type AppEnv string

const (
	EnvLocal AppEnv = "local"
	EnvDev   AppEnv = "dev"
	EnvProd  AppEnv = "prod"
)

type App struct {
	Env AppEnv
}

type Config struct {
	App           App
	HTTP          httpserver.Config
	GRPCClient    grpcclient.Config
	Postgres      postgres.Config
	Logger        logger.Config
	KafkaProducer kafkaproducer.Config
}

func InitConfig() (Config, error) {
	// Try to load .env, but do not fail if it is missing
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using environment variables from the system/Compose")
	}

	c := Config{}

	c, err := initDBConfig(c)

	if err != nil {
		return c, err
	}

	c, err = initAppConfig(c)

	if err != nil {
		return c, err
	}

	c, err = initLoggerConfig(c)

	if err != nil {
		return c, err
	}

	c, err = initHTTPConfig(c)

	if err != nil {
		return c, err
	}

	c, err = initGRPCClientConfig(c)

	if err != nil {
		return c, err
	}

	c, err = initKafkaProducerConfig(c)

	if err != nil {
		return c, err
	}

	return c, nil
}

func initDBConfig(c Config) (Config, error) {
	c.Postgres.Username = os.Getenv("DB_USER")

	if c.Postgres.Username == "" {
		return c, errors.New("missing DB_USER")
	}

	c.Postgres.Password = os.Getenv("DB_PASSWORD")

	if c.Postgres.Password == "" {
		return c, errors.New("missing DB_PASSWORD")
	}

	c.Postgres.Host = os.Getenv("DB_HOST")

	if c.Postgres.Host == "" {
		return c, errors.New("missing DB_HOST")
	}

	c.Postgres.Port = os.Getenv("DB_PORT")

	if c.Postgres.Port == "" {
		return c, errors.New("missing DB_PORT")
	}

	c.Postgres.Database = os.Getenv("DB_NAME")

	if c.Postgres.Database == "" {
		return c, errors.New("missing DB_NAME")
	}

	return c, nil
}

func initAppConfig(c Config) (Config, error) {
	env := os.Getenv("APP_ENV")

	if env == "" {
		return c, errors.New("missing APP_ENV")
	}

	if env == string(EnvLocal) {
		c.App.Env = EnvLocal
	} else if env == string(EnvDev) {
		c.App.Env = EnvDev
	} else if env == string(EnvProd) {
		c.App.Env = EnvProd
	} else {
		return c, errors.New("unknown APP_ENV")
	}

	return c, nil
}

func initLoggerConfig(c Config) (Config, error) {
	if c.App.Env == EnvLocal {
		c.Logger.Level = logger.Debug
	} else if c.App.Env == EnvDev {
		c.Logger.Level = logger.Info
	} else if c.App.Env == EnvProd {
		c.Logger.Level = logger.Error
	}

	return c, nil
}

func initHTTPConfig(c Config) (Config, error) {
	c.HTTP.Addr = os.Getenv("HTTP_ADDR")

	if c.HTTP.Addr == "" {
		return c, errors.New("missing HTTP_ADDR")
	}

	return c, nil
}

func initGRPCClientConfig(c Config) (Config, error) {
	c.GRPCClient.Addr = os.Getenv("INVENTORY_ADDR")

	if c.GRPCClient.Addr == "" {
		return c, errors.New("missing INVENTORY_ADDR")
	}

	return c, nil
}

func initKafkaProducerConfig(c Config) (Config, error) {
	brokers := os.Getenv("KAFKA_BROKERS")

	if brokers == "" {
		return c, errors.New("missing KAFKA_BROKERS")
	}

	c.KafkaProducer.Brokers = []string{brokers}

	return c, nil
}
