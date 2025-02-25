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
	"embed"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"encoding/json"

	"github.com/goodsign/monday"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed data/wdayin/*.json
var wdayinData embed.FS

// Event はデータセット内の各記念日イベントを表現する構造体です。
type Event struct {
	ID          int    `json:"id" yaml:"id"`
	Date        string `json:"date" yaml:"date"`           // "MM-DD"形式で記録（例: "02-22"）
	Recurring   bool   `json:"recurring" yaml:"recurring"` // 毎年繰り返すイベントなら true
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
}

const kDataRoot = "data/wdayin"

var cfgFile string
var showDescription bool
var showAll bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wday",
	Short: "CLI tool revealing daily historical events, notable birthdays, and observances.",
	Long: `whatday is a simple CLI tool that reveals historical events, notable birthdays,
and interesting observances for any given day. Built with Go and Cobra, 
this project serves as a personal learning exercise in crafting concise, effective CLI applications.`,
	Run: func(cmd *cobra.Command, args []string) {
		var datasetPath string
		locale := viper.GetString("locale")
		switch locale {
		case "EnUS":
			datasetPath = filepath.Join(kDataRoot, "EnUS.json")
		case "JaJP":
			datasetPath = filepath.Join(kDataRoot, "JaJP.json")
		default:
			datasetPath = filepath.Join(kDataRoot, "JaJP.json")
		}

		data, err := wdayinData.ReadFile(datasetPath)
		if err != nil {
			fmt.Println("Error reading dataset file:", err)
			return
		}

		// Parse the JSON data into a slice of Event structs
		var events []Event
		if err := json.Unmarshal(data, &events); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON data: %v\n", err)
			return
		}

		today := time.Now()
		var eventsFound []Event
		for _, e := range events {
			d, err := time.Parse("2006-01-02", e.Date)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing events data: %d\n", e.ID)
				os.Exit(1)
				return
			}

			if d.Month() == today.Month() && d.Day() == today.Day() {
				eventsFound = append(eventsFound, e)
			}
		}

		if len(eventsFound) == 0 {
			return
		}

		if showAll {
			for _, e := range eventsFound {
				displayEvent(e, showDescription)
			}
			return
		}

		eventSelected := eventsFound[randomIndex(len(eventsFound))]
		displayEvent(eventSelected, showDescription)
	},
}

func displayEvent(e Event, showDescription bool) {
	fmt.Println(e.Title)
	if showDescription {
		// ロケールに沿って、記念日のフォーマットを調整する
		dateFormatted, err := FormatDateForLocale(e.Date, viper.GetString("locale"))
		if err == nil {
			fmt.Println(dateFormatted)
		}
		fmt.Println(e.Description)
	}
}

func FormatDateForLocale(dateStr, locale string) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	switch locale {
	case "JaJP":
		return monday.Format(t, "2006年1月2日", monday.LocaleJaJP), nil
	case "EnUS":
		return monday.Format(t, "January 2, 2006", monday.LocaleEnUS), nil
	default:
		return t.Format("2006-01-02"), nil
	}
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomIndex(max int) int {
	return rnd.Intn(max)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wday.yaml)")
	rootCmd.PersistentFlags().String("locale", "JaJP", "Locale code for the default dataset (e.g. JaJP, EnUS)")
	viper.BindPFlag("locale", rootCmd.PersistentFlags().Lookup("locale"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&showDescription, "description", "d", false, "Show event descriptions")
	rootCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all events found.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".wday" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".wday")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			var configFile string
			if cfgFile != "" {
				configFile = cfgFile
			} else {
				home, err := os.UserHomeDir()
				cobra.CheckErr(err)
				configFile = home + "/.wday.yaml"
			}
			defaultConfig := []byte("locale: JaJP\n")
			err := os.WriteFile(configFile, defaultConfig, 0644)
			if err != nil {
				fmt.Println("Error creating default config file:", err)
			} else {
				fmt.Println("Default config file created at:", configFile)
				viper.SetConfigFile(configFile)
				if err := viper.ReadInConfig(); err != nil {
					fmt.Fprintln(os.Stderr, "Error reading default config file:", err)
				}
			}
		} else {
			fmt.Fprintln(os.Stderr, "Error reading config file:", err)
		}
	}
}
