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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.etcd.io/bbolt"
)

//go:embed data/wdayin/*.json
var wdayinData embed.FS

const kDataRoot = "data/wdayin"

var cfgFile string
var showAll bool
var inputDate string
var DB *bbolt.DB

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

		var today time.Time
		var matchLevel string
		if inputDate == "" {
			today = time.Now()
			matchLevel = "today"
		} else {
			today, err = time.Parse("2006-01-02", inputDate)
			matchLevel = "exact"
			if err != nil {
				today, err = time.Parse("01-02", inputDate)
				if err == nil {
					today = today.AddDate(9999-today.Year(), 0, 0)
				}
				matchLevel = "month-day"
				if err != nil {
					today, err = time.Parse("02", inputDate)
					if err == nil {
						today = today.AddDate(9999-today.Year(), 1-int(today.Month()), 0)
					}
					matchLevel = "day"
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error parsing input date: %v\n", err)
						os.Exit(1)
						return
					}
				}
			}
		}

		var eventsFound []Event
		for _, e := range events {
			if eventMatches(e, today, matchLevel) {
				eventsFound = append(eventsFound, e)
			}
		}

		if len(eventsFound) == 0 {
			return
		}

		if showAll {
			for _, e := range eventsFound {
				displayEvent(e, false)
			}
			return
		}

		eventSelected := eventsFound[randomIndex(len(eventsFound))]
		displayEvent(eventSelected, false)
	},
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func eventMatches(e Event, today time.Time, matchLevel string) bool {
	eventDate, err := time.Parse("2006-01-02", e.Date)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing event date: %v\n", err)
		os.Exit(1)
		return false
	}

	switch matchLevel {
	case "today":
		if e.Frequency == "yearly" {
			return eventDate.Month() == today.Month() && eventDate.Day() == today.Day()
		}
		return eventDate.Day() == today.Day()
	case "exact":
		return eventDate.Year() == today.Year() && eventDate.Month() == today.Month() && eventDate.Day() == today.Day()
	case "month-day":
		return eventDate.Month() == today.Month() && eventDate.Day() == today.Day()
	case "day":
		return eventDate.Day() == today.Day()
	default:
		return false
	}
}

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
	cobra.OnInitialize(initializeApp)
	cobra.OnFinalize(finalizeApp)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wday.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all events found.")
	rootCmd.Flags().StringVarP(&inputDate, "date", "d", "", "Specify a date to search for events (YYYY-MM-DD, MM-DD or DD).")
}

var configPath string

// initializeApp reads in config file and ENV variables if set.
func initializeApp() {
	// find home directory
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath = filepath.Join(home, ".config", "wday")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := os.Mkdir(configPath, 0755)
		if err != nil {
			fmt.Println("Error creating config directory:", err)
			return
		}
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		cobra.CheckErr(err)

		// Search config in home directory with name ".wday" (without extension).
		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config") // name of config file (without extension)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			var configFile string
			if cfgFile != "" {
				configFile = cfgFile
			} else {
				configFile = filepath.Join(configPath, "config.yaml")
			}

			// デフォルト設定ファイルを作成
			defaultConfig := []byte("locale: JaJP\n")
			err := os.WriteFile(configFile, defaultConfig, 0644)
			if err != nil {
				fmt.Println("Error creating default config file:", err)
				return
			}

			// 作成した設定ファイルを読み込む
			viper.SetConfigFile(configFile)
			if err := viper.ReadInConfig(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading default config file:", err)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Error reading config file:", err)
		}
	}

	// Make database cache directory if it doesn't exist at configPath/db
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	dbDir := filepath.Join(home, ".cache", "wday", "db")
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err := os.MkdirAll(dbDir, 0755)
		if err != nil {
			fmt.Println("Error creating database cache directory:", err)
			return
		}
	}

	locale := viper.GetString("locale")
	dbPath := filepath.Join(dbDir, locale+".db")

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Create the database file
		DB, err = bbolt.Open(dbPath, 0600, nil)
		if err != nil {
			fmt.Println("Error creating database:", err)
			return
		}

		// Write the dataset to the database
		datasetPath := filepath.Join(kDataRoot, locale+".json")

		data, err := wdayinData.ReadFile(datasetPath)
		if err != nil {
			fmt.Println("Error reading dataset file:", err)
			return
		}

		err = DB.Update(func(tx *bbolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("Events"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			return bucket.Put([]byte("data"), data)
		})
		if err != nil {
			fmt.Println("Error writing dataset to database:", err)
			return
		}
	} else {
		DB, err = bbolt.Open(dbPath, 0600, nil)
		if err != nil {
			fmt.Println("Error opening database:", err)
			return
		}
	}
}

func finalizeApp() {
	if DB != nil {
		DB.Close()
	}
}
