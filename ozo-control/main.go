package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"ozo-control/src"
	"ozo-control/src/config"

	"github.com/go-rod/rod"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ozo-control",
		Usage: "Manage OZOの勤怠管理を自動化するCLIツール",
		Commands: []*cli.Command{
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func checkIn() error {
	config, err := config.New()
	if err != nil {
		return err
	}

	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page, err := src.SignIn(browser, config.UserId, config.Password)
	if err != nil {
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
		return err
	}

	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page, err := src.SignIn(browser, config.UserId, config.Password)
	if err != nil {
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
		return err
	}

	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page, err := src.SignIn(browser, config.UserId, config.Password)
	if err != nil {
		return err
	}

	page, err = src.RegisterHoliday(browser, page, override)
	if err != nil {
		return err
	}
	page.Close()

	return nil
}
