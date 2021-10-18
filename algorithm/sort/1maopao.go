package main

//思路，将数据从后向前排，先找到最大的放最后一位，在找到第二大的放最后第二位...直到数据完成排序
//优化思路，使用一个交换标识，如果轮转一圈无顺序变更表示已经是顺序排列

//普通逻辑
//func Bubble(arr []int) {
//	for i := len(arr) - 1; i > 0; i-- {
//		for j := 0; j < i; j++ {
//			if arr[j] > arr[j+1] {
//				arr[j], arr[j+1] = arr[j+1], arr[j]
//			}
//		}
//	}
//}

//优化后逻辑
func Bubble(arr []int) {
	size := len(arr)
	var swapped bool
	for i := size - 1; i > 0; i-- {
		swapped = false //顺序变更标签
		for j := 0; j < i; j++ {
			if arr[j+1] < arr[j] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true //判断是否有顺序变更
			}
		}
		if swapped != true { //循环以便如果已没有顺序变更表示已是标准顺序
			break
		}
	}

}
