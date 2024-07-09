package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UserId   string `required:"true" envconfig:"USER_ID"`
	Password string `required:"true" envconfig:"PASSWORD"`
}

const envFilePath = ".env"

func Init(userId string, password string) error {
	// .envファイルが存在する場合は削除する
	if exist(envFilePath) {
		err := os.Remove(envFilePath)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(envFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("OZO_CONTROL_USER_ID=%s\nOZO_CONTROL_PASSWORD=%s\n", userId, password))
	if err != nil {
		return err
	}

	envFilePathAbs, err := filepath.Abs(envFilePath)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Created .env file with USER_ID: %s, PASSWORD: %s in %s\n", userId, password, envFilePathAbs))
	return nil
}

func Clean() error {
	if !exist(envFilePath) {
		return nil
	}

	err := os.Remove(envFilePath)
	if err != nil {
		return err
	}
	return nil
}

func New() (*Config, error) {
	if !exist(envFilePath) {
		slog.Error("No .env file found")
		return nil, errors.New("no .env file found")
	}

	err := godotenv.Load(envFilePath)
	if err != nil {
		return nil, errors.New("error loading .env file")
	}
	slog.Info("Loaded .env file")

	var c Config
	err = envconfig.Process("ozo_control", &c)
	if err != nil {
		return nil, err
	}

	if c.UserId == "" || c.Password == "" {
		return nil, errors.New("USER_ID and PASSWORD must be set")
	}

	slog.Info(fmt.Sprintf("Loaded config: %+v", c))
	return &c, nil
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
