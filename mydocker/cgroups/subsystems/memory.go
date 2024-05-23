package subsystems

import (
	"fmt"
	"github.com/pkg/errors"
	"mydocker/constant"
	"os"
	"path"
	"strconv"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/18 17:42
@description:
1. 比如限制内存则是往memory.limit_in_bytes里写入指定值
2. 添加某个进程到 cgroup 中就是往对应的 tasks 文件中写入对应的 pid
3. 删除 cgroup 就是把对应目录删掉
**/

type Memory struct {
}

func (m *Memory) Name() string {
	return "memory"
}

// Set 设置cgroupPath对应的cgroup的内存资源限制
func (m *Memory) Set(cgroupPath string, res *ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}

	subCgroupPath, err := getCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	// 将限制写入memory.limit_in_bytes 文件中
	if err := os.WriteFile(path.Join(subCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup memory fail %s", err)
	}
	return nil
}

// Apply 将pid加入到cgroupPath对应的cgroup中
func (m *Memory) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}
	subCgroupPath, err := getCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return errors.Wrapf(err, "get cgroup %s", cgroupPath)
	}
	if err := os.WriteFile(path.Join(subCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %s", err)
	}
	return nil
}

// Remove 删除cgroupPath对应的cgroup
func (m *Memory) Remove(cgroupPath string) error {
	subCgroupPath, err := getCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subCgroupPath)
}
