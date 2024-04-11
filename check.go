package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func checkExtension() {
	// Check if baseDir exists
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		fmt.Println("'WechatDesktop' directory does not exist. Creating it now.")
		err := os.Mkdir(baseDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	extensionDir := filepath.Join(baseDir, "chrome")
	// Check if 'chrome' folder exists
	if _, err := os.Stat(extensionDir); os.IsNotExist(err) {
		fmt.Println("'chrome' directory does not exist. Creating it now.")
		err := os.Mkdir(extensionDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Check if 'manifest.json' exists in the 'chrome' folder
	manifestPath := filepath.Join(extensionDir, "manifest.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		fmt.Println("'manifest.json' does not exist. Downloading and unzipping the file now.")
		err := downloadAndUnzip(downloadURL, baseDir)
		if err != nil {
			panic(err)
		}
	}
}

func checkConfig() string {
	// Check if 'config' file exists
	configFile := filepath.Join(baseDir, "config")
	var chromePath string
	if _, err := os.Stat(configFile); err == nil {
		// File exists, read the first line
		file, err := os.Open(configFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		chromePath, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		chromePath = strings.TrimSpace(chromePath)
	} else {
		// Prompt the user for Chrome path
		fmt.Print("Enter the path to Chrome: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		chromePath = scanner.Text()
		// If chromePath ends with '/' or '\', add "chrome.exe" to the end
		if strings.HasSuffix(chromePath, "/") || strings.HasSuffix(chromePath, "\\") {
			chromePath += "chrome.exe"
		}

		// Write the Chrome path to 'config' file in baseDir
		file, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(chromePath + "\n")
		if err != nil {
			panic(err)
		}
	}
	return chromePath
}
