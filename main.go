package main

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Endpoint to upload a zip file
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Save the uploaded file
		zipPath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, zipPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Extract the uploaded zip file
		extractDir := "extracted_files/" + strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
		if err := extractZip(zipPath, extractDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded and extracted successfully!",
		})
	})

	r.Run(":8080")
}

// Function to extract a zip file
func extractZip(zipPath, extractDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()

		extractedFilePath := filepath.Join(extractDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, os.ModePerm)
		} else {
			if err := os.MkdirAll(filepath.Dir(extractedFilePath), os.ModePerm); err != nil {
				return err
			}

			extractedFile, err := os.Create(extractedFilePath)
			if err != nil {
				return err
			}
			defer extractedFile.Close()

			if _, err := io.Copy(extractedFile, zippedFile); err != nil {
				return err
			}
		}
	}

	return nil
}
