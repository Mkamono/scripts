package src

import (
	"errors"
	"log/slog"

	"github.com/go-rod/rod"
)

const topPageUrl string = "https://manage.ozo-cloud.jp/buysell/default.cfm?version=buysell"

func SignIn(browser *rod.Browser, userID string, password string) (*rod.Page, error) {
	page := browser.MustPage(topPageUrl)
	page.MustElement("input[name='LoginUserID']").MustInput(userID)
	page.MustElement("input[name='LoginPassword']").MustInput(password)
	page.MustElement("input[type='submit']").MustClick()
	page.MustWaitStable()

	if page.MustHas("#logout") {
		slog.Info("Successfully signed in")
	} else {
		return nil, errors.New("failed to sign in")
	}
	return page, nil
}
