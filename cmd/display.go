package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

func displayEvent(e Event, showDescription bool) {
	output := e.Title
	if showDescription {
		// ロケールに沿って、記念日のフォーマットを調整する
		dateFormatted, err := FormatDateForLocale(e, viper.GetString("locale"))
		if err == nil {
			output += "\t" + dateFormatted
		}
	}
	fmt.Println(output)
}
