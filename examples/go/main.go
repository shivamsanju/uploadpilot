package main

import (
	"fmt"
	"os"
	"test/uploader"
)

func main() {
	up := uploader.NewUploader(
		"a3660691-0412-4da7-a537-e63f500899e3",
		"a1e4d1f9-7d0a-4d85-8b25-ea8f90267dc7",
		"up-PoDTH92DEDKN78WD20250316165824",
		nil,
	)

	fileData, err := os.ReadFile("./files/error-dark.png")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	File := &uploader.File{
		Name:        "error-dark.png",
		Data:        fileData,
		ContentType: "image/png",
	}

	success, err := up.Upload(File, map[string]any{"description": "Test file"})
	if err != nil {
		fmt.Println("Upload failed:", err)
	} else {
		fmt.Println("Upload successful:", success)
	}
}
