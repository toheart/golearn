package subsystems

/**
@file:
@author: levi.Tang
@time: 2024/5/18 17:09
@description:
**/

type ResourceConfig struct {
	MemoryLimit string
	CpuCfsQuota int
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	// Name 返回当前Subsystem的名称,比如cpu、memory
	Name() string
	// Set 设置某个cgroup在这个Subsystem中的资源限制
	Set(path string, res *ResourceConfig) error
	// Apply 将进程添加到某个cgroup中
	Apply(path string, pid int, res *ResourceConfig) error
	// Remove 移除某个Cgroup
	Remove(path string) error
}

// SubsystemsIns 通过不同的subsystem初始化实例创建资源限制处理链数组
var SubsystemsIns = []Subsystem{
	&CpuSet{},
	&Memory{},
	&Cpu{},
}
