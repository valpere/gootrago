/*
Copyright © 2025 Valentyn Solomko <valentyn.solomko@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Translate CSV files or specific columns",
	Long: `A flexible CSV translation tool that can translate entire files or specific columns while preserving the original structure. 
Supports both Basic and Advanced Google Cloud Translation APIs and various CSV formats.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("csv called")
	// },
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputFile == outputFile {
			return fmt.Errorf("input file and output file are the same: %v", inputFile)
		}

		csv, err := readCSVToSlice(inputFile, false, csvDelimiter, csvComment)
		if err != nil {
			return fmt.Errorf("failed to read CSV file: %v", err)
		}

		colNumbers, err := decodeColNumbers(csvColumn, len(csv[0]))
		if err != nil {
			return err
		}
		nCols := len(colNumbers)
		for i, row := range csv {
			if nCols == 0 {
				strOut, err := translateEx(row, useAdvanced)
				if err != nil {
					return fmt.Errorf("failed to translate text: %v", err)
				}
				csv[i] = strOut
			} else {
				strInp := make([]string, 0, nCols)
				for _, v := range colNumbers {
					strInp = append(strInp, row[v-1])
				}
				strOut, err := translateEx(strInp, useAdvanced)
				if err != nil {
					return fmt.Errorf("failed to translate text: %v", err)
				}
				for k, v := range strOut {
					row[colNumbers[k]-1] = v
				}
			}
		}

		return writeSliceToCSV(outputFile, csv, nil, csvDelimiter)
	},
}

func init() {
	rootCmd.AddCommand(csvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// csvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// csvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	csvCmd.Flags().StringSliceVarP(&csvColumn, "column", "l", []string{}, "One or many columns number to translate (can be specified multiple times). Numeration starts from '1' or 'A'")
	csvCmd.Flags().StringVarP(&csvDelimiter, "csv-delimiter", "", "", "Delimiter for CSV files")
	csvCmd.Flags().StringVarP(&csvComment, "csv-comment", "", "", "Comment character for CSV files")
}

func decodeColNumbers(csvColumn []string, csvWidth int) ([]int, error) {
	var colNumbers []int = make([]int, 0, len(csvColumn))
	for _, col := range csvColumn {
		col = strings.ToUpper(strings.TrimSpace(col))
		var colNumber int
		var err error
		if (col[0] >= 'A') && (col[0] <= 'Z') {
			colNumber = titleToNumber(col)
		} else {
			colNumber, err = strconv.Atoi(col)
			if err != nil {
				return nil, fmt.Errorf("invalid column number: %v", col)
			}
		}

		if (colNumber < 1) || (colNumber > csvWidth) {
			return nil, fmt.Errorf("column number is out of range: %v", col)
		}
		colNumbers = append(colNumbers, colNumber)
	}

	return colNumbers, nil
}

func titleToNumber(columnTitle string) int {
	l := len(columnTitle)
	if l < 1 {
		return 0
	}

	res := 0
	for _, c := range columnTitle {
		if (c < 'A') || (c > 'Z') {
			return 0
		}

		res = res*26 + int(c-'A'+1)
	}

	return res
}
