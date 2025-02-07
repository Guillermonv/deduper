package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	Dir  = "/Users/guillermovarelli/workspace" // Change this to target your workspace
	Path = "/deduper/src/examples/test_data.csv"
)

type DeduperService struct{}

func NewDeduperService() *DeduperService {
	return &DeduperService{}
}

// FindDuplicates reads the CSV and calculates match score for each row and writes results to output file
func (fs *DeduperService) FindDuplicates() (string, error) {
	// Combine the constants for directory and path
	filePath := Dir + Path

	// Check if the CSV file exists before trying to open it
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("CSV file does not exist at %s", filePath)
		}
		return "", fmt.Errorf("error checking file existence: %w", err)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read CSV: %w", err)
	}

	// Check if the header is as expected (6 columns)
	headers := records[0]
	if len(headers) != 6 {
		return "", fmt.Errorf("unexpected CSV format: expected 6 columns, found %d", len(headers))
	}

	var result []string
	// Loop through each row (skip header row)
	for i, record := range records[1:] {
		if len(record) != 6 {
			return "", fmt.Errorf("invalid record format at row %d: expected 6 fields, found %d", i+1, len(record))
		}

		contactID := record[0]
		name := record[1]
		lname := record[2]
		email := record[3]
		postalZip := record[4]
		address := record[5]

		// Compare this contact with all subsequent contacts (to avoid reverse comparisons)
		for j := i + 1; j < len(records)-1; j++ {
			sourceRecord := records[j+1]
			if len(sourceRecord) != 6 {
				return "", fmt.Errorf("invalid record format at row %d: expected 6 fields, found %d", j+1, len(sourceRecord))
			}

			sourceContactID := sourceRecord[0]
			sourceName := sourceRecord[1]
			sourceLname := sourceRecord[2]
			sourceEmail := sourceRecord[3]
			sourcePostalZip := sourceRecord[4]
			sourceAddress := sourceRecord[5]

			// Calculate the total score based on field comparisons
			var totalScore float64
			totalScore += scoreField(name, sourceName, 5)
			totalScore += scoreField(lname, sourceLname, 5)
			totalScore += scoreField(email, sourceEmail, 20)
			totalScore += scoreField(postalZip, sourcePostalZip, 30)
			totalScore += scoreField(address, sourceAddress, 40)

			// Calculate the match percentage
			percentage := totalScore

			// Store the result in the format: ContactID, Source ContactID, Match Score
			result = append(result, fmt.Sprintf("ContactID: %s, Source ContactID: %s, Match Score: %.2f", contactID, sourceContactID, percentage))
		}
	}

	// Create the output file in the same directory as the CSV file
	outputPath := filepath.Join(filepath.Dir(filePath), "match_results.csv")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Write the results to the output file
	writer := csv.NewWriter(outputFile)
	// Writing headers
	writer.Write([]string{"ContactID", "Source ContactID", "Match Score"})
	for _, res := range result {
		parts := strings.Split(res, ", ")
		writer.Write([]string{
			parts[0][len("ContactID: "):],
			parts[1][len("Source ContactID: "):],
			parts[2][len("Match Score: "):],
		})
	}
	writer.Flush()

	return outputPath, nil
}

// scoreField compares the input with the CSV field and returns the score based on weight
func scoreField(recordValue, inputValue string, weight float64) float64 {
	// Normalize values (lowercase and trim spaces)
	recordValue = strings.TrimSpace(strings.ToLower(recordValue))
	inputValue = strings.TrimSpace(strings.ToLower(inputValue))

	// If exact match, return full weight
	if inputValue == recordValue {
		return weight
	}

	// Compute Levenshtein similarity score
	similarity := 1 - float64(levenshteinDistance(recordValue, inputValue))/float64(max(len(recordValue), len(inputValue)))

	// Apply weight
	return similarity * weight
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Initialize previous row
	prevRow := make([]int, len(s2)+1)
	for i := range prevRow {
		prevRow[i] = i
	}

	// Compute distances row by row
	for i, c1 := range s1 {
		currRow := make([]int, len(s2)+1)
		currRow[0] = i + 1

		for j, c2 := range s2 {
			cost := 0
			if c1 != c2 {
				cost = 1
			}
			currRow[j+1] = min(
				currRow[j]+1,    // Insertion
				prevRow[j+1]+1,  // Deletion
				prevRow[j]+cost, // Substitution
			)
		}
		prevRow = currRow
	}

	return prevRow[len(s2)]
}
