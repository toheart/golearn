package utils

import "os"

/**
@file:
@author: levi.Tang
@time: 2024/5/22 9:52
@description:
**/

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
