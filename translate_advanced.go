package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
	"google.golang.org/api/option"
)

func translateAdvanced() error {
	if projectID == "" {
		return fmt.Errorf("project ID is required for Advanced API")
	}

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
	var client *translate.TranslationClient
	var clientErr error

	if credentials != "" {
		client, clientErr = translate.NewTranslationClient(ctx, option.WithCredentialsFile(credentials))
	} else {
		client, clientErr = translate.NewTranslationClient(ctx)
	}

	if clientErr != nil {
		return fmt.Errorf("failed to create client: %v", clientErr)
	}
	defer client.Close()

	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		Contents:           []string{string(content)},
		TargetLanguageCode: targetLang,
		MimeType:          "text/plain",
	}

	if sourceLang != "auto" {
		req.SourceLanguageCode = sourceLang
	}

	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to translate text: %v", err)
	}

	if len(resp.GetTranslations()) == 0 {
		return fmt.Errorf("no translation returned")
	}

	translatedText := resp.GetTranslations()[0].GetTranslatedText()
	err = ioutil.WriteFile(outputFile, []byte(translatedText), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	fmt.Printf("Successfully translated %s to %s using Advanced API\n", inputFile, outputFile)
	fmt.Printf("Source language: %s, Target language: %s\n", 
		resp.GetTranslations()[0].GetDetectedLanguageCode(), targetLang)
	
	return nil
}
