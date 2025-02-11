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
// Note: This function relies on the global variable 'inputFile' being set
// before calling. Make sure inputFile contains a valid file path.
// --------------------------------------------------------------------------
func readInp() (string, error) {
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
//
// Note: This function relies on the global variable 'outputFile' being set
// before calling. Make sure outputFile contains a valid file path.
// --------------------------------------------------------------------------
func writeOut(strOut []string) error {
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
// readCSVToSlice reads a CSV file and returns a slice of string slices
// where each inner slice represents a row from the CSV file.
// Parameters:
//   - filename: path to the CSV file
//   - hasHeader: true if the CSV file has a header row that should be skipped
//
// Returns:
//   - [][]string: slice of slices containing the CSV data
//   - error: nil if successful, error message if failed
//
// --------------------------------------------------------------------------
func readCSVToSlice(filename string, hasHeader bool) ([][]string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

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
// writeSliceToCSV writes a slice of string slices to a CSV file.
// Parameters:
//   - filename: path where the CSV file should be created/overwritten
//   - data: slice of string slices containing the data to write
//   - header: optional header row (can be nil if no header is needed)
//
// Returns:
//   - error: nil if successful, error message if failed
//
// --------------------------------------------------------------------------
func writeSliceToCSV(filename string, data [][]string, header []string) error {
	// Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

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
