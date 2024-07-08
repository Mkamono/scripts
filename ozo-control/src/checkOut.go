package src

import (
	"errors"
	"log/slog"

	"github.com/go-rod/rod"
)

const (
	checkOutButtonSelector string = "#btn04"                                                                                                                                                                                                                                                                               // 退勤ボタン
	checkOutTimeSelector   string = "body > div.desktop_wrapper > div > table > tbody > tr > td:nth-child(1) > div > div > div > div.clearfix.mytool-bg > form > table:nth-child(79) > tbody > tr > td:nth-child(2) > table > tbody > tr > td:nth-child(2) > table.BaseDesign > tbody > tr:nth-child(3) > td:nth-child(4)" // 退勤時間(実績)
)

func CheckOut(page *rod.Page) (*rod.Page, error) {
	// 退勤済みかどうかを確認
	checkOutTime := page.MustElement(checkOutTimeSelector).MustText()
	if checkOutTime != "\xc2\xa0" { // \xc2\xa0 は &nbsp; のエスケープ
		slog.Warn("Already checked out at " + checkOutTime)
		return page, nil
	}

	// 退勤ボタンをクリック
	checkOutButton := page.MustElement(checkOutButtonSelector)
	checkOutButton.MustClick()

	page.MustWaitStable()

	// 退勤時間を取得
	checkOutTime = page.MustElement(checkOutTimeSelector).MustText()
	if checkOutTime == "\xc2\xa0" {
		return nil, errors.New("failed to check out")
	}

	slog.Info("Successfully checked out at " + checkOutTime)

	return page, nil
}
