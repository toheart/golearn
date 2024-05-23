package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"regexp"
)

/**
@file:
@author: levi.Tang
@time: 2024/8/7 16:28
@description:
**/

func main() {
	free, err := GetDiskFreeSpace("K:\\")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("GetfreeDisk: %.2f \n ", free)
}

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
