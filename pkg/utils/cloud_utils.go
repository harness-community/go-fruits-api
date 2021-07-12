package utils

import (
	"math/rand"
	"os"
)

var (
	clouds = [...]string{"gce", "aws", "azure"}
)

func WhichCloud() string {
	if cloud := os.Getenv("CLOUD_PROVIDER"); cloud == "" {
		r := rand.Intn(3)
		return clouds[r]
	} else {
		return cloud
	}
}
