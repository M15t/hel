package receiver

import (
	"fmt"
	"math"
	"strings"
)

func prettyByteSize(b int64) string {
	bf := float64(b)
	for _, unit := range []string{"", "K", "M", "G", "T", "P", "E", "Z"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", bf)
}

func prettyURL(url string) string {
	// Find the position of the question mark
	index := strings.Index(url, "?")
	if index != -1 {
		// Remove the query parameters and their values by truncating the URL
		url = url[:index]
	}

	return url
}
