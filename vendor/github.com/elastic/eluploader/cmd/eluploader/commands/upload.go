package commands

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"gopkg.in/cheggaaa/pb.v1"
)

const chunkSize = int64(50 * 1048576)

var (
	semaphore  chan struct{}
	pbpool     *pb.Pool
	fileDigest string

	defaultNumWorkers = runtime.NumCPU()
	cancel            = make(chan os.Signal)
)

type UploadCmd struct {
	UploadID      string
	Filepath      string
	ApiURL        string
	NumWorkers    int
	UploadedParts map[string]Part
	mu            sync.Mutex
}

type Part struct {
	UploadID   string
	FileDigest string
	PartNr     int64
	Digest     string
}

func (cmd *UploadCmd) uploadExists() (bool, error) {
	apiURL := cmd.ApiURL + "/api/uploads/" + cmd.UploadID

	res, err := http.Head(apiURL)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, nil
	}
	return true, nil
}

func (cmd *UploadCmd) computeFileDigest(f *os.File, pbar *pb.ProgressBar) (string, error) {
	f.Seek(0, 0)
	return cmd.computeDigest(f, pbar)
}

func (cmd *UploadCmd) computeDigest(r io.Reader, pbar *pb.ProgressBar) (string, error) {
	var dw = sha256.New()
	rr := pbar.NewProxyReader(r)
	_, err := io.Copy(dw, rr)
	if err != nil {
		return "", fmt.Errorf("Failed to compute digest: %s", err)
	}
	return hex.EncodeToString(dw.Sum(nil)), nil
}

func (cmd *UploadCmd) readPart(f *os.File, n int64) (*bytes.Buffer, error) {
	buf := make([]byte, n)
	_, err := f.Read(buf)
	if err != nil {
		return bytes.NewBuffer(buf), fmt.Errorf("Error reading chunk: %s", err)
	}
	return bytes.NewBuffer(buf), nil
}

func (cmd *UploadCmd) uploadPart(partNr int64, bytesCount int64, buf *bytes.Buffer, pbar *pb.ProgressBar) (Part, error) {
	semaphore <- struct{}{}        // acquire a token
	defer func() { <-semaphore }() // release the token

	part := Part{PartNr: partNr, UploadID: cmd.UploadID, FileDigest: fileDigest}

	pbpool.Add(pbar)

	// TEMP: Fake errors
	// if partNr == 3 {
	// 	pbar.Prefix(fmt.Sprintf("Part %-4d[ERROR] ", partNr))
	// 	return part, fmt.Errorf("Fake error")
	// }

	bw := bytes.NewBuffer([]byte{})
	tr := io.TeeReader(buf, bw)

	partDigest, err := cmd.computeDigest(tr, pb.New64(bytesCount))
	if err != nil {
		pbar.Finish()
		return part, fmt.Errorf("Failed to compute part digest: %s", err)
	}
	part.Digest = partDigest

	key := strings.Join([]string{part.UploadID, part.FileDigest, part.Digest}, "-")

	cmd.mu.Lock()
	part, ok := cmd.UploadedParts[key]
	cmd.mu.Unlock()
	if ok {
		pbar.Prefix(fmt.Sprintf("Part %-4d[SKIPPED] ", partNr))
		pbar.Format("[x>.]")
		pbar.Finish()
		return part, nil
	}

	r := pbar.NewProxyReader(bw)

	client := &http.Client{}

	params := url.Values{}
	params.Set("filename", filepath.Base(cmd.Filepath))
	params.Set("file_digest", fileDigest)
	params.Set("part_number", fmt.Sprintf("%d", partNr))
	params.Set("part_digest", partDigest)

	apiURL := cmd.ApiURL + "/api/uploads/" + cmd.UploadID + "?" + params.Encode()
	req, err := http.NewRequest("PUT", apiURL, r)
	if err != nil {
		pbar.Prefix(fmt.Sprintf("Part %-4d[ERROR] ", partNr))
		return part, fmt.Errorf("Failed to create HTTP request: %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		pbar.Prefix(fmt.Sprintf("Part %-4d[ERROR] ", partNr))
		return part, fmt.Errorf("Failed to get HTTP response: %s", err)
	}
	if resp.StatusCode != http.StatusCreated {
		pbar.Prefix(fmt.Sprintf("Part %-4d[ERROR] ", partNr))
		return part, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return part, nil
}

func (cmd *UploadCmd) storeUploadedParts() (string, error) {
	filename := fmt.Sprintf("%s.json", cmd.UploadID)
	b, err := json.Marshal(cmd.UploadedParts)
	if err != nil {
		return filename, fmt.Errorf("failed to encode metadata: %s", err)
	}
	if err = ioutil.WriteFile(filename, b, 0644); err != nil {
		return filename, fmt.Errorf("failed to store metadata: %s", err)
	}
	return filename, nil
}

func (cmd *UploadCmd) loadUploadedParts() error {
	filename := fmt.Sprintf("%s.json", cmd.UploadID)

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			cmd.UploadedParts = make(map[string]Part)
			return nil
		}
		return fmt.Errorf("failed to check metadata file: %s", err)
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to load metadata file: %s", err)
	}

	if err = json.Unmarshal(b, &cmd.UploadedParts); err != nil {
		return fmt.Errorf("failed to parse metadata: %s", err)
	}

	return nil
}

