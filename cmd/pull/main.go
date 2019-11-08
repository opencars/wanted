package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const (
	bom0 = 0xef
	bom1 = 0xbb
	bom2 = 0xbf
)

type Resp struct {
	Result Result `json:"result"`
}

type Revision struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type Result struct {
	Revisions []Revision `json:"resource_revisions"`
}

// NewReader returns an io.Reader that will skip over initial UTF-8 byte order marks.
func NewReader(r io.Reader) io.Reader {
	buf := bufio.NewReader(r)
	b, err := buf.Peek(3)
	if err != nil {
		// not enough bytes
		return buf
	}
	if b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
		buf.Discard(3)
	}
	return buf
}

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, NewReader(resp.Body))
	return err
}

func main() {
	res, err := http.Get("https://data.gov.ua/api/3/action/resource_show?id=06e65b06-3120-4713-8003-7905a83f95f5")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var resp Resp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}

	revisions := make(chan Revision)
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			for {
				rev, ok := <-revisions
				if !ok {
					wg.Done()
					return
				}
				if err := DownloadFile(rev.Name, rev.URL); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Println(rev.URL)
			}
		}()
	}

	for i, revision := range resp.Result.Revisions {
		revision.Name = fmt.Sprintf("./data/%d-%s", i+1, revision.Name)
		revisions <- revision
	}

	close(revisions)
	wg.Wait()

	fmt.Println("Amount of revisions:", len(resp.Result.Revisions))
}
