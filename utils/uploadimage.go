package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UploadImage(r *http.Request, fieldName string, oldImagePath string) (string, error) {
	file, fileHeader, err := r.FormFile(fieldName)
	fmt.Println(">> r.FormFile() dipanggil")

	if err != nil {
		fmt.Println(">> error FormFile:", err)
		if err == http.ErrMissingFile {
			fmt.Println(">> Tidak ada gambar di-upload")
			return oldImagePath, nil // Gak upload gambar baru, pakai gambar lama
		}
		return "", err
	}

	fmt.Printf(">> fileHeader.Filename: %s\n", fileHeader.Filename)

	defer file.Close()

	if oldImagePath != "" && !strings.Contains(oldImagePath, "default") {
		oldFile := "" + oldImagePath
		_ = os.Remove(oldFile)
	}

	ext := filepath.Ext(fileHeader.Filename)
	newFileName := time.Now().Format("20060102150405") + ext
	savePath := "./public/assets/img/" + newFileName

	if ext == "" {
		ext = ".png" // atau .png default aman
	}

	fmt.Println(">> ext:", ext)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}

	fmt.Println(">> dst:", dst)
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/assets/img/" + newFileName, nil

}
