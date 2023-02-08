package tinypng

import (
	"fmt"
	"io"
	"os"

	"github.com/go-zoox/fetch"
	"github.com/go-zoox/headers"
)

// document: https://tinypng.com/developers
const (
	user = "api"
	url  = "https://api.tinypng.com/shrink"
)

type Config struct {
	ApiKey     string `json:"api_key"`
	InputFile  *os.File
	OutputFile *os.File
}

func TinyPNG(cfg *Config) error {
	response, err := fetch.New().
		SetMethod(fetch.POST).
		SetURL(url).
		SetHeader(headers.ContentType, "application/octet-stream").
		SetBody(cfg.InputFile).
		SetBasicAuth(user, cfg.ApiKey).
		Execute()
	if err != nil {
		return err
	}

	if !response.Ok() {
		return fmt.Errorf("failed to tinypng(1): %s", response.String())
	}

	// 302
	compressedImageURL := response.Headers.Get("location")
	if compressedImageURL == "" {
		return fmt.Errorf("failed to tinypng(2): cannot get location")
	}
	response, err = fetch.Stream(compressedImageURL)
	if err != nil {
		return err
	}

	if !response.Ok() {
		return fmt.Errorf("failed to tinypng(3): %s", response.String())
	}

	if _, err := io.Copy(cfg.OutputFile, response.Stream); err != nil {
		return err
	}

	return nil
}
