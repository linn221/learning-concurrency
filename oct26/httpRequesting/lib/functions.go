package lib

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Result struct {
	Url        string
	StatusCode int
	Body       string
	Latency    time.Duration
	Error      error
}

func (r *Result) Print() {
	if r.Error != nil {
		fmt.Printf("Url: %s\nError: %v\n\n", r.Url, r.Error)
	} else {
		fmt.Printf("Url: %s\nStatusCode: %d\nLatency: %v\n\n", r.Url, r.StatusCode, r.Latency)
	}
}

func FetchUrl(url string) Result {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return Result{
			Url:   url,
			Error: err,
		}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	latency := time.Since(start).Round(time.Millisecond)
	return Result{
		Url:        url,
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Latency:    latency,
	}
}
