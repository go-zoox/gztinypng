package main

import (
	"fmt"
	"os"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/fs/type/yaml"
	"github.com/go-zoox/gztinypng/tinypng"
)

var configDir = fs.JoinHomeDir(".config/go-zoox")
var tinypngConfigPath = fs.JoinPath(configDir, "gztinypng.yml")

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
		Name:    "gztinypng",
		Usage:   "TinyPNG cli, easy way to compress images",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Usage:    "input image file",
				Aliases:  []string{"i"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Usage:    "output image file",
				Aliases:  []string{"o"},
				Required: true,
			},
			&cli.StringFlag{
				Name:  "api-key",
				Usage: "tinypng api key",
			},
		},
	})

	app.Command(func(ctx *cli.Context) (err error) {
		cfg := tinypng.Config{}

		if !fs.IsExist(configDir) {
			if err := fs.Mkdirp(configDir); err != nil {
				return fmt.Errorf("failed to create config directory: %v", err)
			}
		}

		if fs.IsExist(tinypngConfigPath) {
			if err := yaml.Read(tinypngConfigPath, &cfg); err != nil {
				return fmt.Errorf("failed to read config file: %v", err)
			}
		}

		if ctx.String("api-key") != "" {
			cfg.ApiKey = ctx.String("api-key")
		}

		input := ctx.String("input")
		output := ctx.String("output")

		cfg.InputFile, err = fs.Open(input)
		if err != nil {
			return fmt.Errorf("failed to open input image: %v", err)
		}
		defer cfg.InputFile.Close()

		cfg.OutputFile, err = os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to open output image: %v", err)
		}
		defer cfg.OutputFile.Close()

		if cfg.ApiKey == "" {
			return fmt.Errorf("api key is required")
		}

		return tinypng.TinyPNG(&cfg)
	})

	app.Run()
}
