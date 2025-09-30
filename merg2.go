package main

import "fmt"

func main() {
	count := new(int)
	_ = MergeSort([]int{2, 4, 1, 3, 5}, count)
	fmt.Println("No of inversions", *count)

}
func MergeSort(arr []int, count *int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	leftHalf := arr[:mid]
	rightHalf := arr[mid:]

	sortedLeft := MergeSort(leftHalf, count)
	sortedRight := MergeSort(rightHalf, count)

	return merge(sortedLeft, sortedRight, count)
}

func merge(left []int, right []int, count *int) []int {
	result := []int{}
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i = i + 1
		} else {
			result = append(result, right[j])
			*count += len(left) - i
			j = j + 1
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}
