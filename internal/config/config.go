package config

import (
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/yesyash/yaskin-backend/internal/logger"
)

type Config struct {
	Port       int    `validate:"required"`
	AppEnv     string `validate:"appEnv"`
	DbHost     string `validate:"required"`
	DbPort     int    `validate:"required"`
	Database   string `validate:"required"`
	DbUsername string `validate:"required"`
	DbPassword string `validate:"required"`
}

var (
	Port       int
	AppEnv     string // prod | staging | dev
	DbHost     string
	DbPort     int
	Database   string
	DbUsername string
	DbPassword string
)

/*
* Loads the environment variables from the .env file
* --
 */
func loadEnv() {
	env := os.Getenv("APP_ENV")

	/*
	  for "prod" or "staging" we don't want to load the .env file
	  we assume that the environment variables are already set
	*/
	if env == "prod" || env == "staging" {
		return
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

/*
* Converts a string to an integer
* --
* @param value - string
* @returns - integer on successful conversion or panics if the conversion fails
 */
func convertToInt(value string) int {
	intValue, err := strconv.Atoi(value)

	if err != nil {
		panic(err)
	}

	return intValue
}

/*
* Validates the APP_ENV environment variable
* --
 */
func validateAppEnv(fl validator.FieldLevel) bool {
	env := fl.Field().String()

	if env != "prod" && env != "staging" && env != "dev" {
		return false
	}

	return true
}

func init() {
	loadEnv()

	config := Config{
		Port:       convertToInt(os.Getenv("PORT")),
		AppEnv:     os.Getenv("APP_ENV"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     convertToInt(os.Getenv("DB_PORT")),
		Database:   os.Getenv("DB_DATABASE"),
		DbUsername: os.Getenv("DB_USERNAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
	}

	validate := validator.New()

	validate.RegisterValidation("appEnv", validateAppEnv)

	err := validate.Struct(config)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			logger.Error(err.Error())
		}

		logger.Fatal("Invalid environment variables")
	}

	Port = config.Port
	AppEnv = config.AppEnv
	DbPort = config.DbPort
	Database = config.Database
	DbUsername = config.DbUsername
	DbPassword = config.DbPassword
	DbHost = config.DbHost
}
