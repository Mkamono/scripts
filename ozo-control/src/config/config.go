package config

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UserId   string `required:"true" envconfig:"USER_ID"`
	Password string `required:"true" envconfig:"PASSWORD"`
}

func New() (*Config, error) {
	// .envファイルが存在しない場合は作成する
	if !exist(".env") {
		slog.Error("No .env file found")
		err := copy(".env.template", ".env")
		if err != nil {
			log.Fatal(err)
		}
		slog.Error("Created .env file from .env.template")
		slog.Error("Please fill in the .env file")
		return nil, errors.New("no .env file found")
	}

	err := godotenv.Load(".env")
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

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func copy(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}
