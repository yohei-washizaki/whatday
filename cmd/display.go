package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

func displayEvent(e Event, showDescription bool) {
	fmt.Println(e.Title)
	if !showDescription {
		return
	}

	// ロケールに沿って、記念日のフォーマットを調整する
	dateFormatted, err := FormatDateForLocale(e, viper.GetString("locale"))
	if err == nil {
		fmt.Println(dateFormatted)
	}
	fmt.Println(e.Description)
}
