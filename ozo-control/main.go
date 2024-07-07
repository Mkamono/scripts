package main

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	// ヘッドレスブラウザを起動する
	url := launcher.New().MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()

	// スクレイピング対象のページを指定する
	page := browser.MustPage("https://zenn.dev/kenta_goto")

	// ページが完全にロードされるのを待つ
	page.MustWaitLoad()

	// セレクターに一致する要素を取得する
	elements := page.MustElements("h3")

	// 各要素のテキストを出力する
	for _, element := range elements {
		fmt.Println(element.MustText())
	}
}
