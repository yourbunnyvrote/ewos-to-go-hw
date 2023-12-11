package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type PreResult struct {
	Company              string
	CountValidOperations int
	Balance              int
	InvalidOperations    []OperationsWithDate
}

type Result struct {
	Company              string        `json:"company"`
	CountValidOperations int           `json:"valid_operations_count"`
	Balance              int           `json:"balance"`
	InvalidOperations    []interface{} `json:"invalid_operations,omitempty"`
}

type OperationsWithDate struct {
	ID        interface{}
	CreatedAt OperationCreatedAt
}

type Operation struct {
	Type      OperationType      `json:"type,omitempty"`
	Value     OperationValue     `json:"value,omitempty"`
	ID        OperationID        `json:"id,omitempty"`
	CreatedAt OperationCreatedAt `json:"created_at,omitempty"`
}

type FinancialRecord struct {
	Company   string             `json:"company,omitempty"`
	Operation *Operation         `json:"operation,omitempty"`
	Type      OperationType      `json:"type,omitempty"`
	Value     OperationValue     `json:"value,omitempty"`
	ID        OperationID        `json:"id,omitempty"`
	CreatedAt OperationCreatedAt `json:"created_at,omitempty"`
}

type OperationType string

type OperationValue string

type OperationID struct {
	ID interface{}
}

type OperationCreatedAt string

const (
	Plus  OperationType = "+"
	Minus OperationType = "-"
)

func (oca *OperationCreatedAt) UnmarshalJSON(data []byte) error {
	var value string
	_ = json.Unmarshal(data, &value)

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		*oca = ""
		return nil
	} else {
		*oca = OperationCreatedAt(t.Format(time.RFC3339))
	}
	return nil
}

func (ot *OperationType) UnmarshalJSON(data []byte) error {
	var value string
	_ = json.Unmarshal(data, &value)

	switch value {
	case "income", "+":
		*ot = Plus
		break
	case "outcome", "-":
		*ot = Minus
		break
	}
	return nil
}

func (ov *OperationValue) UnmarshalJSON(data []byte) error {
	var rawID json.RawMessage
	if err := json.Unmarshal(data, &rawID); err != nil {
		*ov = ""
		return nil
	}

	var numericValue int
	if err := json.Unmarshal(rawID, &numericValue); err == nil {
		*ov = OperationValue(strconv.Itoa(numericValue))
		return nil
	}

	var floatValue float64
	if err := json.Unmarshal(rawID, &floatValue); err == nil {
		if math.Mod(floatValue, 1) != 0 {
			*ov = ""
		} else {
			*ov = OperationValue(strconv.Itoa(int(floatValue)))
		}
		return nil
	}

	var stringID string
	if err := json.Unmarshal(rawID, &stringID); err == nil {
		if numericValue, err = strconv.Atoi(stringID); err == nil {
			*ov = OperationValue(strconv.Itoa(numericValue))
			return nil
		}
		floatValue, err = strconv.ParseFloat(stringID, 64)
		if err != nil {
			*ov = ""
			return nil
		}
		if math.Mod(floatValue, 1) != 0 {
			*ov = ""
		} else {
			*ov = OperationValue(strconv.Itoa(int(floatValue)))
		}
		return nil
	}
	return nil
}

func (oid *OperationID) UnmarshalJSON(data []byte) error {
	var rawID json.RawMessage
	if err := json.Unmarshal(data, &rawID); err != nil {
		oid.ID = ""
		return nil
	}

	var numericID int
	if err := json.Unmarshal(rawID, &numericID); err == nil {
		oid.ID = numericID
		return nil
	}

	var stringID string
	if err := json.Unmarshal(rawID, &stringID); err == nil {
		oid.ID = stringID
		return nil
	}

	oid.ID = ""
	return nil
}

func processFile(filePath string) {
	jsonData, err := os.ReadFile(filePath)
	var records []FinancialRecord
	err = json.Unmarshal([]byte(jsonData), &records)

	for i := 0; i < len(records); i++ {
		if records[i].Operation != nil {
			if records[i].Type == "" {
				records[i].Type = records[i].Operation.Type
			}

			if records[i].Value == "" {
				records[i].Value = records[i].Operation.Value
			}

			if records[i].CreatedAt == "" {
				records[i].CreatedAt = records[i].Operation.CreatedAt
			}

			if records[i].ID.ID == nil {
				records[i].ID.ID = records[i].Operation.ID.ID
			}
		}
	}

	validOperations := make(map[string]PreResult)

	for _, record := range records {
		if record.ID.ID != nil && record.Company != "" && record.CreatedAt != "" {
			existingResult := validOperations[record.Company]
			existingResult.Company = record.Company

			if record.Type == "" || record.Value == "" {
				existingResult.InvalidOperations = append(existingResult.InvalidOperations, OperationsWithDate{record.ID.ID, record.CreatedAt})
			} else {
				existingResult.CountValidOperations += 1

				existingResult.Balance += addBalance(record.Type, record.Value)
			}
			validOperations[record.Company] = existingResult
		}
	}

	for _, key := range validOperations {
		sort.Slice(validOperations[key.Company].InvalidOperations, func(i, j int) bool {
			tmp1, _ := time.Parse(time.RFC3339, string(validOperations[key.Company].InvalidOperations[i].CreatedAt))
			tmp2, _ := time.Parse(time.RFC3339, string(validOperations[key.Company].InvalidOperations[j].CreatedAt))

			return tmp1.Before(tmp2)
		})
	}

	result := make([]Result, 0, len(validOperations))
	for _, val := range validOperations {
		var tmp Result
		tmp.Company = val.Company
		tmp.CountValidOperations = val.CountValidOperations
		tmp.Balance = val.Balance
		tmp.InvalidOperations = make([]interface{}, len(val.InvalidOperations))
		for i, value := range val.InvalidOperations {
			tmp.InvalidOperations[i] = value.ID
		}
		result = append(result, tmp)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Company < result[j].Company
	})

	newFile, err := os.Create("out.json")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer newFile.Close()

	encoder := json.NewEncoder(newFile)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(result)
	if err != nil {
		fmt.Println("Ошибка при сериализации и записи в файл:", err)
		return
	}

	fmt.Println("Данные успешно записаны в файл output.json")
}

func addBalance(operationType OperationType, value OperationValue) int {
	sign := 1
	if operationType == "-" {
		sign = -1
	}
	result, _ := strconv.Atoi(string(value))
	return sign * result
}

func main() {
	fileName := *flag.String("file", "", "--file 'filename'")
	flag.Parse()

	fmt.Println(fileName)

	file, err := os.Open(fileName)
	if err == nil {
		processFile(fileName)
		defer file.Close()
		return
	}

	envVarValue := os.Getenv("FILE")
	if envVarValue != "" {
		fileName = filepath.Base(envVarValue)
		file, err = os.Open(fileName)
		if err == nil {
			processFile(fileName)
			defer file.Close()
			return
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите путь к файлу: ")
	scanner.Scan()
	inputPath := scanner.Text()

	if inputPath != "" {
		fileName = filepath.Base(inputPath)
		file, err = os.Open(fileName)
		if err == nil {
			processFile(fileName)
			defer file.Close()
			return
		}
	}
	fmt.Println(err)

	return
}
