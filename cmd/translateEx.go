package cmd

import (
	"context"
	"fmt"
	"os"

	translateBas "cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"

	translateAdv "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
)

// **************************************************************************
// translate handles translation using the Basic Google Translate API
// --------------------------------------------------------------------------
func translateEx(strInp []string) (strOut []string, err error) {
	// Choose between Basic and Advanced API based on the flag
	if useAdvanced {
		strOut, err = translateAdvanced(strInp)
		// fmt.Printf("Successfully translated %s to %s using Advanced API\n", inputFile, outputFile)
	} else {
		strOut, err = translateBasic(strInp)
		// fmt.Printf("Successfully translated %s to %s using Basic API\n", inputFile, outputFile)
	}

	if err != nil {
		return strOut, fmt.Errorf("failed to translate text: %v", err)
	}

	return strOut, nil
}

// translateBasic handles translation using the Basic Google Translate API
func translateBasic(strInp []string) (strOut []string, err error) {
	// Set up Google Cloud credentials if provided
	if credentials != "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentials)
	}

	// Initialize the translation client
	ctx := context.Background()
	var client *translateBas.Client
	var clientErr error

	// Create client with or without explicit credentials
	if credentials != "" {
		client, clientErr = translateBas.NewClient(ctx, option.WithCredentialsFile(credentials))
	} else {
		client, clientErr = translateBas.NewClient(ctx)
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

	var translations []translateBas.Translation
	if sourceLang == "auto" {
		// If source language is auto, let the API detect it
		translations, err = client.Translate(ctx,
			strInp,
			targetLangTag,
			&translateBas.Options{
				Format: translateBas.Text,
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
			&translateBas.Options{
				Source: sourceLangTag,
				Format: translateBas.Text,
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

// **************************************************************************
// translateAdvanced handles translation using the Advanced Google Translate API
// --------------------------------------------------------------------------
func translateAdvanced(strInp []string) (strOut []string, err error) {
	// Verify project ID is provided (required for Advanced API)
	if projectID == "" {
		return strOut, fmt.Errorf("project ID is required for Advanced API")
	}

	// Set up Google Cloud credentials if provided
	if credentials != "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentials)
	}

	// Initialize the translation client
	ctx := context.Background()
	var client *translateAdv.TranslationClient
	var clientErr error

	// Create client with or without explicit credentials
	if credentials != "" {
		client, clientErr = translateAdv.NewTranslationClient(ctx,
			option.WithCredentialsFile(credentials))
	} else {
		client, clientErr = translateAdv.NewTranslationClient(ctx)
	}

	if clientErr != nil {
		return strOut, fmt.Errorf("failed to create client: %v", clientErr)
	}
	defer client.Close()

	// Prepare the translation request
	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		Contents:           strInp,
		TargetLanguageCode: targetLang,
		MimeType:           "text/plain", // Specify plain text format
	}

	// Add source language if specified (not auto)
	if sourceLang != "auto" {
		req.SourceLanguageCode = sourceLang
	}

	// Perform the translation
	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return strOut, fmt.Errorf("failed to translate text: %v", err)
	}

	if len(resp.GetTranslations()) == 0 {
		return strOut, fmt.Errorf("no translation returned")
	}

	for _, tra := range resp.GetTranslations() {
		strOut = append(strOut, tra.GetTranslatedText())
	}

	// fmt.Printf("Source language: %s, Target language: %s\n", resp.GetTranslations()[0].GetDetectedLanguageCode(), targetLang)

	return strOut, nil
}
