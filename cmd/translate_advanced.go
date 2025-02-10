package cmd

/*
This file implements translation using the Advanced Google Translate API (v3).
The Advanced API requires a project ID but offers more features and control over the translation process.
It's suitable for production environments where more detailed control and monitoring is needed.
*/

import (
	"context"
	"fmt"
	"os"

	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"google.golang.org/api/option"
)

// translateAdvanced handles translation using the Advanced Google Translate API
func translateAdvanced(strInp string) (string, error) {
	// Verify project ID is provided (required for Advanced API)
	if projectID == "" {
		return "", fmt.Errorf("project ID is required for Advanced API")
	}

	// Set up Google Cloud credentials if provided
	if credentials != "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentials)
	}

	// Initialize the translation client
	ctx := context.Background()
	var client *translate.TranslationClient
	var clientErr error

	// Create client with or without explicit credentials
	if credentials != "" {
		client, clientErr = translate.NewTranslationClient(ctx,
			option.WithCredentialsFile(credentials))
	} else {
		client, clientErr = translate.NewTranslationClient(ctx)
	}

	if clientErr != nil {
		return "", fmt.Errorf("failed to create client: %v", clientErr)
	}
	defer client.Close()

	// Prepare the translation request
	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		Contents:           []string{strInp},
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
		return "", fmt.Errorf("failed to translate text: %v", err)
	}

	if len(resp.GetTranslations()) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	translatedText := resp.GetTranslations()[0].GetTranslatedText()

	fmt.Printf("Source language: %s, Target language: %s\n", resp.GetTranslations()[0].GetDetectedLanguageCode(), targetLang)

	return translatedText, nil
}
