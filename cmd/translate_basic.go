package cmd

/*
This file implements translation using the Basic Google Translate API.
The Basic API is simpler to use and doesn't require a project ID, but has fewer features compared to the Advanced API.
*/

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

// translateBasic handles translation using the Basic Google Translate API
func translateBasic(strInp []string) (strOut []string, err error) {
	// Set up Google Cloud credentials if provided
	if credentials != "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentials)
	}

	// Initialize the translation client
	ctx := context.Background()
	var client *translate.Client
	var clientErr error

	// Create client with or without explicit credentials
	if credentials != "" {
		client, clientErr = translate.NewClient(ctx, option.WithCredentialsFile(credentials))
	} else {
		client, clientErr = translate.NewClient(ctx)
	}

	if clientErr != nil {
		return strOut, fmt.Errorf("failed to create client: %v", clientErr)
	}
	defer client.Close()

	// Parse the target language code
	targetLangTag, err := language.Parse(targetLang)
	if err != nil {
		return strOut, fmt.Errorf("invalid target language code: %v", err)
	}

	var translations []translate.Translation
	if sourceLang == "auto" {
		// If source language is auto, let the API detect it
		translations, err = client.Translate(ctx,
			strInp,
			targetLangTag,
			&translate.Options{
				Format: translate.Text,
			})
	} else {
		// If source language is specified, parse and use it
		sourceLangTag, err := language.Parse(sourceLang)
		if err != nil {
			return strOut, fmt.Errorf("invalid source language code: %v", err)
		}

		translations, err = client.Translate(ctx,
			strInp,
			targetLangTag,
			&translate.Options{
				Source: sourceLangTag,
				Format: translate.Text,
			})
	}

	if err != nil {
		return strOut, fmt.Errorf("failed to translate text: %v", err)
	}

	if len(translations) == 0 {
		return strOut, fmt.Errorf("no translation returned")
	}

	// if sourceLang == "auto" {
	// 	detectedLang := translations[0].Source
	// 	fmt.Printf("Detected source language: %s, Target language: %s\n", detectedLang, targetLang)
	// } else {
	// 	fmt.Printf("Source language: %s, Target language: %s\n", sourceLang, targetLang)
	// }

	for _, tra := range translations {
		strOut = append(strOut, tra.Text)
	}

	return strOut, nil
}
