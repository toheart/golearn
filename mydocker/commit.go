package main

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/utils"
	"os/exec"
)

/*
*
@file:
@author: levi.Tang
@time: 2024/5/21 9:31
@description:
*
*/
var ErrImageAlreadyExists = errors.New("Image Already Exists")

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit container to image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing image name")
		}
		containerId := context.Args().Get(0)
		imageName := context.Args().Get(1)
		return commitContainer(containerId, imageName)
	},
}

func commitContainer(containerId string, imageName string) error {
	mntPath := utils.GetMerged(containerId)
	imageTar := utils.GetImage(imageName)
	exists, err := utils.PathExists(imageTar)
	if err != nil {
		return errors.WithMessagef(err, "check is image [%s/%s] exist failed", imageName, imageTar)
	}
	if exists {
		return ErrImageAlreadyExists
	}
	log.Infof("commitContainer imagetar: %s", imageTar)

	if _, err := exec.Command("tar", "-zcf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		log.Errorf("tar folder %s error %s", mntPath, err)
	}
	return nil
}
