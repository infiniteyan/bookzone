package util

import "fmt"

func FormatBytes(size int64) string {
	units := []string{" B", " KB", " MB", " GB", " TB"}

	s := float64(size)

	i := 0

	for ; s >= 1024 && i < 4; i++ {
		s /= 1024
	}

	return fmt.Sprintf("%.2f%s", s, units[i])
}
