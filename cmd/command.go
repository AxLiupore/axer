package main

import (
	"axer/cmd/main/cgroup"
	"axer/cmd/main/container"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroups limit ie: axer run -it [image] [command]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cs",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "c",
			Usage: "cpu limit",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		var cmdArray []string
		for _, arg := range ctx.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := ctx.Bool("it")
		limit := &cgroup.Limit{
			CPU:    ctx.String("c"),
			CPUSet: ctx.String("cs"),
			Memory: ctx.String("m"),
		}
		// run this cmd
		Run(tty, cmdArray, limit)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(ctx *cli.Context) error {
		logrus.Info("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}
