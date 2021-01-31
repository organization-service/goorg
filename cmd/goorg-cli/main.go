package main

import (
	"log"
	"os"

	"github.com/organization-service/goorg"
	generator "github.com/organization-service/goorg/cmd/goorg-cli/di-generator"
	"github.com/urfave/cli/v2"
)

func diCommand() *cli.Command {
	const (
		inputDir  = "input"
		outputDir = "output"
	)
	return &cli.Command{
		Name:  "di",
		Usage: "dependency injection command",
		Subcommands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"gene"},
				Usage:   "Generate di register",
				Action: func(c *cli.Context) error {
					return generator.Action(c.String(inputDir), c.String(outputDir))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    inputDir,
						Aliases: []string{"i"},
						Value:   "./",
						Usage:   "解析したいディレクトリ",
					},
					&cli.StringFlag{
						Name:    outputDir,
						Aliases: []string{"o"},
						Value:   "./registry",
						Usage:   "生成されるディレクトリ",
					},
				},
			},
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "goorg-cli"
	app.Version = goorg.Version
	app.Usage = "goorg framework command"
	app.Commands = []*cli.Command{
		diCommand(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
