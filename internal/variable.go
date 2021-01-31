package internal

import "os"

func GetApmName() string {
	return os.Getenv("APM_NAME")
}
