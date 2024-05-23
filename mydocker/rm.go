package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/container"
)

/*
*
@file:
@author: levi.Tang
@time: 2024/5/21 21:56
@description:
*
*/
var removeCommand = cli.Command{
	Name:  "rm",
	Usage: "remove unused containers,e.g. mydocker rm 1234567890",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "f", // 强制删除
			Usage: "force delete running container,",
		}},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}
		containerId := context.Args().Get(0)
		force := context.Bool("f")
		removeContainer(containerId, force)
		return nil
	},
}

func removeContainer(containerId string, force bool) {
	containerInfo, err := getInfoByContainerId(containerId)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerId, err)
		return
	}

	switch containerInfo.Status {
	case container.STOP: // STOP 状态容器直接删除即可
		if err = container.DeleteContainerInfo(containerId); err != nil {
			log.Errorf("Remove file %s error %v", containerId, err)
			return
		}
		container.DeleteWorkSpace(containerId, containerInfo.Name)
	case container.RUNNING: // RUNNING 状态容器如果指定了 force 则先 stop 然后再删除
		if !force {
			log.Errorf("Couldn't remove running container [%s], Stop the container before attempting removal or"+
				" force remove", containerId)
			return
		}
		log.Infof("force delete running container [%s]", containerId)
		stopContainer(containerId)
		removeContainer(containerId, force)
	default:
		log.Errorf("Couldn't remove container,invalid status %s", containerInfo.Status)
		return
	}
}
