package subsystems

import (
	"bufio"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

/*
*
@file:
@author: levi.Tang
@time: 2024/5/18 17:23
@description:
*
*/
const mountPointIndex = 4

func getCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	// 不需要创建就直接返回
	cgroupRoot := findCgroupMountPoint(subsystem)
	absPath := path.Join(cgroupRoot, cgroupPath)
	if !autoCreate {
		return absPath, nil
	}
	// 指定自动创建
	_, err := os.Stat(absPath)
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(absPath, os.ModePerm)
		return absPath, err
	}
	// 其他错误或者没有错误都直接返回，如果err=nil,那么errors.Wrap(err, "")也会是nil
	return absPath, errors.Wrap(err, "create cgroup")
}

func findCgroupMountPoint(subsystem string) string {
	// /proc/self/mountinfo 为当前进程的Mountinfo信息
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// 举例 36 26 0:30 / /sys/fs/cgroup rw,nosuid,nodev,noexec,relatime - cgroup2 cgroup2 rw,nsdelegate
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		subsystems := strings.Split(fields[len(fields)-1], ",")
		for _, opt := range subsystems {
			if opt == subsystem {
				// 找到文件路径
				return fields[mountPointIndex]
			}
		}
	}
	if err = scanner.Err(); err != nil {
		log.Errorf("read err: %s", err)
		return ""
	}
	return ""
}
