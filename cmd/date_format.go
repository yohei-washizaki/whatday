package cmd

import (
	"fmt"
	"time"

	"github.com/goodsign/monday"
)

func FormatDateForLocale(e Event, locale string) (string, error) {
	switch e.Frequency {
	case "yearly":
		return FormatDateForLocaleYearly(e.Date, locale)
	case "monthly":
		return FormatDateForLocaleMonthly(e.Date, locale)
	default:
		return "", fmt.Errorf("unsupported frequency: %s", e.Frequency)
	}
}

func FormatDateForLocaleYearly(dateStr, locale string) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	var dateFormatted string
	var error error
	switch locale {
	case "JaJP":
		dateFormatted, error = monday.Format(t, "2006年1月2日", monday.LocaleJaJP), nil
	case "EnUS":
		dateFormatted, error = monday.Format(t, "January 2, 2006", monday.LocaleEnUS), nil
	default:
		dateFormatted, error = t.Format("2006-01-02"), nil
	}
	return dateFormatted, error
}

func FormatDateForLocaleMonthly(dateStr, locale string) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	var dateFormatted string
	switch locale {
	case "JaJP":
		dateFormatted, err = monday.Format(t, "2日", monday.LocaleJaJP), nil
		dateFormatted = "毎月" + dateFormatted
	case "EnUS":
		dateFormatted, err = monday.Format(t, "2", monday.LocaleEnUS), nil
		dateFormatted = "Every month on the " + dateFormatted
	default:
		dateFormatted, err = t.Format("02"), nil
	}
	return dateFormatted, err
}
