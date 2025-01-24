package utils

import "log"

func LogError(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func LogInfo(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}
