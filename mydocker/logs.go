package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"io"
	"mydocker/container"
	"os"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/21 10:28
@description:
**/

var logCommand = cli.Command{
	Name:  "logs",
	Usage: "print logs of a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("please input your container id")
		}
		containerName := context.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}

func logContainer(containerId string) {
	logFileLocation := fmt.Sprintf(container.InfoLocFormat, containerId) + container.GetLogfile(containerId)
	file, err := os.Open(logFileLocation)
	defer file.Close()
	if err != nil {
		log.Errorf("Log container open file %s error %v", logFileLocation, err)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		log.Errorf("Log container read file %s error %v", logFileLocation, err)
		return
	}
	_, err = fmt.Fprint(os.Stdout, string(content))
	if err != nil {
		log.Errorf("Log container Fprint  error %v", err)
		return
	}
}
