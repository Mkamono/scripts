package config

import (
	"errors"
	"fmt"
	"log"
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

func New() (*Config, error) {
	// .envファイルが存在しない場合は作成する
	envFilePath := ".env"
	if !exist(envFilePath) {
		slog.Error("No .env file found")
		err := createEnvFile(envFilePath)
		if err != nil {
			log.Fatal(err)
		}
		slog.Error("Created .env file")
		abs, err := filepath.Abs(envFilePath)
		if err != nil {
			return nil, err
		}

		slog.Error(fmt.Sprintf("Please fill in USER_ID and PASSWORD in %s", abs))
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

	slog.Info(fmt.Sprintf("Loaded config: %+v\n", c))

	return &c, nil
}

func Uninstall() error {
	envFilePath := ".env"
	if !exist(envFilePath) {
		return nil
	}

	err := os.Remove(envFilePath)
	if err != nil {
		return err
	}
	return nil
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func createEnvFile(filename string) error {
	f, e := os.Create(filename)
	if e != nil {
		return e
	}
	defer f.Close()

	_, e = f.WriteString("OZO_CONTROL_USER_ID=\nOZO_CONTROL_PASSWORD=\n")
	if e != nil {
		return e
	}
	return nil
}
