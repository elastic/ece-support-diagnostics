package helpers

import (
	"fmt"
)

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func ClearStdoutLine() {
	fmt.Printf("\033[F") // back to previous line
	fmt.Printf("\033[K") // clear line
}

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

// func setupFolders(folders []string) {
// 	for _, folder := range folders {
// 		f, _ := filepath.Abs(folder)
// 		flog := logp.NewLogger("folders")
// 		flog.Info("Temp folder setup: ", f)
// 		os.MkdirAll(f, os.ModePerm)
// 	}
// }

// func stringInSlice(a string, list []string) bool {
// 	for _, b := range list {
// 		if b == a {
// 			return true
// 		}
// 	}
// 	return false
// }

// func writeFile(filepath string, data []byte) error {
// 	return ioutil.WriteFile(filepath, data, 0644)
// }
