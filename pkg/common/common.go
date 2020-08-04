package common

import "os"

func GetEnvVariableString(key, defaultValue string) string {
	var res string
	if res = os.Getenv(key); res == "" {
		res = defaultValue
	}
	return res
}
