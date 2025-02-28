/*
Copyright © 2025 Yohei WASHIZAKI <yohei.washizaki@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setLocale bool

// localeCmd represents the locale command
var localeCmd = &cobra.Command{
	Use:   "locale",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if setLocale {
			localeCode := args[0]
			_, ok := GetLocaleByCode(localeCode)
			if !ok {
				fmt.Fprintf(os.Stderr, "Unsupported locale: %s\n", localeCode)
				fmt.Fprintf(os.Stderr, "Supported locales are:\n")
				for _, loc := range SupportedLocalesMap {
					printLocale(loc, true, true)
				}
				os.Exit(1)
				return
			}
			viper.Set("locale", localeCode)
			viper.WriteConfig()
			return
		}

		localeCode := viper.GetString("locale")
		locale, ok := GetLocaleByCode(localeCode)

		if !ok {
			return
		}

		printLocale(locale, showDesc, true)
	},
}

type Locale struct {
	Code        string
	DisplayName string
}

var SupportedLocalesMap = map[string]Locale{
	"JaJP": {Code: "JaJP", DisplayName: "日本語"},
	"EnUS": {Code: "EnUS", DisplayName: "English(US)"},
}

func GetLocaleByCode(code string) (Locale, bool) {
	loc, ok := SupportedLocalesMap[code]
	return loc, ok
}

var showDesc bool

func printLocale(locale Locale, displayName bool, toStdErr bool) {
	if displayName {
		if toStdErr {
			fmt.Fprintf(os.Stderr, "%s, %s\n", locale.Code, locale.DisplayName)
		} else {
			fmt.Printf("%s, %s\n", locale.Code, locale.DisplayName)
		}
		return
	}
	if toStdErr {
		fmt.Fprintf(os.Stderr, "%s\n", locale.Code)
	} else {
		fmt.Printf("%s\n", locale.Code)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, loc := range SupportedLocalesMap {
			printLocale(loc, showDesc, true)
		}
	},
}

var setCmd = &cobra.Command{
	Use:   "set [locale code]",
	Short: "Set the locale",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		localeCode := args[0]
		_, ok := GetLocaleByCode(localeCode)
		if !ok {
			fmt.Fprintf(os.Stderr, "Unsupported locale: %s\n", localeCode)
			fmt.Fprintf(os.Stderr, "Supported locales are:\n")
			for _, loc := range SupportedLocalesMap {
				printLocale(loc, true, true)
			}
			os.Exit(1)
			return
		}
		viper.Set("locale", localeCode)
		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(localeCmd)
	localeCmd.AddCommand(listCmd)
	localeCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	localeCmd.Flags().BoolVar(&showDesc, "desc", false, "Show description of locale.")
	localeCmd.Flags().BoolVar(&setLocale, "set", false, "Set locale.")
}
