package main

import (
	log "github.com/sirupsen/logrus"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strings"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/17 21:08
@description:
**/

func Run(tty bool, cmdArr []string, res *subsystems.ResourceConfig, volume string, name string, imageName string) {
	containerId := container.GenerateContainerID()
	parent, writePipe := container.NewParentProcess(tty, volume, containerId, imageName)
	if parent == nil {
		log.Error("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	// record info
	info := container.NewInfo(parent.Process.Pid, containerId, name, strings.Join(cmdArr, ""))
	if err := info.RecordInfo(); err != nil {
		log.Errorf("Record container info error %s", err)
		return
	}
	// 创建cgroup manager
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	_ = cgroupManager.Set(res)
	_ = cgroupManager.Apply(parent.Process.Pid, res)
	sendInitCommand(cmdArr, writePipe)
	if tty {
		_ = parent.Wait()
		container.DeleteWorkSpace("/root/", volume)
		container.DeleteContainerInfo(containerId)
	}
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s ", command)

	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
