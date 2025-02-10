/*
Copyright © 2025 Valentyn Solomko <valentyn.solomko@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global variables to store command-line flags and configuration
var (
	cfgFile     string // Path to configuration file
	inputFile   string // Path to input file for translation
	outputFile  string // Path where translated text will be saved
	sourceLang  string // Source language code (e.g., 'en' for English)
	targetLang  string // Target language code (e.g., 'es' for Spanish)
	projectID   string // Google Cloud Project ID (required for Advanced API)
	credentials string // Path to Google Cloud credentials JSON file
	useAdvanced bool   // Flag to switch between Basic and Advanced APIs
)

func readInp() (string, error) {
	// Read the strInp of the input file
	strInp, err := os.ReadFile(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read input file: %v", err)
	}

	return string(strInp), nil
}

func writeOut(strOut []string) error {
	// Write the translated text to the output file
	fh, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer fh.Close()

	for _, str := range strOut {
		_, err = fh.WriteString(str)
		if err != nil {
			return fmt.Errorf("failed to write to output file: %v", err)
		}
	}

	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gootrago",
	Short: "CLI Google Translator written on Golang",
	Long: `A CLI application that translates text files using Google Translate API.
It supports both Basic and Advanced Google Translate APIs and various language options.
The Basic API is simpler but has fewer features, while the Advanced API offers more control but requires a Google Cloud Project ID.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	// RunE is used instead of Run to allow error handling
	RunE: func(cmd *cobra.Command, args []string) error {
		strInp, err := readInp()
		if err != nil {
			return fmt.Errorf("failed to read input file: %v", err)
		}

		var strOut []string
		// Choose between Basic and Advanced API based on the flag
		if useAdvanced {
			strOut, err = translateAdvanced([]string{strInp})
			// fmt.Printf("Successfully translated %s to %s using Advanced API\n", inputFile, outputFile)
		} else {
			strOut, err = translateBasic([]string{strInp})
			// fmt.Printf("Successfully translated %s to %s using Basic API\n", inputFile, outputFile)
		}

		if err != nil {
			return fmt.Errorf("failed to translate text: %v", err)
		}

		// Ensure the output directory exists
		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}

		return writeOut(strOut)
	},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gootrago.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Local flags (only available to this command)
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file to translate (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for translation (required)")
	rootCmd.Flags().StringVarP(&sourceLang, "source", "s", "auto", "Source language code (e.g., 'en' for English)")
	rootCmd.Flags().StringVarP(&targetLang, "target", "t", "", "Target language code (e.g., 'uk' for Ukrainian) (required)")
	rootCmd.Flags().StringVarP(&projectID, "project", "p", "", "Google Cloud Project ID (required for advanced API)")
	rootCmd.Flags().StringVarP(&credentials, "credentials", "c", "", "Path to Google Cloud credentials JSON file")
	rootCmd.Flags().BoolVarP(&useAdvanced, "advanced", "a", false, "Use Advanced Google Translate API")

	// Mark required flags
	// These flags must be provided or the application will show an error
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
	rootCmd.MarkFlagRequired("target")
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

		// Search config in home directory with name ".gootrago" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gootrago")
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig() // Find and read the config file
	if err == nil {
		// If a config file is found, read it in.
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

}
