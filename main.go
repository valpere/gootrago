package main

/*
This application provides a command-line interface for translating text files using either
the Basic or Advanced Google Translate API. It uses Cobra for CLI handling and Viper for
configuration management, supporting both command-line flags and configuration files.

The application is split into three main parts:
1. CLI setup and configuration (main.go)
2. Basic API implementation (translate_basic.go)
3. Advanced API implementation (translate_advanced.go)
*/

import (
	"fmt"
	"os"

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

func main() {
	Execute()
}

// rootCmd represents the base command when called without any subcommands
// This is the main entry point for the CLI application
var rootCmd = &cobra.Command{
	Use:   "translator",
	Short: "A text file translator using Google Translate",
	Long: `A CLI application that translates text files using Google Translate API.
It supports both Basic and Advanced Google Translate APIs and various language options.
The Basic API is simpler but has fewer features, while the Advanced API offers more
control but requires a Google Cloud Project ID.`,
	// RunE is used instead of Run to allow error handling
	RunE: func(cmd *cobra.Command, args []string) error {
		// Choose between Basic and Advanced API based on the flag
		if useAdvanced {
			return translateAdvanced()
		}
		return translateBasic()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init is called automatically by Cobra to initialize the command-line flags
func init() {
	// Initialize configuration before executing any commands
	cobra.OnInitialize(initConfig)

	// Global persistent flags (available to all commands)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.translator.yaml)")

	// Local flags (only available to this command)
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "",
		"Input file to translate (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "",
		"Output file for translation (required)")
	rootCmd.Flags().StringVarP(&sourceLang, "source", "s", "auto",
		"Source language code (e.g., 'en' for English)")
	rootCmd.Flags().StringVarP(&targetLang, "target", "t", "",
		"Target language code (e.g., 'es' for Spanish) (required)")
	rootCmd.Flags().StringVarP(&projectID, "project", "p", "",
		"Google Cloud Project ID (required for advanced API)")
	rootCmd.Flags().StringVarP(&credentials, "credentials", "c", "",
		"Path to Google Cloud credentials JSON file")
	rootCmd.Flags().BoolVarP(&useAdvanced, "advanced", "a", false,
		"Use Advanced Google Translate API")

	// Mark required flags
	// These flags must be provided or the application will show an error
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
	rootCmd.MarkFlagRequired("target")
}

// initConfig reads in config file and ENV variables if set.
// This allows for persistent configuration across runs.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag if provided
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory for default config location
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".translator" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".translator")
	}

	// Read environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
