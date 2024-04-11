package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const downloadURL = "https://mirror.ghproxy.com/https://github.com/lqzhgood/wechat-need-web/releases/download/1.1.1/chrome.zip"
const baseDir = "WechatDesktop"

func main() {
	checkExtension()

	chromePath := checkConfig()

	// Open Chrome with the specified arguments
	executable, _ := os.Executable()
	execPath := filepath.Dir(executable)
	extensionPath := filepath.Join(execPath, baseDir, "chrome")
	args := []string{
		"--app=https://wx2.qq.com/?target=t",
		"-incognito",
		"--load-extension=" + extensionPath,
		"--start-maximized",
	}
	cmd := exec.Command(chromePath, args...)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	// Exit the program
	fmt.Println("Chrome started successfully.")
}
