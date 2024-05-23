package subsystems

import (
	"fmt"
	"github.com/pkg/errors"
	"mydocker/constant"
	"os"
	"path"
	"strconv"
)

/*
*
@file:
@author: levi.Tang
@time: 2024/5/19 8:52
@description:
*
*/
type CpuSet struct {
}

func (c *CpuSet) Name() string {
	return "cpuset"
}

func (c *CpuSet) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}
	subsysCgroupPath, err := getCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup cpuset fail %v", err)
	}
	return nil
}

func (c *CpuSet) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}
	subsysCgroupPath, err := getCgroupPath(c.Name(), cgroupPath, false)
	if err != nil {
		return errors.Wrapf(err, "get cgroup %c", cgroupPath)

	}
	if err := os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}
	return nil
}

func (c *CpuSet) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(c.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}
