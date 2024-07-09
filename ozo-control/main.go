package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"ozo-control/src"
	"ozo-control/src/config"

	"github.com/MatusOllah/slogcolor"
	"github.com/go-rod/rod"
	"github.com/urfave/cli/v2"
)

func main() {
	ops := slogcolor.DefaultOptions
	ops.SrcFileMode = slogcolor.ShortFile
	ops.SrcFileLength = 30
	slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, ops)))
	app := &cli.App{
		Name:  "ozo-control",
		Usage: "Manage OZOの勤怠管理を自動化するCLIツール",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "設定ファイルを作成します",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "user-id",
						Aliases:  []string{"u"},
						Usage:    "OZOのユーザーID",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Aliases:  []string{"p"},
						Usage:    "OZOのパスワード",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					userId := cCtx.String("user-id")
					password := cCtx.String("password")
					browser := rod.New().MustConnect()
					_, err := src.SignIn(browser, userId, password)
					if err != nil {
						slog.Error("Failed to sign in")
						slog.Error("Please check your user ID and password")
						config.Clean()
						return err
					}
					err = config.Init(userId, password)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "check-in",
				Aliases: []string{"i"},
				Usage:   "出勤ボタンを押します",
				Action: func(cCtx *cli.Context) error {
					err := checkIn()
					if err != nil {
						log.Fatal(err)
					}
					slog.Info("Successfully checked in")
					return nil
				},
			},
			{
				Name:    "check-out",
				Aliases: []string{"o"},
				Usage:   "退勤ボタンを押します",
				Action: func(cCtx *cli.Context) error {
					err := checkOut()
					if err != nil {
						log.Fatal(err)
					}
					slog.Info("Successfully checked out")
					return nil
				},
			},
			{
				Name:    "register-holiday",
				Aliases: []string{"r"},
				Usage:   "休暇を登録します",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "override",
						Aliases:     []string{"o"},
						Usage:       "登録済みの勤務時間を上書きします",
						DefaultText: "false",
					},
				},
				Action: func(cCtx *cli.Context) error {
					override := cCtx.Bool("override")
					err := registerHoliday(override)
					if err != nil {
						log.Fatal(err)
					}
					slog.Info("Successfully registered holiday")
					return nil
				},
			},
			{
				Name:    "clean",
				Aliases: []string{"c"},
				Usage:   "設定ファイルと実行ファイルを削除します",
				Action: func(cCtx *cli.Context) error {
					err := config.Clean()
					if err != nil {
						log.Fatal(err)
					}
					slog.Info("Successfully uninstalled")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

const recreateErrorMsg = "Please recreate the config file using `ozo-control init -u <user-id> -p <password>`"

func checkIn() error {
	config, err := config.New()
	if err != nil {
		slog.Error(recreateErrorMsg)
		return err
	}

	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page, err := src.SignIn(browser, config.UserId, config.Password)
	if err != nil {
		slog.Error(recreateErrorMsg)
		return err
	}

	page, err = src.CheckIn(page)
	if err != nil {
		return err
	}
	page.Close()

	return nil
}

func checkOut() error {
	config, err := config.New()
	if err != nil {
		slog.Error(recreateErrorMsg)
		return err
	}

	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page, err := src.SignIn(browser, config.UserId, config.Password)
	if err != nil {
		slog.Error(recreateErrorMsg)
		return err
	}

	page, err = src.CheckOut(page)
	if err != nil {
		return err
	}
	page.Close()

	return nil
}

func registerHoliday(override bool) error {
	slog.Info(fmt.Sprintf("override: %v", override))

	config, err := config.New()
	if err != nil {
		slog.Error(recreateErrorMsg)
		return err
	}

	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page, err := src.SignIn(browser, config.UserId, config.Password)
	if err != nil {
		slog.Error(recreateErrorMsg)
		return err
	}

	page, err = src.RegisterHoliday(browser, page, override)
	if err != nil {
		return err
	}
	page.Close()

	return nil
}
