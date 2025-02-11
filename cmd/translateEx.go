/*
This file provides translation functionality using both Basic and Advanced Google
Translate APIs. The package implements a flexible translation system that can handle
multiple input strings and switch between API versions based on requirements.

The translation system is built around three main functions:
1. translateEx - The main entry point that orchestrates the translation process
2. translateBasic - Handles translation using the Basic Google Translate API
3. translateAdvanced - Handles translation using the Advanced Google Translate API v3
*/
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
// translateEx serves as the main entry point for the translation system,
// orchestrating the translation process by delegating to either the Basic
// or Advanced Google Translate API based on user preference.
//
// This function acts as a facade, abstracting the complexity of choosing
// and using different translation APIs behind a simple interface.
//
// Parameters:
//   - strInp []string: A slice of strings to be translated. Each string
//     in the slice will be translated individually while maintaining order.
//   - useAdvanced bool: Determines which API version to use:
//   - true: Uses the Advanced API (requires projectID)
//   - false: Uses the Basic API
//
// Returns:
//   - []string: A slice containing the translated strings, maintaining
//     the same order as the input slice
//   - error: An error if any occurred during translation, nil otherwise
//
// The function relies on several global variables being set:
//   - credentials: Path to Google Cloud credentials file (optional)
//   - projectID: Required for Advanced API
//   - sourceLang: Source language code (or "auto" for detection)
//   - targetLang: Target language code
//
// Usage example:
//
//	input := []string{"Hello", "World"}
//	translated, err := translateEx(input, false)
//	if err != nil {
//	    log.Fatalf("Translation failed: %v", err)
//	}
//
// Note: This function preserves the order of translations, ensuring that
// each translated string corresponds to its original input string.
// --------------------------------------------------------------------------
func translateEx(strInp []string, useAdvanced bool) (strOut []string, err error) {
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

// **************************************************************************
// translateBasic handles translation using the Basic Google Translate API.
// This implementation is simpler and doesn't require a project ID, making
// it suitable for basic translation needs.
//
// The function handles:
// 1. Client initialization with optional credentials
// 2. Language parsing and validation
// 3. Automatic language detection when sourceLang is "auto"
// 4. Batch translation of multiple strings
//
// Parameters:
//   - strInp []string: Slice of strings to translate
//
// Returns:
//   - []string: Slice of translated strings
//   - error: Any error that occurred during translation
//
// The function uses several global variables:
//   - credentials: Path to Google Cloud credentials (optional)
//   - sourceLang: Source language code or "auto"
//   - targetLang: Target language code
//
// Error cases:
//   - Invalid credentials
//   - Invalid language codes
//   - API communication failures
//   - Empty translation results
//
// Note: The Basic API is often sufficient for simple translation needs
// and doesn't require project setup in Google Cloud.
// --------------------------------------------------------------------------
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
// translateAdvanced handles translation using the Advanced Google Translate API (v3).
// This implementation provides additional features and control but requires
// a Google Cloud project ID.
//
// The function provides:
// 1. Advanced translation features
// 2. Project-based authentication
// 3. Detailed translation metadata
// 4. Batch processing capabilities
//
// Parameters:
//   - strInp []string: Slice of strings to translate
//
// Returns:
//   - []string: Slice of translated strings
//   - error: Any error that occurred during translation
//
// Required global variables:
//   - projectID: Google Cloud project ID (mandatory)
//   - credentials: Path to credentials file (optional)
//   - sourceLang: Source language code or "auto"
//   - targetLang: Target language code
//
// The Advanced API is recommended when you need:
// - Enterprise-level translation features
// - Detailed translation metadata
// - Integration with other Google Cloud services
// - Advanced monitoring and logging
//
// Error handling:
// - Validates project ID presence
// - Handles authentication errors
// - Manages API-specific errors
// - Validates translation results
//
// Note: This function requires proper Google Cloud project setup
// and appropriate API permissions.
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
