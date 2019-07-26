package helpers

import (
	"fmt"
)

// PanicError - progress over perfection, but sometimes you just need to panic!
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

// ClearStdoutLine clear previous stdout line
func ClearStdoutLine() {
	fmt.Printf("\033[F") // back to previous line
	fmt.Printf("\033[K") // clear line
}

// ByteCountDecimal provides a human string for storage based on 1000
func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// ByteCountBinary provides a human string for storage based on 1024
func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
