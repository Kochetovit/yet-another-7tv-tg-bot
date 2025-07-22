package tg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	url7tv    = "https://cdn.7tv.app/emote/"
	gifSuffix = "/4x.gif"
	pngSuffix = "/4x.png"
)

// Use with defer os.Remove(tmpFile.Name())
func downloadImage(url string) (string, error) {
	// Create temporary file with "image-*.tmp" pattern
	tmpFile, err := os.CreateTemp("", "image-*.png")
	if err != nil {
		slog.Error("failed to create temp file", "err", err.Error())

		return "", err
	}
	//nolint:errcheck // don't care
	defer tmpFile.Close() // Close file when function exits

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("failed to create request", "err", err.Error())

		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	client := &http.Client{Timeout: 10 * time.Second}

	// Fetch image from URL
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed to fetch image new", "err", err.Error())

		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("failed to fetch image", "status", resp.StatusCode)

		return "", errors.New("wrong status code")
	}

	// Copy image data to temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		slog.Error("failed to copy image data", "err", err.Error())

		return "", err
	}

	return tmpFile.Name(), nil
}

func convertPNG(inputPath string) (string, error) {
	tmpOutputFile, err := os.CreateTemp("", "output-*.png")
	if err != nil {
		slog.Error("failed to create temp file", "err", err.Error())

		return "", err
	}

	tmpOutputPath := tmpOutputFile.Name()
	tmpOutputFile.Close()

	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", inputPath,
		"-vf", `scale='if(eq(iw,ih),512,if(gt(iw,ih),512,-2))':'if(eq(iw,ih),512,if(gt(ih,iw),512,-2))'`,
		"-c:v", "png",
		"-update", "1",
		tmpOutputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		os.Remove(tmpOutputPath)
		slog.Error("failed to convert image", "err", err.Error(), "stderr", stderr.String())

		return "", err
	}

	return tmpOutputPath, nil
}

func convertGIF(inputPath string) (string, error) {
	tmpOutputFile, err := os.CreateTemp("", "output-*.webm")
	if err != nil {
		slog.Error("failed to create temp file", "err", err.Error())

		return "", err
	}

	tmpOutputPath := tmpOutputFile.Name()
	tmpOutputFile.Close() // Close immediately so ffmpeg can write to it

	durationCmd := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-show_entries",
		"format=duration",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		inputPath,
	)

	var stderr bytes.Buffer
	durationCmd.Stderr = &stderr

	durationOut, err := durationCmd.Output()
	if err != nil {
		slog.Error("failed to get duration zz", "err", err.Error(), "stderr", stderr.String())

		return "", err
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(durationOut)), 64)
	if err != nil {
		slog.Error("failed to parse duration", "err", err.Error())

		return "", err
	}

	slog.Info("input duration", "duration", duration)

	speed := 1.0
	if duration > 3.0 {
		speed = duration / 3.0
	}

	vf := fmt.Sprintf("setpts=PTS/%.4f,fps=30,scale='if(eq(iw,ih),512,if(gt(iw,ih),512,-2))':'if(eq(iw,ih),512,if(gt(ih,iw),512,-2))',format=yuva420p", speed)

	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", inputPath,
		"-vf", vf,
		"-c:v", "libvpx-vp9",
		"-crf", "40", // Increase CRF to reduce file size
		"-b:v", "200K", // Lower bitrate to fit size constraint
		"-quality", "good",
		"-cpu-used", "2",
		"-auto-alt-ref", "0",
		"-row-mt", "1",
		"-an",
		"-loop", "0",
		"-r", "30", // Force output to max 30fps
		"-fs", "256K", // Max file size is 256kb.
		tmpOutputPath,
	)

	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		os.Remove(tmpOutputPath)
		slog.Error("failed to convert image", "err", err.Error(), "stderr", stderr.String())

		return "", err
	}

	return tmpOutputPath, nil
}
