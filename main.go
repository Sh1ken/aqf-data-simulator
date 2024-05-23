/*
 * 1. Loop through all files in data folder via config.json
 * 2. Make a copy of current file
 * 3. Delete first line
 * 4. Add a new line with randomized data (consider interval)
 * 5. Save the file and overwrite the original
 * 6. Sleep for a minute and do the process again
 */
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

func generateRandomRow(tempFileRoute string, file File) error {
	fmt.Println("| Generating random row for", file.Name, "in", tempFileRoute)
	/*
	 * 1. Open the file
	 * 2. Delete the first row
	 * 3. Read last row
	 * 4. Create new record based on the last row and config
	 * 5. Write the new record to the file
	 */
	return nil
}

func main() {
	// Retrieve config from config.json
	var config Config = retrieveConfig()

	// For each station
	for _, file := range config.Files {
		fmt.Println("| Processing", file.Name)

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

		// fmt.Println(file.Name, file.Separator, file.Interval, file.DateConfig.Format, file.DateConfig.Gmt)
		// for _, column := range file.Columns {
		// 	fmt.Println(column.Name, column.Type, column.Position)
		// }
	}
}
