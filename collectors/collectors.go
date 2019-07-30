package collector

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/elastic/ece-support-diagnostics/config"
)

// StartCollector StartCollector
func StartCollector(fn interface{}, messageCh chan<- string, cfg *config.Config, wg *sync.WaitGroup) {
	start := time.Now()
	// tR := rand.Intn(10)
	// timeout := time.Duration(tR) * time.Second
	timeout := 3 * time.Minute

	// set ctx timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// this channel is used by a function that can be canceled
	//  it avoids having the canceled function try to write to the
	//  main channel that might already be closed (and cause a panic)
	ch := make(chan string, 1)

	// provide function signature and pass arguments
	go fn.(func(chan<- string, *config.Config))(ch, cfg)

	// Timeout or Task is Done
	select {
	case <-ctx.Done():
		funcName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		messageCh <- fmt.Sprintf("\u2715 %s, canceled by %s timeout (took: %s)", funcName, timeout, time.Since(start))
	case res := <-ch:
		messageCh <- fmt.Sprintf("%s (took: %s)", res, time.Since(start))
	}

	wg.Done()
}
