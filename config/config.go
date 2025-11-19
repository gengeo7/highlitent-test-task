package config

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvEnum int

const (
	Development EnvEnum = iota
	Production
)

type Config struct {
	Env                    EnvEnum
	Host                   string
	Port                   int
	PostgresHost           string
	PostgresPort           int
	PostgresUser           string
	PostgresPassword       string
	PostgresDatabase       string
	PGAdminDefaultEmail    string
	PGAdminDefaultPassword string
	MigrationPath          string
}

var Conf Config

func getEnvCustom(env string, f func(string) error) (string, error) {
	val, err := getEnv(env)
	if err != nil {
		return "", err
	}

	err = f(val)
	if err != nil {
		return "", err
	}

	return val, nil
}

func getEnv(env string) (string, error) {
	val, have := os.LookupEnv(env)
	if !have {
		return "", fmt.Errorf("%s environment variable not found", env)
	}
	return val, nil
}

func getEnvInt(env string, min, max int) (int, error) {
	val, err := getEnv(env)
	if err != nil {
		return 0, err
	}
	valInt, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", env)
	}
	if valInt < min {
		return 0, fmt.Errorf("%s must be bigger than %d", env, min)
	}
	if valInt > max {
		return 0, fmt.Errorf("%s must be smaller than %d", env, max)
	}
	return valInt, nil
}

func getEnvEnum[T any](env string, constraints map[string]T) (*T, error) {
	val, err := getEnv(env)
	if err != nil {
		return nil, err
	}
	if t, have := constraints[val]; !have {
		return nil, fmt.Errorf("%s must be in constraints %v", env, constraints)
	} else {
		return &t, nil
	}
}

func Initialize() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	errs := make([]error, 0)
	port, err := getEnvInt("PORT", 1000, 100000)
	if err != nil {
		errs = append(errs, err)
	}

	postgresPort, err := getEnvInt("POSTGRES_PORT", 1000, 100000)
	if err != nil {
		errs = append(errs, err)
	}

	env, err := getEnvEnum("ENV", map[string]EnvEnum{"development": Development, "production": Production})
	if err != nil {
		errs = append(errs, err)
	}

	host, err := getEnv("HOST")
	if err != nil {
		errs = append(errs, err)
	}

	postgresHost, err := getEnv("POSTGRES_HOST")
	if err != nil {
		errs = append(errs, err)
	}

	postgresUser, err := getEnv("POSTGRES_USER")
	if err != nil {
		errs = append(errs, err)
	}

	postgresPassword, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		errs = append(errs, err)
	}

	postgresDB, err := getEnv("POSTGRES_DB")
	if err != nil {
		errs = append(errs, err)
	}

	pgAdminDefaultEmail, err := getEnvCustom("PGADMIN_DEFAULT_EMAIL", func(s string) error {
		if _, err := mail.ParseAddress(s); err != nil {
			return fmt.Errorf("PGADMIN_DEFAULT_EMAIL must be a valid email address")
		}
		return nil
	})
	if err != nil {
		errs = append(errs, err)
	}

	pgAdminDefaulPassword, err := getEnv("PGADMIN_DEFAULT_PASSWORD")
	if err != nil {
		errs = append(errs, err)
	}

	migrationPath, err := getEnv("MIGRATION_PATH")
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	Conf = Config{
		Env:                    *env,
		Host:                   host,
		Port:                   port,
		PostgresHost:           postgresHost,
		PostgresPort:           postgresPort,
		PostgresUser:           postgresUser,
		PostgresPassword:       postgresPassword,
		PostgresDatabase:       postgresDB,
		PGAdminDefaultPassword: pgAdminDefaulPassword,
		PGAdminDefaultEmail:    pgAdminDefaultEmail,
		MigrationPath:          migrationPath,
	}

	return nil
}
