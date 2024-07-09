package src

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/samber/lo"

	"github.com/go-rod/rod"
)

const listPageUrl string = "https://manage.ozo-cloud.jp/buysell/default.cfm?version=buysell&app_cd=329&fuseaction=knt"
const attendanceSelector string = "#frmSearch > table > tbody > tr:nth-child(2) > td > table > tbody > tr:has(td > span > a)" // > span > a:nth-child(1)"
const dateSelector string = "td:nth-child(1) > span > a:nth-child(1)"
const reasonSelector string = "td:nth-child(3)"

func RegisterHoliday(browser *rod.Browser, page *rod.Page, override bool) (*rod.Page, error) {
	page.MustNavigate(listPageUrl)
	page.MustWaitLoad()

	attendances := page.MustElements(attendanceSelector)
	holidays, err := getPublicHoliday()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	slog.Info("running parallel...")
	for i, att := range attendances {
		dateStr := att.MustElement(dateSelector).MustText()

		date, err := parseRespToDate(dateStr)
		if err != nil {
			return nil, err
		}

		c := lo.CountBy(holidays, func(h time.Time) bool {
			return h.Equal(date)
		})

		reason := att.MustElement(reasonSelector).MustText()
		isAlreadyRegistered := reason != "\xc2\xa0"
		if isAlreadyRegistered && !override {
			slog.Info(dateStr + " is skipped because it is already registered")
			continue
		}

		if !isWeekend(date) && c == 0 {
			slog.Info(dateStr + " is skipped because it is weekday")
			continue // 休日登録は平日のみ
		}

		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			p := browser.MustPage(listPageUrl)
			p.MustWaitStable()

			a := p.MustElements(attendanceSelector)[i]
			dateStr := a.MustElement(dateSelector).MustText()

			date, err := parseRespToDate(dateStr)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			if isWeekend(date) {
				slog.Info(dateStr + " is weekend. Registering...")
				a.MustElement(dateSelector).MustClick()
				p.MustWaitStable()
				resisterHoliday(p)
				slog.Info("Registered holiday " + dateStr)
				return
			}

			if c > 0 {
				slog.Info(dateStr + " is public holiday. Registering...")
				a.MustElement(dateSelector).MustClick()
				p.MustWaitStable()
				resisterHoliday(p)
				slog.Info("Registered holiday " + dateStr)
				return
			}
		}(i)
	}
	wg.Wait()
	return page, nil
}

func parseRespToDate(date string) (time.Time, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Time{}, err
	}
	layout := "01月02日"
	t, err := time.ParseInLocation(layout, date[0:10], jst)
	if err != nil {
		return time.Time{}, err
	}
	t = t.AddDate(time.Now().Year(), 0, 0).In(jst)
	return t, nil
}

func isWeekend(date time.Time) bool {
	return date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
}

func getPublicHoliday() ([]time.Time, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("https://holidays-jp.github.io/api/v1/%d/date.json", time.Now().Year())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	json.Unmarshal(b, &result)

	holidays := []time.Time{}
	for k := range result {
		t, err := time.ParseInLocation("2006-01-02", k, jst)
		if err != nil {
			return nil, err
		}
		holidays = append(holidays, t)
	}

	return holidays, nil
}

const holidaySelectorInList string = "#tmp_JIYUU_SYUJITU"
const registerButtonSelector string = "#btnRegist" // cspell:disable-line

func resisterHoliday(page *rod.Page) *rod.Page {
	selector := page.MustElement(holidaySelectorInList)
	selector.MustSelect("公休")

	wait, handle := page.MustHandleDialog()
	go page.MustElement(registerButtonSelector).MustClick()
	wait()
	handle(true, "")
	page.MustWaitStable()
	return page
}
