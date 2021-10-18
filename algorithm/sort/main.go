package main

import (
	"math/rand"
	"time"
)

const (
	num      = 10000  // 测试数组的长度
	rangeNum = 100000 // 数组元素大小范围
)

func main() {
	//arr := GenerateRand() //生成随机数组
	//排序前 复制原数组
	//org_arr := make([]int, num)
	//copy(org_arr, arr)
	//冒泡排序
	//Bubble(arr)
	// 选择排序
	//SelectSort(arr)
	// 插入排序
	//InsertSort(arr)
	//快速排序
	//QuickSort(arr, 0, len(arr)-1)
	// 归并排序
	Guibing()
	//MergeSort(arr, 0, len(arr)-1)
	// 堆排序
	//HeapSort(arr)
	//sort.Ints(org_arr) //使sort模块对原数组排序
	//sort.Sort(sort.IntSlice(org_arr))
	//fmt.Println(arr, org_arr, IsSame(arr, org_arr))
	//打印前15个数,并对比排序是否正确
	//fmt.Println(arr[:15], org_arr[:15], IsSame(arr, org_arr))
}

//生成随机数组
func GenerateRand() []int {
	randSeed := rand.New(rand.NewSource(time.Now().Unix() + time.Now().UnixNano()))
	arr := make([]int, num)
	for i := 0; i < num; i++ {
		arr[i] = randSeed.Intn(rangeNum)
	}
	return arr
}

//比较两个切片
func IsSame(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	if (slice1 == nil) != (slice2 == nil) {
		return false
	}

	for i, num := range slice1 {
		if num != slice2[i] {
			return false
		}
	}
	return true
}
