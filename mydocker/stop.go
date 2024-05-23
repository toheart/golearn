package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/container"
	"os"
	"path"
	"syscall"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/21 21:31
@description:
**/

var stopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container,e.g. mydocker stop 1234567890",
	Action: func(context *cli.Context) error {
		// 期望输入是：mydocker stop 容器Id，如果没有指定参数直接打印错误
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}
		containerId := context.Args().Get(0)
		stopContainer(containerId)
		return nil
	},
}

func stopContainer(containerId string) {
	// 1. 根据容器Id查询容器信息
	containerInfo, err := getInfoByContainerId(containerId)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerId, err)
		return
	}
	// 2.发送SIGTERM信号
	if err = syscall.Kill(containerInfo.Pid, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s error %v", containerId, err)
		return
	}
	// 3.修改容器信息，将容器置为STOP状态，并清空PID
	containerInfo.Status = container.STOP
	containerInfo.Pid = 0
	if err := containerInfo.RecordInfo(); err != nil {
		log.Errorf("RecordInfo %s error %s", containerId, err)
		return
	}
}

func getInfoByContainerId(containerId string) (*container.Info, error) {
	dirPath := fmt.Sprintf(container.InfoLocFormat, containerId)
	configFilePath := path.Join(dirPath, container.ConfigName)
	contentBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "read file %s", configFilePath)
	}
	var containerInfo container.Info
	if err = json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return nil, err
	}
	return &containerInfo, nil
}
