/*
Copyright © 2023 XieYanke

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/xieyanke/tcpdumpc/pkg/tcpdump"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print the tcpdumpc version",
	}

	app := &cli.App{
		Name:    "tcpdumpc",
		Usage:   "A wrapper of tcpdump for containers",
		Version: "0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Usage: "specify contianer id which to capture",
			},
			&cli.StringFlag{
				Name:  "runtime",
				Usage: "specify container runtime",
				Value: "docker",
				Action: func(ctx *cli.Context, v string) error {
					switch v {
					case "docker":
						return nil
					default:
						return fmt.Errorf("unsupport container runtime: %s", v)
					}
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "tcpdump",
				Aliases: []string{"dump"},
				Usage:   "execute tcpdump, eg: tcpdumpc --id XXXX dump -n -i eth0",
				Action: func(ctx *cli.Context) error {
					if !tcpdump.CheckTcpdumpExist() {
						return fmt.Errorf("'tcpdump' executable file not found in $PATH")
					}

					containerID := ctx.String("id")
					if containerID == "" {
						return fmt.Errorf("the container id, global flag '--id' must be specified")
					}

					containerRuntime := ctx.String("runtime")

					tcpdumpC := tcpdump.NewTcpdumpC(containerID, containerRuntime, os.Stdout, os.Stderr, ctx.Args().Slice())
					return tcpdumpC.Run()
				},
				SkipFlagParsing: true,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
