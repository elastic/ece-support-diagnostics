package ecediag

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"testing"
)

func TestFile_zeroByteCheck(t *testing.T) {

	content := []byte("")
	emptyTmpFile, err := ioutil.TempFile("", "testemptyfile")
	if err != nil {
		log.Fatal(err)
	}

	fullTmpFile, err := ioutil.TempFile("", "testfullfile")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(emptyTmpFile.Name()) // clean up
	defer os.Remove(fullTmpFile.Name())  // clean up

	if _, err := emptyTmpFile.Write(content); err != nil {
		log.Fatal(err)
	}
	content = []byte("hello world")
	if _, err := fullTmpFile.Write(content); err != nil {
		log.Fatal(err)
	}

	emptyStat, _ := emptyTmpFile.Stat()
	fullStat, _ := fullTmpFile.Stat()

	empty := File{
		info:     emptyStat,
		filepath: emptyTmpFile.Name(),
	}
	file := File{
		info:     fullStat,
		filepath: fullTmpFile.Name(),
	}

	tests := []struct {
		name string
		file File
		want bool
	}{
		{
			name: "0 byte file",
			file: empty,
			want: true,
		},
		{
			name: "file with bytes",
			file: file,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.file.zeroByteCheck(tt.name); got != tt.want {
				t.Errorf("File.zeroByteCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_findPattern(t *testing.T) {

	files := Files{}

	tests := []struct {
		name     string
		files    *Files
		path     string
		re       *regexp.Regexp
		expected []string
	}{
		{
			name:  "test files",
			files: &files,
			path:  "test_data/ece_api",
			re:    regexp.MustCompile(`.*`),
			expected: []string{
				"test_data/ece_api/allocators",
				"test_data/ece_api/cluster_plan_activity",
				"test_data/ece_api/clusters_elasticsearch",
				"test_data/ece_api/clusters_elasticsearch_clean",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.files.findPattern(tt.path, tt.re)
			for i, v := range files {
				if v.filepath != tt.expected[i] {
					t.Errorf("Files.findPattern() = %v, want %v", v.filepath, tt.expected[i])
				}
			}
			// fmt.Printf("%+v\n", files)
		})
	}
}
