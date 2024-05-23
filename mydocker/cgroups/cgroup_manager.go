package cgroups

import (
	log "github.com/sirupsen/logrus"
	"mydocker/cgroups/subsystems"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/19 8:58
@description:
**/

type CgroupManager struct {
	Path string

	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{Path: path}
}

func (c *CgroupManager) Apply(pid int, res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		err := subSysIns.Apply(c.Path, pid, res)
		if err != nil {
			log.Errorf("apply subsystem:%s, err: %s", subSysIns.Name(), err)
		}
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		err := subSysIns.Set(c.Path, res)
		if err != nil {
			log.Errorf("set Subsystem: %s, err: %s", subSysIns.Name(), err)
		}
	}
	return nil
}

// Destroy 释放cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
