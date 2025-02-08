package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

func translateBasic() error {
	if credentials != "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentials)
	}

	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	ctx := context.Background()
	var client *translate.Client
	var clientErr error

	if credentials != "" {
		client, clientErr = translate.NewClient(ctx, option.WithCredentialsFile(credentials))
	} else {
		client, clientErr = translate.NewClient(ctx)
	}

	if clientErr != nil {
		return fmt.Errorf("failed to create client: %v", clientErr)
	}
	defer client.Close()

	targetLangTag, err := language.Parse(targetLang)
	if err != nil {
		return fmt.Errorf("invalid target language code: %v", err)
	}

	var translations []translate.Translation
	if sourceLang == "auto" {
		translations, err = client.Translate(ctx,
			[]string{string(content)},
			targetLangTag,
			&translate.Options{
				Format: translate.Text,
			})
	} else {
		sourceLangTag, err := language.Parse(sourceLang)
		if err != nil {
			return fmt.Errorf("invalid source language code: %v", err)
		}

		translations, err = client.Translate(ctx,
			[]string{string(content)},
			targetLangTag,
			&translate.Options{
				Source: sourceLangTag,
				Format: translate.Text,
			})
	}

	if err != nil {
		return fmt.Errorf("failed to translate text: %v", err)
	}

	if len(translations) == 0 {
		return fmt.Errorf("no translation returned")
	}

	err = ioutil.WriteFile(outputFile, []byte(translations[0].Text), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	fmt.Printf("Successfully translated %s to %s using Basic API\n", inputFile, outputFile)
	
	if sourceLang == "auto" {
		detectedLang := translations[0].Source
		fmt.Printf("Detected source language: %s, Target language: %s\n", 
			detectedLang, targetLang)
	} else {
		fmt.Printf("Source language: %s, Target language: %s\n", 
			sourceLang, targetLang)
	}
	
	return nil
}