func newProgressBar(prefix string, total int64) *pb.ProgressBar {
	bar := pb.New64(total).Prefix(prefix)
	bar.ShowCounters = true
	bar.ShowTimeLeft = false
	bar.ShowFinalTime = false
	bar.ShowSpeed = false
	bar.SetUnits(pb.U_BYTES)
	bar.SetMaxWidth(80)
	bar.Format("[~>-]")
	return bar
}

func (cmd *UploadCmd) Execute() {
	if ok, err := cmd.uploadExists(); !ok {
		fmt.Printf("[!] The upload with ID %q does not exist\n", cmd.UploadID)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		os.Exit(1)
	}

	// Semaphore; https://github.com/adonovan/gopl.io/blob/master/ch8/crawl2/findlinks.go
	semaphore = make(chan struct{}, cmd.NumWorkers)
	defer close(semaphore)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Upload ID: %s\n", cmd.UploadID)
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("File:      %s\n", cmd.Filepath)
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\nComputing SHA256 digest")
	fmt.Println(strings.Repeat("-", 80))

	f, err := os.Open(cmd.Filepath)
	if err != nil {
		fmt.Printf("[!] Error opening file: %s", err)
		os.Exit(1)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("[!] Error getting file information: %s\n", err)
		os.Exit(1)
	}

	fileSize := fi.Size()

	// ---- SHA256 --------------------------------------------------------------
	//
	pbar := newProgressBar("", fileSize)
	pbar.Start()

	fileDigest, err = cmd.computeFileDigest(f, pbar)
	if err != nil {
		fmt.Printf("[!] Error computing digest: %s\n", err)
		os.Exit(1)
	}
	pbar.Finish()
	fmt.Println("------", fileDigest, "--------")

	// ---- Upload chunks -------------------------------------------------------
	//
	var n int64
	var partNr int64
	var wg sync.WaitGroup
	var numParts = int(math.Ceil(float64(fileSize) / float64(chunkSize)))
	var errors = make([]error, 0)

	if err := cmd.loadUploadedParts(); err != nil {
		fmt.Printf("[!] Failed to load previously uploaded parts: %s\n", err)
	}

	fmt.Printf("\nUploading the file in %d parts with %d workers\n", numParts, cmd.NumWorkers)
	fmt.Println(strings.Repeat("-", 80))

	f.Seek(0, 0)
	pbpool = pb.NewPool()
	pbpool.Start()
	signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cancel
		pbpool.Stop()
		fmt.Println("\nExiting...")
		os.Exit(1)
	}()

	for n = 0; n < fileSize; n = n + chunkSize {
		partNr++

		nb := chunkSize
		if nb > fileSize {
			nb = fileSize
		}
		if n+chunkSize > fileSize {
			nb = fileSize - n
		}

		buf, err := cmd.readPart(f, nb)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		wg.Add(1)
		go func(partNr int64, buf *bytes.Buffer, nb int64) {
			defer wg.Done()

			bar := newProgressBar(fmt.Sprintf("Part %-4d", partNr), nb)
			bar.ShowCounters = false
			defer bar.Finish()

			part, err := cmd.uploadPart(partNr, nb, buf, bar)
			if err != nil {
				errors = append(errors, fmt.Errorf("Uploading part %d failed: %s\n", partNr, err))
			} else {
				key := strings.Join([]string{part.UploadID, part.FileDigest, part.Digest}, "-")
				cmd.mu.Lock()
				cmd.UploadedParts[key] = part
				cmd.mu.Unlock()
			}
		}(partNr, buf, nb)
	}

	wg.Wait()
	pbpool.Stop()

	if len(errors) > 0 {
		fmt.Printf("\n[ERROR] Failed to upload %d parts\n", len(errors))
		fmt.Println(strings.Repeat("-", 80))
		for _, err := range errors {
			fmt.Printf("* %s", err)
		}
		filename, err := cmd.storeUploadedParts()
		if err != nil {
			fmt.Println("[!] Error storing metadata about succesfully uploaded parts")
		} else {
			fmt.Printf("------ Upload metadata stored to %s -----\n\n", filename)
			fmt.Println(strings.Repeat("=", 80))
			fmt.Println("\n             Run the command again to retry uploading failed parts")
			fmt.Println(strings.Repeat("=", 80))
		}
		os.Exit(1)
	}

	fmt.Printf("\nFinalizing the file\n")
	fmt.Println(strings.Repeat("-", 80))

	apiURL := cmd.ApiURL + "/api/uploads/" + cmd.UploadID + "/" + fileDigest + "/_finalize"
	client := &http.Client{}
	resp, err := client.Post(apiURL, "application/json", nil)
	if err != nil {
		fmt.Printf("Failed to get HTTP response: %s\n", err)
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error when finalizing the file: %s\n", resp.Status)
		os.Exit(1)
	}

	fmt.Print("File sucessfully uploaded and finalized\n")

	j := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		fmt.Printf("Error loading service JSON response: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Printf("--> %s/d/%s\n", cmd.ApiURL, j["slug"])
	os.Exit(0)
}
