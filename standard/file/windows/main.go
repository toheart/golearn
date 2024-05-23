package main

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

/**
@file:
@author: levi.Tang
@time: 2024/8/6 17:20
@description:
**/

func main() {

	localGameDir := "G:\\9105"
	srcDir := "H:\\9105"
	var err error
	// 判断原始目录已存在
	if !PathIsDir(srcDir) {
		fmt.Errorf("dir not found: %s \n", srcDir)
		return
	}
	cmd := exec.Command("powershell", "robocopy", "\""+srcDir+"\"", "\""+localGameDir+"\"", "/e", "/purge", "/R:0", "/mt:12")
	defer func() {
		if err == nil {
			fmt.Printf("powershell robocopy %s %s /e /purge /R:0 /mt:12 \n", srcDir, localGameDir)
		}
	}()
	var sin, sout, serr bytes.Buffer
	cmd.Stdin = &sin
	cmd.Stdout = &sout
	cmd.Stderr = &serr
	if err := cmd.Run(); err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		fmt.Printf("exitCode: %d \n", exitCode)
		fmt.Printf("err: %s", err)
		// robocopy 任何等于或大于 8 的值表示复制操作期间至少发生了一次失败。
		if exitCode < 8 {
			return
		}
		fmt.Errorf("sin: %s ||sout: %s || serr: %s", gbk_to_utf8(sin.Bytes()),
			gbk_to_utf8(sout.Bytes()),
			gbk_to_utf8(serr.Bytes()))
	}
	return
}

func gbk_to_utf8(s []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	if d, e := ioutil.ReadAll(reader); e == nil {
		return d
	}
	return s
}

func PathIsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if fileInfo.IsDir() {
		return true
	}
	return false
}

// GetDiskFreeSpace
//
//	@Description: 获取输出的磁盘的空间大小
//	@param disk	 盘符名称
//	@return int	 返回磁盘大小
//	@return error
func GetDiskFreeSpace(disk string) (float64, error) {
	regex := regexp.MustCompile(`^[A-Za-z]:\\{0,2}`)
	if !regex.MatchString(disk) {
		return 0, fmt.Errorf("please input disk path: %s \n", disk)
	}

	var freeBytesAvailable uint64
	var totalNumberOfBytes uint64
	var totalNumberOfFreeBytes uint64

	err := windows.GetDiskFreeSpaceEx(windows.StringToUTF16Ptr(disk),
		&freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return 0, err
	}

	return float64(freeBytesAvailable / 1024 / 1024 / 1024), nil
}
