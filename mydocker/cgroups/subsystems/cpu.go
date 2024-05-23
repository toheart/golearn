package subsystems

import (
	"fmt"
	"mydocker/constant"
	"os"
	"path"
	"strconv"
)

/*
*
@file:
@author: levi.Tang
@time: 2024/5/18 19:44
@description:
*
*/
const (
	PeriodDefault = 100000
	Percent       = 100
)

type Cpu struct {
}

func (c *Cpu) Name() string {
	return "cpu"
}

func (c *Cpu) Set(cGroupPath string, res *ResourceConfig) error {
	if res.CpuCfsQuota == 0 && res.CpuShare == "" {
		return nil
	}
	subCgroupPath, err := getCgroupPath(c.Name(), cGroupPath, true)
	if err != nil {
		return err
	}
	// cpu.shares 控制的是CPU使用的比例而不是绝对值
	if res.CpuShare != "" {
		if err = os.WriteFile(path.Join(subCgroupPath, "cpu.shares"), []byte(res.CpuShare), constant.Perm0644); err != nil {
			return fmt.Errorf("set cgroup cpu share fail : %s ", err)
		}
	}
	// cpu.cfs_period_us & cpu.cfs_quota_us 控制的是CPU使用时间，单位是微秒，比如每1秒钟，这个进程只能使用200ms，相当于只能用20%的CPU
	if res.CpuCfsQuota != 0 {
		// cpu.cfs_period_us 默认为100000，即100ms
		if err = os.WriteFile(path.Join(subCgroupPath, "cpu.cfs_period_us"), []byte(strconv.Itoa(PeriodDefault)), constant.Perm0644); err != nil {
			return fmt.Errorf("set cgroup cpu cfs_period_us err: %s", err)
		}
		// cpu.cfs_quota_us 则根据用户传递的参数来控制，比如参数为20，就是限制为20%CPU，所以把cpu.cfs_quota_us设置为cpu.cfs_period_us的20%就行
		// 这里只是简单的计算了下，并没有处理一些特殊情况，比如负数什么的
		if err = os.WriteFile(path.Join(subCgroupPath, "cpu.cfs_quota_us"), []byte(strconv.Itoa(PeriodDefault)), constant.Perm0644); err != nil {
			return fmt.Errorf("set cgroup cpu cfs_quota_us err: %s", err)
		}
	}
	return nil
}

func (c *Cpu) Apply(cGroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuCfsQuota == 0 && res.CpuShare == "" {
		return nil
	}

	subCgroupPath, err := getCgroupPath(c.Name(), cGroupPath, false)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cGroupPath, err)
	}

	// 把pid写入进行就行
	if err = os.WriteFile(path.Join(subCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}
	return nil
}

func (c *Cpu) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(c.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}
