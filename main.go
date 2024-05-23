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

func main() {

	// Retrieve config from config.json
	var config Config = retrieveConfig()

	// For each station
	for _, file := range config.Files {
		fmt.Println("| Processing", file.Name)

		// We make a temporary copy of the file
		currentFile := config.FolderData + file.Name
		tempFile := config.FolderTemp + "temp_" + file.Name

		err := copyFile(currentFile, tempFile)

		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Remove(tempFile)
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
