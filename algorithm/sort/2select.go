package main

//思路，先找最小的放在第一位，在找到第二小的放在第二位，即从小到大排列，直到数据完成排序

func SelectSort(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j <= len(arr)-1; j++ {
			if arr[j] < arr[i] {
				arr[j], arr[i] = arr[i], arr[j]
			}
		}
	}
}
