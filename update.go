package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// getDownloadURL returns the os-specific url to download apollo from.
func getDownloadURL() string {
	// * determining the os-specific url

	switch runtime.GOOS {
	case "windows":
		return "https://github.com/Shravan-1908/apollo/releases/latest/download/apollo-windows-amd64.exe"
	case "linux":
		return "https://github.com/Shravan-1908/apollo/releases/latest/download/apollo-linux-amd64"
	case "darwin":
		return "https://github.com/Shravan-1908/apollo/releases/latest/download/apollo-darwin-amd64"
	default:
		return ""
	}
}

// download downloads the apollo executable from github and returns the path to the downloaded executable.
func download() string {
	fmt.Println("Downloading the apollo executable...")
	url := getDownloadURL()
	if url == "" {
		fmt.Println("Your OS isn't supported by apollo.")
		return ""
	}

	// * sending a request
	res, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: Unable to download the executable. Check your internet connection.")
		fmt.Println(err.Error())
		return ""
	}
	defer res.Body.Close()

	// * determining the executable path
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	downloadPath := filepath.Join(dir, "Downloads", "apollo_temp")
	if runtime.GOOS == "windows" {
		downloadPath += ".exe"
	}

	exe, er := os.Create(downloadPath)
	if er != nil {
		fmt.Println("Error: Unable to access file permissions.")
		fmt.Println(er.Error())
		return ""
	}
	defer exe.Close()

	bar := progressbar.DefaultBytes(res.ContentLength, "progress:")

	// * writing the received content to the apollo executable
	_, errr := io.Copy(io.MultiWriter(exe, bar), res.Body)
	if errr != nil {
		fmt.Println("Error: Unable to write the executable.")
		fmt.Println(errr.Error())
		return ""
	}

	return downloadPath
}

// Update updates apollo by downloading the latest executable from github, and renaming the
// old executable to `apollo-old` so that it can be deleted by `deletePreviousInstallation`.
func update() {
	fmt.Println("Updating apollo...")

	downloadPath := download()
	if downloadPath == "" {
		fmt.Println("Aborting update due to the above error.")
		return
	}

	originalPath := filepath.Join(getApolloDir(), "apollo")
	if runtime.GOOS == "windows" {
		originalPath += ".exe"
	}
	os.Rename(originalPath, originalPath+"-old")

	os.Chmod(downloadPath, 0700)
	os.Rename(downloadPath, originalPath)

	fmt.Println("apollo was updated successfully.")
}

// DeletePreviousInstallation deletes previous installation if it exists.
func deletePreviousInstallation() {
	apolloDir := getApolloDir()

	files, _ := ioutil.ReadDir(apolloDir)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "-old") {
			os.Remove(apolloDir + "/" + f.Name())
		}
	}
}
