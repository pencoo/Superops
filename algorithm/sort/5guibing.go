package main

import "fmt"

//归并排序：
//

func Guibing() {
	arr1 := []int{1, 3, 5, 7, 9, 11, 13,17,29,38,41}
	arr2 := []int{2, 4, 6, 8, 10, 12,14,15,16,19,31}
	i, j := 0, 0
	arr := make([]int, 0)
	for i < len(arr1) {
		for j < len(arr2) {
			if arr1[i] < arr2[j] {
				arr = append(arr, arr1[i])
				i++
				break
			} else if arr1[i] > arr2[j] {
				arr = append(arr, arr2[j])
				j++
			} else {
				arr = append(arr, arr1[i])
				arr = append(arr, arr2[j])
				i++
				j++
				break
			}
		}
		if j == len(arr2) {
			break
		}
	}
	if i <= len(arr1)-1 {
		for i <= len(arr1)-1 {
			arr = append(arr, arr1[i])
			i++
		}
	}
	if j <= len(arr2)-1 {
		for i <= len(arr2)-1 {
			arr = append(arr, arr2[j])
			j++
		}
	}
	fmt.Println(arr)
}
