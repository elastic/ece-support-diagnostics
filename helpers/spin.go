package helpers

// based on https://github.com/caarlos0/spin

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	// ClearLine go to the beginning of the line and clear it
	ClearLine = "\r\033[K"
	// PreviousLine moves back to the prior line
	PreviousLine = "\033[F"
)

// Spinner types.
var (
	Box1    = `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
	Default = Box1
)

// Spinner main type
type Spinner struct {
	frames []rune
	pos    int
	active uint64
	text   string
}

// NewSpinner Spinner with args
func NewSpinner(text string) *Spinner {
	s := &Spinner{
		text: ClearLine + text,
	}
	s.Set(Default)
	return s
}

// Set frames to the given string which must not use spaces.
func (s *Spinner) Set(frames string) {
	s.frames = []rune(frames)
}

// Start shows the spinner.
func (s *Spinner) Start() *Spinner {
	if atomic.LoadUint64(&s.active) > 0 {
		return s
	}
	atomic.StoreUint64(&s.active, 1)
	go func() {
		for atomic.LoadUint64(&s.active) > 0 {
			fmt.Printf(s.text, s.next())
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return s
}

// Stop hides the spinner.
func (s *Spinner) Stop() bool {
	if x := atomic.SwapUint64(&s.active, 0); x > 0 {
		fmt.Printf(ClearLine)
		return true
	}
	return false
}

func (s *Spinner) next() string {
	r := s.frames[s.pos%len(s.frames)]
	s.pos++
	return string(r)
}
