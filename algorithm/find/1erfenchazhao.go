package main

import (
	"fmt"
	"time"
)

func binary_search_only(I []int, n int) int {
	i, j := 0, len(I)
	for i < j {
		if I[(i+j)/2] > n {
			j = (i + j) / 2
		} else if I[(i+j)/2] < n {
			i = (i + j) / 2
		} else if I[(i+j)/2] == n {
			return (i + j) / 2
		}
		fmt.Println("i:", i, " j:", j, " n:", n)
		time.Sleep(time.Second * 1)
	}
	return -1
}
