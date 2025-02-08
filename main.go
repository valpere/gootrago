package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	inputFile   string
	outputFile  string
	sourceLang  string
	targetLang  string
	projectID   string
	credentials string
	useAdvanced bool
)

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "translator",
	Short: "A text file translator using Google Translate",
	Long: `A CLI application that translates text files using Google Translate API.
It supports both Basic and Advanced Google Translate APIs and various language options.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if useAdvanced {
			return translateAdvanced()
		}
		return translateBasic()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.translator.yaml)")

	// Local flags
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file to translate (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for translation (required)")
	rootCmd.Flags().StringVarP(&sourceLang, "source", "s", "auto", "Source language code (e.g., 'en' for English)")
	rootCmd.Flags().StringVarP(&targetLang, "target", "t", "", "Target language code (e.g., 'es' for Spanish) (required)")
	rootCmd.Flags().StringVarP(&projectID, "project", "p", "", "Google Cloud Project ID (required for advanced API)")
	rootCmd.Flags().StringVarP(&credentials, "credentials", "c", "", "Path to Google Cloud credentials JSON file")
	rootCmd.Flags().BoolVarP(&useAdvanced, "advanced", "a", false, "Use Advanced Google Translate API")

	// Required flags
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
	rootCmd.MarkFlagRequired("target")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".translator")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
