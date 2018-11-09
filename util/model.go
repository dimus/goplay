package util

import (
	"runtime"
	"time"
)

type Model struct {
	Verifier
}

type Verifier struct {
	URL                string
	BatchSize          int
	Workers            int
	WaitTimeout        time.Duration
	Sources            []int
	Verify             bool
	AdvancedResolution bool
}

func NewModel() *Model {
	m := &Model{
		Verifier{
			URL:         "http://index.globalnames.org/api/graphql",
			WaitTimeout: 90 * time.Second,
			BatchSize:   500,
			Workers:     runtime.NumCPU(),
		},
	}
	return m
}
