package utils

import "github.com/lib/pq"

func Contain(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func Index(arr []string, val string) int {
	for index, v := range arr {
		if v == val {
			return index
		}
	}
	return -1
}

func Remove(arr pq.StringArray, index int) pq.StringArray {
	return append(arr[:index], arr[index+1:]...)
}
