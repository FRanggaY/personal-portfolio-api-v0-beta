package helper

import "time"

func GetStringTimeNow() string {
	// Get the current time
	currentTime := time.Now()

	// Format the current time into a string (e.g., "2006-01-02_15-04-05")
	timeString := currentTime.Format("2006-01-02_15-04-05")

	return timeString
}
