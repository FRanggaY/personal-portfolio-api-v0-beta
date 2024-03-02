package helper

import (
	"strconv"
)

func ParsePageSize(pageSizeStr string) int {
	// Parse page size from string, or use default value
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		return 5 // Default page size
	}
	return pageSize
}

func ParsePageNumber(pageNumberStr string) int {
	// Parse page number from string, or use default value
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNumber <= 0 {
		return 1 // Default page number
	}
	return pageNumber
}

func ParseUserID(userIdString string) int64 {
	userID, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return 0
	}
	return userID
}
