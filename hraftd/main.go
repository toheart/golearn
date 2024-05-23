package main

/**
@file:
@author: levi.Tang
@time: 2024/7/31 15:35
@description:
**/

import (
	"github.com/hashicorp/raft"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	// 1. 设置raft的configure
	config := raft.DefaultConfig()
	config.LocalID = "node01"
	// 2. 设置raft的通信方式
	addr, err := net.ResolveIPAddr("tcp", ":10001")
	if err != nil {
		panic(err)
	}
	transport, err := raft.NewTCPTransport(":10001", addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		panic(err)
	}
	// 设置快照
	snapshots, err := raft.NewFileSnapshotStore("./node01", 2, os.Stderr)
	if err != nil {
		panic(err)
	}
	var logStore raft.LogStore
	var stableStore raft.StableStore
	// 设置log存储, 测试使用内存方式
	logStore = raft.NewInmemStore()
	// 设置stable存储, 测试使用内存方式
	stableStore = raft.NewInmemStore()
	// 初始化系统
	fms := Fsm{}
	ra, err := raft.NewRaft(config, fms, logStore, stableStore, snapshots, transport)
	if err != nil {
		panic(err)
	}
	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	}
	ra.BootstrapCluster(configuration)
}

type Fsm struct {
}

func (f Fsm) Apply(log *raft.Log) interface{} {
	return nil
}

func (f Fsm) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (f Fsm) Restore(snapshot io.ReadCloser) error {
	return nil
}
