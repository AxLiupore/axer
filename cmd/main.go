package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = "This axer is a simple container runtime implementation. The purpose of this project is to learn how to docker works and how to write a docker by myself. Enjoy it, just for fun."

func main() {
	app := cli.NewApp()
	app.Name = "axer"
	app.Usage = usage

	// define command
	app.Commands = []cli.Command{
		runCommand,
		initCommand,
	}

	// before app run this func
	app.Before = func(ctx *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
