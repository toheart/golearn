package container

import (
	"encoding/json"
	"fmt"
	errors2 "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"mydocker/constant"
	"os"
	"path"
	"time"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/21 9:51
@description:
**/

const (
	RUNNING       = "running"
	STOP          = "stopped"
	Exit          = "exited"
	InfoLoc       = "/var/lib/mydocker/containers/"
	InfoLocFormat = InfoLoc + "%s/"
	ConfigName    = "config.json"
	IDLength      = 10
	LogFile       = "%s-json.log"
)

type Info struct {
	Pid         int    `json:"pid"`        // 容器的init进程在宿主机上的 PID
	Id          string `json:"id"`         // 容器Id
	Name        string `json:"name"`       // 容器名
	Command     string `json:"command"`    // 容器内init运行命令
	CreatedTime string `json:"createTime"` // 创建时间
	Status      string `json:"status"`     // 容器的状态
}

func NewInfo(pid int, id string, name string, command string) *Info {
	// 预防错误, 未指定名称可以使用id作为名称
	if name == "" {
		name = id
	}

	info := &Info{
		Pid:         pid,
		Id:          id,
		Name:        name,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:      RUNNING,
		Command:     command}

	return info
}

func (i *Info) RecordInfo() error {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return errors2.WithMessage(err, "container info marshal failed")
	}
	jsonStr := string(jsonBytes)
	// 拼接容器存储信息的文件路径
	dirPath := fmt.Sprintf(InfoLocFormat, i.Id)
	if err := os.MkdirAll(dirPath, constant.Perm0622); err != nil {
		return errors2.WithMessagef(err, "mkdir %s failed", dirPath)
	}
	// 将容器信息写入文件
	fileName := path.Join(dirPath, ConfigName)
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return errors2.WithMessagef(err, "create file %s failed", fileName)
	}
	if _, err = file.WriteString(jsonStr); err != nil {
		return errors2.WithMessagef(err, "write container info to  file %s failed", fileName)
	}
	return nil
}

func DeleteContainerInfo(containerID string) error {
	dirPath := fmt.Sprintf(InfoLocFormat, containerID)
	if err := os.RemoveAll(dirPath); err != nil {
		log.Errorf("Remove dir %s error %s", dirPath, err)
		return errors2.WithMessagef(err, "remove dir %s failed", dirPath)
	}
	return nil
}

func GenerateContainerID() string {
	return randStringBytes(IDLength)
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// GetLogfile build logfile name by containerId
func GetLogfile(containerId string) string {
	return fmt.Sprintf(LogFile, containerId)
}
