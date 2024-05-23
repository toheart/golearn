package tnet

import (
	"fmt"
	"net/url"
	"testing"
)

/**
@file:
@author: levi.Tang
@time: 2024/4/16 10:49
@description:
**/

func TestUrl(t *testing.T) {
	// 解析url地址
	u, err := url.Parse("http://1.2.3.44:999/search?q=dotnet")
	if err != nil {
		panic(err)
	}

	// 打印格式化的地址信息
	fmt.Println(u.Scheme)     // 返回协议
	fmt.Println(u.Hostname()) // 返回域名
	fmt.Println(u.Path)       // 返回路径部分
	fmt.Println(u.Port())
	fmt.Println(u.RawQuery) // 返回url的参数部分

	params := u.Query() // 以url.Values数据类型的形式返回url参数部分,可以根据参数名读写参数

	fmt.Println(params.Get("q")) // 读取参数q的值
}
