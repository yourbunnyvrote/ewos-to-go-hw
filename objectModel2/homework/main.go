package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type OperationType string

type OperationValue string

type OperationID struct {
	Value interface{}
}

type OperationCreatedAt string

const (
	Income  OperationType = "income"
	Outcome OperationType = "outcome"
	Plus    OperationType = "+"
	Minus   OperationType = "-"
)

type FinancialRecord struct {
	Company   string             `json:"company,omitempty"`
	Operation *Operation         `json:"operation,omitempty"`
	Type      OperationType      `json:"type,omitempty"`
	Value     OperationValue     `json:"value,omitempty"`
	ID        OperationID        `json:"id,omitempty"`
	CreatedAt OperationCreatedAt `json:"created_at,omitempty"`
}

type Operation struct {
	Type      OperationType      `json:"type,omitempty"`
	Value     OperationValue     `json:"value,omitempty"`
	ID        OperationID        `json:"id,omitempty"`
	CreatedAt OperationCreatedAt `json:"created_at,omitempty"`
}

type OperationResult struct {
	CountValidOperations int
	Balance              int
	InvalidOperations    []DateWithID
}

type DateWithID struct {
	ID        interface{}
	CreatedAt OperationCreatedAt
}

type ResultJSON struct {
	Company              string        `json:"company"`
	CountValidOperations int           `json:"valid_operations_count"`
	Balance              int           `json:"balance"`
	InvalidOperationsID  []interface{} `json:"invalid_operations,omitempty"`
}

func (ot *OperationType) UnmarshalJSON(data []byte) error {
	var value string
	_ = json.Unmarshal(data, &value)

	switch OperationType(value) {
	case Income, Plus:
		*ot = Plus
		break
	case Outcome, Minus:
		*ot = Minus
		break
	}
	return nil
}

func (ov *OperationValue) UnmarshalJSON(data []byte) error {
	var rawID json.RawMessage
	if err := json.Unmarshal(data, &rawID); err != nil {
		return nil
	}

	// Attempted unmarshal into a number
	if err := ov.unmarshalNumeric(rawID); err == nil {
		return nil
	}

	// Attempted unmarshal into a string
	if err := ov.unmarshalString(rawID); err == nil {
		return nil
	}

	return nil
}

func (ov *OperationValue) unmarshalNumeric(rawID json.RawMessage) error {
	var numericValue int
	if err := json.Unmarshal(rawID, &numericValue); err == nil {
		*ov = OperationValue(strconv.Itoa(numericValue))
		return nil
	}

	var floatValue float64
	if err := json.Unmarshal(rawID, &floatValue); err == nil {
		// Check if the number is an integer
		if math.Mod(floatValue, 1) != 0 {
			*ov = ""
		} else {
			*ov = OperationValue(strconv.Itoa(int(floatValue)))
		}
		return nil
	}

	return errors.New("ID type is not int or float")
}

func (ov *OperationValue) unmarshalString(rawID json.RawMessage) error {
	var stringValue string
	if err := json.Unmarshal(rawID, &stringValue); err != nil {
		return errors.New("invalid string value format")
	}

	// Check value type is int
	if intValue, err := strconv.Atoi(stringValue); err == nil {
		*ov = OperationValue(strconv.Itoa(intValue))
		return nil
	}

	// Check value type is float
	floatValue, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return nil
	}

	// Check if the number is an integer
	if math.Mod(floatValue, 1) != 0 {
		*ov = ""
	} else {
		*ov = OperationValue(strconv.Itoa(int(floatValue)))
	}

	return nil
}

func (oid *OperationID) UnmarshalJSON(data []byte) error {
	var rawID json.RawMessage
	if err := json.Unmarshal(data, &rawID); err != nil {
		return nil
	}

	// Attempted unmarshal into a number
	var intID int
	if err := json.Unmarshal(rawID, &intID); err == nil {
		oid.Value = intID
		return nil
	}

	// Attempted unmarshal into a string
	var stringID string
	if err := json.Unmarshal(rawID, &stringID); err == nil {
		oid.Value = stringID
		return nil
	}

	return nil
}

func (oca *OperationCreatedAt) UnmarshalJSON(data []byte) error {
	var value string
	_ = json.Unmarshal(data, &value)

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		*oca = ""
		return nil
	}

	*oca = OperationCreatedAt(t.Format(time.RFC3339))

	return nil
}

