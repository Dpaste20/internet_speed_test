package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const downloadURL = "http://speedtest.tele2.net/10MB.zip"
const tempFile = "temp_download.bin"

func downloadFile(url, filepath string) (float64, error) {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return 0, err
	}

	elapsed := time.Since(start).Seconds()
	mbps := (float64(written) * 8 / 1_000_000) / elapsed

	return mbps, nil
}

func uploadTest() (float64, error) {
	data := make([]byte, 5*1024*1024)

	start := time.Now()
	resp, err := http.Post("https://httpbin.org/post", "application/octet-stream",
		io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	elapsed := time.Since(start).Seconds()
	mbps := (float64(len(data)) * 8 / 1_000_000) / elapsed

	return mbps, nil
}

func main() {
	fmt.Println("Running Speed Test...\n")

	// DOWNLOAD TEST → saves file
	fmt.Print("Testing Download Speed : ")
	downloadMbps, err := downloadFile(downloadURL, tempFile)
	if err != nil {
		fmt.Println("failed:", err)
	} else {
		fmt.Printf("%.2f Mbps\n", downloadMbps)
	}

	os.Remove(tempFile)

	fmt.Print("Testing Upload Speed : ")
	uploadMbps, err := uploadTest()
	if err != nil {
		fmt.Println("failed:", err)
	} else {
		fmt.Printf("%.2f Mbps\n", uploadMbps)
	}

	fmt.Println("\nSpeed Test Complete. Temporary file deleted.")
}
