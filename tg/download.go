package tg

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const (
	url7tv    = "https://7tv.app/emotes/"
	gifSuffix = "/4x.gif"
	pngSuffix = "/4x.png"
)

// Use with defer os.Remove(tmpFile.Name())
func downloadImage(url string) (string, error) {
	// Create temporary file with "image-*.tmp" pattern
	tmpFile, err := os.CreateTemp("", "image-*.png")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close() // Close file when function exits

	// Fetch image from URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	// Copy image data to temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func convertPNG(inputPath string) (string, error) {
	tmpOutputFile, err := os.CreateTemp("", "output-*.png")
	if err != nil {
		return "", err
	}

	tmpOutputPath := tmpOutputFile.Name()
	tmpOutputFile.Close()

	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", inputPath,
		"-vf", `scale='if(gt(iw,ih),512,-1)':'if(gt(ih,iw),512,-1)'`,
		"-c:v", "png",
		tmpOutputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		os.Remove(tmpOutputPath)
		return "", err
	}

	return tmpOutputPath, nil
}

func convertToWebM(inputPath string) (string, error) {
	tmpOutputFile, err := os.CreateTemp("", "output-*.webm")
	if err != nil {
		return "", err
	}

	tmpOutputPath := tmpOutputFile.Name()
	tmpOutputFile.Close() // Close immediately so ffmpeg can write to it

	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", inputPath,
		"-vf", `scale='if(gt(iw,ih),512,-1)':'if(gt(ih,iw),512,-1)',format=yuva420p`,
		"-c:v", "libvpx-vp9",
		"-crf", "30",
		"-b:v", "0",
		"-quality", "good",
		"-cpu-used", "2",
		"-auto-alt-ref", "0",
		"-row-mt", "1",
		"-an",
		"-loop", "0",
		"-r", "30",
		tmpOutputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		os.Remove(tmpOutputPath)

		return "", err
	}

	return tmpOutputPath, nil
}
