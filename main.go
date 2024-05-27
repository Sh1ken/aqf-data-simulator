package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func retrieveConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return Config{}
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	return config
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func deleteFirstLine(tempFileRoute string) error {
	cmd := exec.Command("sed", "-i", "-e", "1d", tempFileRoute)

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func getLastRow(tempFileRoute string) (string, error) {
	cmd := exec.Command("tail", "-n", "1", tempFileRoute)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func appendToFile(filename, text string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' >> %s", text, filename))
	return cmd.Run()
}

func removeNewlines(text string) string {
	processedText := strings.ReplaceAll(text, "\r", "")
	return strings.ReplaceAll(processedText, "\n", "")
}

func addIntervalToCurrentDate(currentDateString string, dateConfig DateConfig) (string, error) {
	// Convert the last value to a date
	currentDate, err := time.Parse(dateConfig.Format, strings.ReplaceAll(currentDateString, "\"", ""))
	if err != nil {
		return "", err
	}

	// Add the interval to the date
	currentDate = currentDate.Add(time.Minute * time.Duration(dateConfig.Interval))

	return "\"" + currentDate.Format(dateConfig.Format) + "\"", nil
}

func generateRandomDate(columnName string, lastValue string, dateConfig DateConfig) (string, error) {
	lastValueDate := ""
	err := error(nil)

	switch columnName {
	case "timestamp":
		lastValueDate, err = addIntervalToCurrentDate(lastValue, dateConfig)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("Unknown column '%s'", columnName)
	}

	return lastValueDate, nil
}

func generateRandomInt(columnName string, lastValue string) (string, error) {
	switch columnName {
	case "id":
		lastValueInt, err := strconv.Atoi(lastValue)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(lastValueInt + 1), nil
	default:
		return "", fmt.Errorf("Unknown column '%s'", columnName)
	}
}

func generateBiasedFloat(lastValueFloat float64, maxBiasPercentage float64) float64 {
	// Generate a random float between -maxBiasPercentage and maxBiasPercentage
	bias := (rand.Float64() * maxBiasPercentage * 2) - maxBiasPercentage

	// Add the bias to the last value
	valueVariation := lastValueFloat * (bias / 100)
	lastValueFloat += valueVariation

	// If the value is negative, return 0
	if lastValueFloat < 0 {
		return 0
	}

	// Return the new value
	return lastValueFloat
}

func generateRandomFloat(columnName string, lastValue string) (string, error) {
	lastValueFloat, err := strconv.ParseFloat(lastValue, 64)

	if err != nil {
		return "", err
	}

	if lastValueFloat < 0 {
		return "", fmt.Errorf("Negative value for column '%s'", columnName)
	}

	// 70% of the time, the value will still be 0
	if lastValueFloat == 0 {
		if rand.Float64() < 0.7 {
			lastValueFloat = 0
		} else {
			lastValueFloat = 1 + 4*rand.Float64()
		}

	}

	valuePercentageVariation := 0.0

	switch columnName {
	case "battery":
		valuePercentageVariation = 10.0
	case "temp":
		valuePercentageVariation = 10.0
	case "level":
		valuePercentageVariation = 10.0
	case "rain":
		valuePercentageVariation = 15.0
	case "turbidity":
		valuePercentageVariation = 20.0
	case "caudal_ls":
		valuePercentageVariation = 5.0
	case "velocity":
		valuePercentageVariation = 5.0
	default:
		return "", fmt.Errorf("Unknown column '%s'", columnName)
	}

	biasedValue := generateBiasedFloat(lastValueFloat, valuePercentageVariation)

	return fmt.Sprintf("%.2f", biasedValue), nil
}

func generateRandomRow(tempFileRoute string, file File) error {
	// First, we delete the first line of the file
	err := deleteFirstLine(tempFileRoute)
	if err != nil {
		return err
	}

	// Secondly, we get the last line of the file
	lastRow, err := getLastRow(tempFileRoute)
	if err != nil {
		return err
	}

	// Remove any new lines from the last row
	lastRow = removeNewlines(lastRow)
	fmt.Println("| Last row:", lastRow)

	// Then we generate data based on the last row that we just got
	lastRowColumns := strings.Split(lastRow, file.Separator)
	newRowData := []string{}
	currentColumnContent := ""

	for index, column := range file.Columns {
		fmt.Println("| Generating data for", column.Name)
		fmt.Println("| Old:", lastRowColumns[index])

		switch column.Type {
		case "datetime":
			currentColumnContent, err = generateRandomDate(column.Name, lastRowColumns[index], file.DateConfig)
		case "int":
			currentColumnContent, err = generateRandomInt(column.Name, lastRowColumns[index])
		case "float":
			currentColumnContent, err = generateRandomFloat(column.Name, lastRowColumns[index])
		default:
			return fmt.Errorf("Unknown type '%s'", column.Type)
		}

		if err != nil {
			return err
		}

		fmt.Println("| New:", currentColumnContent)

		newRowData = append(newRowData, currentColumnContent)
	}

	// Finally, we append a new line with randomized data
	newRow := strings.Join(newRowData, file.Separator)
	fmt.Println("| New row:", newRow)
	err = appendToFile(tempFileRoute, newRow)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Retrieve config from config.json
	var config Config = retrieveConfig()

	// For each station
	for _, file := range config.Files {
		fmt.Println("|----")
		fmt.Println("| Initializing", file.Name)

		// Create a temporary copy of the file
		originalFileRoute := config.DataFolder + file.Name
		tempFileRoute := config.TempFolder + "temp_" + file.Name

		err := copyFile(originalFileRoute, tempFileRoute)

		if err != nil {
			fmt.Println(err)
			return
		}

		// Generate a random in the tempFile
		err = generateRandomRow(tempFileRoute, file)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Copy the temp file to the original file
		err = copyFile(tempFileRoute, originalFileRoute)

		if err != nil {
			fmt.Println(err)
			return
		}

		// And finally remove the temp file after the process
		err = os.Remove(tempFileRoute)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("| Closing", file.Name)
	}

	fmt.Println("|----")
}
