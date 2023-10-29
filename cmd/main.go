package main

import (
	"axer/cmd/main/constant"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := entry.Level.String()
	msg := entry.Message

	logLine := "[" + timestamp + "] [" + level + "] " + msg + "\n"
	return []byte(logLine), nil
}

func main() {
	app := cli.NewApp()
	app.Name = constant.Name
	app.Usage = constant.Usage

	// define command
	app.Commands = []cli.Command{
		runCommand,
		initCommand,
		commitCommand,
		psCommand,
	}

	// before app run this func
	app.Before = func(ctx *cli.Context) error {
		logrus.SetFormatter(&CustomFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
