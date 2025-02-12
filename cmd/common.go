/*
Package provides file operations for reading input and writing output files.
These functions are designed to handle basic file operations with proper error
handling and resource management.
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

// **************************************************************************
// readInp reads the contents of the input file specified by the global variable
// inputFile and returns them as a string. This function is designed to handle
// the entire file content at once, making it suitable for small to medium-sized
// text files that can fit comfortably in memory.
//
// The function performs the following steps:
// 1. Reads the entire content of the file using os.ReadFile
// 2. Converts the byte slice to a string
// 3. Returns the string along with any potential error
//
// Returns:
//   - string: The entire contents of the input file as a string
//   - error: An error if any occurred during file reading, nil otherwise
//
// Error cases:
//   - File does not exist
//   - Permission denied
//   - File is not readable
//   - Memory allocation fails for very large files
//
// Usage example:
//
//	content, err := readInp()
//	if err != nil {
//	    log.Fatalf("Error reading input: %v", err)
//	}
//
// --------------------------------------------------------------------------
func readInp(inputFile string) (string, error) {
	strInp, err := os.ReadFile(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read the input file: %v", err)
	}

	return string(strInp), nil
}

// **************************************************************************
// writeOut writes a slice of strings to the output file specified by the global
// variable outputFile. Each string in the slice is written sequentially to the
// file. This function is designed for cases where you need to write multiple
// strings to a file while maintaining fine control over the writing process.
//
// The function performs the following steps:
// 1. Creates (or truncates) the output file
// 2. Writes each string from the input slice sequentially
// 3. Automatically closes the file when done using defer
//
// Parameters:
//   - strOut []string: A slice of strings to be written to the file
//
// Returns:
//   - error: An error if any occurred during file operations, nil otherwise
//
// Error cases:
//   - Cannot create output file (permissions, invalid path)
//   - Write operation fails (disk full, file system errors)
//   - File system becomes read-only during writing
//
// Usage example:
//
//	strings := []string{"First line\n", "Second line\n"}
//	err := writeOut(strings)
//	if err != nil {
//	    log.Fatalf("Error writing output: %v", err)
//	}
//
// Important considerations:
// 1. The function creates a new file or truncates an existing one
// 2. Each string is written exactly as provided - no automatic newlines
// 3. The file handle is properly closed even if errors occur
// 4. The function uses defer for safe resource cleanup
// --------------------------------------------------------------------------
func writeOut(outputFile string, strOut []string) error {
	fh, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create the output file: %v", err)
	}
	defer fh.Close()

	for _, str := range strOut {
		_, err = fh.WriteString(str)
		if err != nil {
			return fmt.Errorf("failed to write to the output file: %v", err)
		}
	}

	return nil
}

// **************************************************************************
// readCSVToSlice reads a CSV file and converts its contents into a slice of string slices.
// Each inner slice represents one row from the CSV file, with individual fields as elements.
// The function provides flexibility in handling different CSV formats through customizable
// delimiter and comment characters.
//
// Parameters:
//   - filename string: The path to the CSV file to be read
//   - hasHeader bool: Indicates whether the first row should be treated as a header
//     and excluded from the returned data. When true, the first row is skipped.
//   - csvDelimiter string: The character to use as the field separator. If empty,
//     the default comma (,) is used. Only the first character of the string is used
//     if multiple characters are provided.
//   - csvComment string: The character to use for marking comment lines. If empty,
//     no comment handling is performed. Only the first character of the string is
//     used if multiple characters are provided.
//
// Returns:
//   - [][]string: A two-dimensional slice where each inner slice represents a row
//     from the CSV file. Each element in the inner slice represents a field from
//     that row.
//   - error: An error if any occurred during file operations or CSV parsing.
//     The error includes context about what operation failed.
//
// The function handles several scenarios:
//  1. Missing or inaccessible files
//  2. Malformed CSV data
//  3. Empty files
//  4. Files with or without headers
//  5. Custom delimiters and comment characters
//
// Example usage:
//
//	// Reading a standard CSV file with header
//	data, err := readCSVToSlice("data.csv", true, "", "")
//
//	// Reading a tab-delimited file with comments
//	data, err := readCSVToSlice("data.tsv", false, "\t", "#")
//
// Note: The function loads the entire file into memory. For very large files,
// consider implementing a streaming approach instead.
// --------------------------------------------------------------------------
func readCSVToSlice(filename string, hasHeader bool, csvDelimiter string, csvComment string) ([][]string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)
	if csvDelimiter != "" {
		reader.Comma = rune(csvDelimiter[0])
	}
	if csvComment != "" {
		reader.Comment = rune(csvComment[0])
	}

	// Read all records at once
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	// If file has header and there are records, skip the first row
	if hasHeader && len(records) > 0 {
		records = records[1:]
	}

	return records, nil
}

// **************************************************************************
// writeSliceToCSV writes data from a slice of string slices to a CSV file.
// The function can optionally include a header row and use a custom delimiter.
// It's designed to handle the inverse operation of readCSVToSlice, allowing
// for round-trip processing of CSV data.
//
// Parameters:
//   - filename string: The path where the CSV file should be created or overwritten
//   - data [][]string: The data to write, structured as a slice of rows, where
//     each row is a slice of strings representing individual fields
//   - header []string: Optional slice of strings to use as the header row.
//     If nil, no header row is written
//   - csvDelimiter string: The character to use as the field separator. If empty,
//     the default comma (,) is used. Only the first character of the string is
//     used if multiple characters are provided.
//
// Returns:
//   - error: An error if any occurred during file creation, writing operations,
//     or CSV encoding. The error includes context about what operation failed.
//
// The function handles several important aspects:
//  1. Creates a new file or truncates an existing one
//  2. Properly handles file closing through defer
//  3. Manages CSV encoding with custom delimiters
//  4. Provides clear error messages for different failure scenarios
//
// Example usage:
//
//	// Writing a simple CSV file with header
//	data := [][]string{
//	    {"John", "30", "New York"},
//	    {"Alice", "25", "Los Angeles"},
//	}
//	header := []string{"Name", "Age", "City"}
//	err := writeSliceToCSV("output.csv", data, header, "")
//
//	// Writing a tab-delimited file without header
//	err := writeSliceToCSV("output.tsv", data, nil, "\t")
//
// Important considerations:
//  1. The function overwrites any existing file
//  2. All rows should have the same number of fields
//  3. The function automatically handles proper CSV encoding,
//     including escaping special characters
//  4. The Flush operation ensures all data is written to disk
//
// Note: For very large datasets, consider implementing a streaming
// approach that writes rows incrementally.
// --------------------------------------------------------------------------
func writeSliceToCSV(filename string, data [][]string, header []string, csvDelimiter string) error {
	// Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if csvDelimiter != "" {
		writer.Comma = rune(csvDelimiter[0])
	}

	// Write header if provided
	if header != nil {
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing header: %v", err)
		}
	}

	// Write all records
	if err := writer.WriteAll(data); err != nil {
		return fmt.Errorf("error writing data: %v", err)
	}

	return nil
}

func indicator(shutdownCh <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Print(".")
		case <-shutdownCh:
			return
		}
	}
}
