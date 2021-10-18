package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, v := range arr {
		m := binary_search_only(arr, v)
		fmt.Println(m)
	}
}
