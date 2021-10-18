package main

//思路，在切片中选出一个基数，从数组中查找将小于基数的数放左边，大于基数的数放右边
//递归上面操作从而使切片有序
func QuickSort(arr []int, l, r int) {
	if l < r {
		pivot := arr[r]          //基数
		i := l - 1               //小于基数的元素切片指针位置
		for j := l; j < r; j++ { //循环切片内所有元素
			if arr[j] <= pivot { //找出当前切片中小于基数的数
				i++                             //找到小于基数的数后，此指针右移一位
				arr[j], arr[i] = arr[i], arr[j] //将小于基数的数移到较小一边
			}
		}
		i++                             //当前i和i左边的值都是小于基数的数，将指针加一定位到下一个元素后与基数对调，这样基数左边都小于基数，基数右边都大于基数
		arr[r], arr[i] = arr[i], arr[r] //对调
		QuickSort(arr, l, i-1)          //递归小于基数的切片
		QuickSort(arr, i+1, r)          //递归大于基数的切片
	}
}

//基础版，上面是改进版。既节省内存有能快速完成排序
func quickSort(arr []int) []int {
	n := len(arr)
	//如果n为1，即数组只要一个元素，不需要排序，返回即可
	if n < 2 {
		return arr
	} else {
		middle := arr[0] //获取比较的参考值（基准值）
		var low []int    //小于基准值的数据组成的切片
		var high []int   //大于基准值的数据组成的切片
		for i := 1; i < n; i++ {
			if arr[i] < middle {
				low = append(low, arr[i]) //小于基准值的数据归类
			} else {
				high = append(high, arr[i]) //大于基准值的数据归类
			}
		}
		/*
			以下为组合成的切片：quickSort（low）+middle+quickSort(high)
			返回值为：   递归小于部分 + 基准值  +  递归大于部分
		*/
		lowSlice := quickSort(low)
		highSlice := quickSort(high)
		res := append(lowSlice, middle)
		for _, data := range highSlice {
			res = append(res, data)
		}
		return res
	}
}
