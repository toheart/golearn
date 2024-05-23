package file

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

/**
@file:
@author: levi.Tang
@time: 2024/4/9 16:22
@description:
**/

func TestFileInfo(t *testing.T) {
	entries, err := os.ReadDir("D:\\tests\\G")
	if err != nil {
		t.Errorf("读取目录内容 %s, error: %s", "D:\\tests\\G", err)
	}

	// 遍历目录内容，筛选出子目录
	for _, entry := range entries {
		if entry.IsDir() {
			_, err := strconv.Atoi(entry.Name())
			if err != nil {
				continue
			}
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			fmt.Println(info.Name(), "is a symbolic link")
		}
	}
}
