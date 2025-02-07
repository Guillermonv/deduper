package test

import (
	"encoding/csv"
	"os"
	"testing"

	"deduper/src/internal"

	"github.com/stretchr/testify/assert"
)

// Mocking a function to return predefined records
func mockCSVData() [][]string {
	return [][]string{
		{"ContactID", "Name", "Lname", "Email", "PostalZip", "Address"},
		{"1001", "John", "Doe", "john.doe@example.com", "12345", "123 Main St"},
		{"1002", "Jane", "Doe", "jane.doe@example.com", "12346", "124 Main St"},
		{"1003", "John", "Doe", "john.doe@example.com", "12345", "123 Main St"},
	}
}

// Mocking `os.Open` to avoid actual file reading
func mockOpenFile(fileName string) (*os.File, error) {
	// Create a temporary file with mock CSV content
	mockData := mockCSVData()

	// Creating a temporary file to simulate CSV reading
	tempFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		return nil, err
	}

	// Write mock data to the temporary file
	writer := csv.NewWriter(tempFile)
	defer writer.Flush()

	for _, record := range mockData {
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	// Rewind to the beginning of the file for reading
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

// Test case for FindDuplicates with a successful scenario
func TestFindDuplicates_Success(t *testing.T) {
	// Create an instance of the DeduperService
	deduper := internal.NewDeduperService()

	// Mock opening the file (using mockOpenFile)
	file, err := mockOpenFile("mocked.csv")
	if err != nil {
		t.Fatalf("Failed to mock open file: %v", err)
	}
	defer file.Close()

	// Here we mock the method `FindDuplicates` by passing a mocked file
	// and simulating a result of calling `FindDuplicates`
	outputPath, err := deduper.FindDuplicates()

	// Validate if there were no errors
	assert.NoError(t, err, "FindDuplicates should run without errors")

	// Check that the outputPath contains the expected result file name
	assert.Contains(t, outputPath, "match_results.csv", "Output path should contain 'match_results.csv'")

	// Simulate checking the output by reading the result CSV file
	outputFile, err := os.Open(outputPath)
	assert.NoError(t, err, "Should be able to open the result file")

	reader := csv.NewReader(outputFile)
	records, err := reader.ReadAll()
	assert.NoError(t, err, "Error reading the CSV output file")

	// Validate that the expected number of results is present
	assert.Greater(t, len(records), 1, "Expected at least one match result")

	// Check that the first result contains the correct information format
	expectedHeader := []string{"ContactID", "Source ContactID", "Match Score"}
	assert.Equal(t, records[0], expectedHeader, "The header row should match the expected output format")
}

// Test case to check if only the address matches, so the score should be 40
func TestFindDuplicates_AddressMatchScore(t *testing.T) {
	deduper := internal.NewDeduperService()

	// Mock file for testing
	file, err := mockOpenFile("mocked_address.csv")
	if err != nil {
		t.Fatalf("Failed to mock open file: %v", err)
	}
	defer file.Close()

	// Run FindDuplicates to simulate result
	outputPath, err := deduper.FindDuplicates()

	// Validate if no errors occurred
	assert.NoError(t, err, "FindDuplicates should run without errors")

	// Open and check output CSV
	outputFile, err := os.Open(outputPath)
	assert.NoError(t, err, "Should be able to open the result file")

	reader := csv.NewReader(outputFile)
	records, err := reader.ReadAll()
	assert.NoError(t, err, "Error reading the CSV output file")

	// Validate that the score for address match is 40
	for _, record := range records[1:] {
		contactID := record[0]
		sourceContactID := record[1]
		matchScore := record[2]

		// Assuming contactID 1001 and sourceContactID 1002 match on the address
		if contactID == "1001" && sourceContactID == "1002" {
			assert.Equal(t, matchScore, "40.00", "Address match should have a score of 40")
		}
	}
}

// Test case to check if only the name matches, so the score should be 5
func TestFindDuplicates_NameMatchScore(t *testing.T) {
	deduper := internal.NewDeduperService()

	// Mock file for testing
	file, err := mockOpenFile("mocked_name.csv")
	if err != nil {
		t.Fatalf("Failed to mock open file: %v", err)
	}
	defer file.Close()

	// Run FindDuplicates to simulate result
	outputPath, err := deduper.FindDuplicates()

	// Validate if no errors occurred
	assert.NoError(t, err, "FindDuplicates should run without errors")

	// Open and check output CSV
	outputFile, err := os.Open(outputPath)
	assert.NoError(t, err, "Should be able to open the result file")

	reader := csv.NewReader(outputFile)
	records, err := reader.ReadAll()
	assert.NoError(t, err, "Error reading the CSV output file")

	// Validate that the score for name match is 5
	for _, record := range records[1:] {
		contactID := record[0]
		sourceContactID := record[1]
		matchScore := record[2]

		// Assuming contactID 1003 and sourceContactID 1004 match on the name
		if contactID == "1003" && sourceContactID == "1004" {
			assert.Equal(t, matchScore, "5.00", "Name match should have a score of 5")
		}
	}
}

// Test case to check if only the postal zip matches, so the score should be 30
func TestFindDuplicates_PostalZipMatchScore(t *testing.T) {
	deduper := internal.NewDeduperService()

	// Mock file for testing
	file, err := mockOpenFile("mocked_postal_zip.csv")
	if err != nil {
		t.Fatalf("Failed to mock open file: %v", err)
	}
	defer file.Close()

	// Run FindDuplicates to simulate result
	outputPath, err := deduper.FindDuplicates()

	// Validate if no errors occurred
	assert.NoError(t, err, "FindDuplicates should run without errors")

	// Open and check output CSV
	outputFile, err := os.Open(outputPath)
	assert.NoError(t, err, "Should be able to open the result file")

	reader := csv.NewReader(outputFile)
	records, err := reader.ReadAll()
	assert.NoError(t, err, "Error reading the CSV output file")

	// Validate that the score for postal zip match is 30
	for _, record := range records[1:] {
		contactID := record[0]
		sourceContactID := record[1]
		matchScore := record[2]

		// Assuming contactID 1005 and sourceContactID 1006 match on the postal zip
		if contactID == "1005" && sourceContactID == "1006" {
			assert.Equal(t, matchScore, "30.00", "Postal Zip match should have a score of 30")
		}
	}
}

// Test case to check if multiple fields match and score is combined accordingly
func TestFindDuplicates_MultipleFieldsMatchScore(t *testing.T) {
	deduper := internal.NewDeduperService()

	// Mock file for testing
	file, err := mockOpenFile("mocked_multiple_fields.csv")
	if err != nil {
		t.Fatalf("Failed to mock open file: %v", err)
	}
	defer file.Close()

	// Run FindDuplicates to simulate result
	outputPath, err := deduper.FindDuplicates()

	// Validate if no errors occurred
	assert.NoError(t, err, "FindDuplicates should run without errors")

	// Open and check output CSV
	outputFile, err := os.Open(outputPath)
	assert.NoError(t, err, "Should be able to open the result file")

	reader := csv.NewReader(outputFile)
	records, err := reader.ReadAll()
	assert.NoError(t, err, "Error reading the CSV output file")

	// Validate that the score for multiple fields match is combined (name + address)
	for _, record := range records[1:] {
		contactID := record[0]
		sourceContactID := record[1]
		matchScore := record[2]

		// Assuming contactID 1007 and sourceContactID 1008 match on both name (5) and address (40)
		if contactID == "1007" && sourceContactID == "1008" {
			assert.Equal(t, matchScore, "45.00", "Name and address match should have a combined score of 45")
		}
	}
}
