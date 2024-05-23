package cpu

import "testing"

/**
@file:
@author: levi.Tang
@time: 2024/8/23 11:21
@description:
**/

func sum2(s []int64) int64 {
	var total int64
	for i := 0; i < len(s); i += 2 {
		total += s[i]
	}
	return total
}

func sum8(s []int64) int64 {
	var total int64
	for i := 0; i < len(s); i += 8 {
		total += s[i]
	}
	return total
}

func BenchmarkSum2(b *testing.B) {
	sum := make([]int64, 0, 100000)
	for i := 0; i < 100000; i++ {
		sum = append(sum, int64(i))
	}

	for i := 0; i < b.N; i++ {
		sum2(sum)
	}
}

func BenchmarkSum8(b *testing.B) {
	sum := make([]int64, 0, 100000)
	for i := 0; i < 100000; i++ {
		sum = append(sum, int64(i))
	}

	for i := 0; i < b.N; i++ {
		sum8(sum)
	}
}
