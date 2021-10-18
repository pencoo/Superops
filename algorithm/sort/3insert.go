package main

//思路，从第一个元素开始，与上一个元素对比，如果小于前一个元素，则将此元素迁移，并在此将此元素与前面元素对比。
//直到元素插入的位置使前面已排序部分为正确序列。依次直到列表排序结束

func InsertSort(arr []int) {
	for i := 1; i <= len(arr)-1; i++ {
		for j := i; j > 0; j-- {
			if arr[j-1] > arr[j] {
				arr[j-1], arr[j] = arr[j], arr[j-1]
			}
		}
	}
}
