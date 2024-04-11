package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Downloads and unzips a file from a URL to a specified destination.
func downloadAndUnzip(url, dest string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create a temporary file to write the zip content
	tmpFile, err := os.CreateTemp("", "download-*.zip")
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name()) // Clean up the file afterwards

	// Write the body to the file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return err
	}

	// Unzip the file to the destination
	_, err = tmpFile.Seek(0, 0) // Rewind the file
	if err != nil {
		return err
	}
	err = unzip(tmpFile, dest)
	return err
}

// Unzips a file to a specified destination.
func unzip(src *os.File, dest string) error {
	fi, err := src.Stat()
	if err != nil {
		return err
	}

	r, err := zip.NewReader(src, fi.Size())
	if err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to handle error
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
