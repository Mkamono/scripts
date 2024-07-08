package src

import (
	"errors"
	"log/slog"

	"github.com/go-rod/rod"
)

const (
	checkInButtonSelector string = "#btn03"                                                                                                                                                                                                                                                                               // 出勤ボタン
	checkInTimeSelector   string = "body > div.desktop_wrapper > div > table > tbody > tr > td:nth-child(1) > div > div > div > div.clearfix.mytool-bg > form > table:nth-child(79) > tbody > tr > td:nth-child(2) > table > tbody > tr > td:nth-child(2) > table.BaseDesign > tbody > tr:nth-child(3) > td:nth-child(3)" // 出勤時間(実績)
)

func CheckIn(page *rod.Page) (*rod.Page, error) {
	// 出勤済みかどうかを確認
	checkInTime := page.MustElement(checkInTimeSelector).MustText()
	if checkInTime != "\xc2\xa0" { // \xc2\xa0 は &nbsp; のエスケープ
		slog.Warn("Already checked in at " + checkInTime)
		return page, nil
	}

	// 出勤ボタンをクリック
	checkInButton := page.MustElement(checkInButtonSelector)
	checkInButton.MustClick()

	page.MustWaitStable()

	// 出勤時間を取得
	checkInTime = page.MustElement(checkInTimeSelector).MustText()
	if checkInTime == "\xc2\xa0" {
		return nil, errors.New("failed to check in")
	}

	slog.Info("Successfully checked in at " + checkInTime)

	return page, nil
}