func processFile(filePath string) {
	records, err := readDataFromFile(filePath)
	if err != nil {
		fmt.Println("Error while reading data from a file: ", err)
		return
	}

	for i := 0; i < len(records); i++ {
		copyOperationToFinancialRecord(&records[i])
	}

	operationsByCompany := processOperations(records)
	sortInvalidOperations(operationsByCompany)
	result := createResultJSON(operationsByCompany)

	err = writeResultToFile(result, "out.json")
	if err != nil {
		fmt.Println("Error when writing to a file: ", err)
		return
	}

	fmt.Println("The data was successfully written to the file out.json")
}

func readDataFromFile(filePath string) ([]FinancialRecord, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var records []FinancialRecord
	err = json.Unmarshal(jsonData, &records)
	return records, err
}

func copyOperationToFinancialRecord(record *FinancialRecord) {
	if record.Operation != nil {
		if record.Type == "" {
			record.Type = record.Operation.Type
		}

		if record.Value == "" {
			record.Value = record.Operation.Value
		}

		if record.CreatedAt == "" {
			record.CreatedAt = record.Operation.CreatedAt
		}

		if record.ID.Value == nil {
			record.ID.Value = record.Operation.ID.Value
		}
	}
}

func processOperations(records []FinancialRecord) map[string]OperationResult {
	operationsByCompany := make(map[string]OperationResult)

	for _, record := range records {
		if record.ID.Value != nil && record.Company != "" && record.CreatedAt != "" {
			existingResult := operationsByCompany[record.Company]

			if record.Type == "" || record.Value == "" {
				existingResult.InvalidOperations = append(existingResult.InvalidOperations, DateWithID{record.ID.Value, record.CreatedAt})
			} else {
				existingResult.CountValidOperations += 1

				existingResult.Balance += addBalance(record.Type, record.Value)
			}
			operationsByCompany[record.Company] = existingResult
		}
	}

	return operationsByCompany
}

func sortInvalidOperations(operationsByCompany map[string]OperationResult) {
	for key := range operationsByCompany {
		sort.Slice(operationsByCompany[key].InvalidOperations, func(i, j int) bool {
			tmp1, _ := time.Parse(time.RFC3339, string(operationsByCompany[key].InvalidOperations[i].CreatedAt))
			tmp2, _ := time.Parse(time.RFC3339, string(operationsByCompany[key].InvalidOperations[j].CreatedAt))

			return tmp1.Before(tmp2)
		})
	}
}

func createResultJSON(operationsByCompany map[string]OperationResult) []ResultJSON {
	result := make([]ResultJSON, 0, len(operationsByCompany))
	for key := range operationsByCompany {
		var tmp ResultJSON
		tmp.Company = key
		tmp.CountValidOperations = operationsByCompany[key].CountValidOperations
		tmp.Balance = operationsByCompany[key].Balance
		tmp.InvalidOperationsID = make([]interface{}, len(operationsByCompany[key].InvalidOperations))
		for i, value := range operationsByCompany[key].InvalidOperations {
			tmp.InvalidOperationsID[i] = value.ID
		}
		result = append(result, tmp)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Company < result[j].Company
	})

	return result
}

func writeResultToFile(result []ResultJSON, fileName string) error {
	newFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	encoder := json.NewEncoder(newFile)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(result)
	return err
}

func addBalance(operationType OperationType, value OperationValue) int {
	sign := 1
	if operationType == "-" {
		sign = -1
	}
	result, _ := strconv.Atoi(string(value))
	return sign * result
}

func processFileByName(fileName string) error {
	file, err := os.Open(fileName)
	if err == nil {
		processFile(fileName)
		defer file.Close()
	}
	return err
}

func processCommandLineFlag() error {
	fileName := flag.String("file", "", "--file 'filename'")
	flag.Parse()

	return processFileByName(*fileName)
}

func processEnvironmentVariable() error {
	envVarValue := os.Getenv("FILE")
	if envVarValue != "" {
		fileName := filepath.Base(envVarValue)
		return processFileByName(fileName)
	}
	return errors.New("environment variable FILE is empty")
}

func processStdin() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path to the file:")
	scanner.Scan()
	inputPath := scanner.Text()

	if inputPath != "" {
		fileName := filepath.Base(inputPath)
		return processFileByName(fileName)
	}
	return errors.New("file not found")
}

func main() {
	if err := processCommandLineFlag(); err == nil {
		return
	}

	if err := processEnvironmentVariable(); err == nil {
		return
	}

	err := processStdin()
	if err == nil {
		return
	}

	fmt.Println(err)
}
