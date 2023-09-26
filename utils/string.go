package utils

import (
	"fmt"
	"strings"
)

func GetTransactionDate(originalDate string) string {
	parts := strings.Split(originalDate, "-")
	newstr := fmt.Sprintf("%s-%s", parts[0], parts[1])
	return newstr
}