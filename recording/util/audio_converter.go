// recording/util/audio_converter.go

package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func AudioFileConvert(originalPath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(originalPath))

	if ext == ".wav" {
		return originalPath, nil
	}

	wavPath := strings.TrimSuffix(originalPath, ext) + ".wav"

	// run 'ffmpeg' command to convert
	cmd := exec.Command("ffmpeg", "-i", originalPath, "-ac", "1", "-ar", "44100", wavPath, "-y")
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to convert audio with ffmpeg: %v\nOutput: %s", err, string(cmdOutput))
	}

	_ = os.Remove(originalPath)

	return wavPath, nil
}