package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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
	return strings.ReplaceAll(text, "\n", "")
}

// TODO: random data generation
func generateRandomDate(columnName string, lastValue string) string {
	return lastValue
}

// TODO: random data generation
func generateRandomInt(columnName string, lastValue string) string {
	return lastValue
}

// TODO: random data generation
func generateRandomFloat(columnName string, lastValue string) string {
	return lastValue
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

	for index, column := range file.Columns {
		fmt.Println("| Generating data for", column.Name)
		switch column.Type {
		case "datetime":
			newRowData = append(newRowData, generateRandomDate(column.Name, lastRowColumns[index]))
		case "int":
			newRowData = append(newRowData, generateRandomInt(column.Name, lastRowColumns[index]))
		case "float":
			newRowData = append(newRowData, generateRandomFloat(column.Name, lastRowColumns[index]))
		default:
			return fmt.Errorf("Unknown type '%s'", column.Type)
		}
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
