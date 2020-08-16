package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/opencars/govdata"

	"github.com/opencars/wanted/pkg/bom"
	"github.com/opencars/wanted/pkg/logger"
)

// DownloadFile downloads file from the url to the specified filepath.
func DownloadFile(filepath string, url string) error {
	// Get the data.
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file.
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file.
	utf, err := bom.NewReader(resp.Body)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, utf)

	return err
}

type downloader struct {
	mp        sync.Mutex
	progress  int
	revisions chan govdata.Revision

	wg sync.WaitGroup
}

// Download all revisions from the gov.data.ua resource.
func (d *downloader) Download(resource *govdata.Resource) {
	for i := 0; i < 10; i++ {
		d.wg.Add(1)

		go func() {
			for {
				rev, ok := <-d.revisions
				if !ok {
					d.wg.Done()
					return
				}

				if err := DownloadFile(rev.Name, rev.URL); err != nil {
					logger.Fatal(err)
				}

				d.mp.Lock()
				d.progress++
				d.mp.Unlock()
			}
		}()
	}

	for i := len(resource.Revisions); i > 0; i-- {
		parts := strings.Split(resource.Revisions[i-1].URL, "/")
		resource.Revisions[i-1].Name = fmt.Sprintf("./data/%s.json", parts[len(parts)-1])
		d.revisions <- resource.Revisions[i-1]
	}

	close(d.revisions)
	d.wg.Wait()
}

// ProgressBar shows information about the installation progress.
func (d *downloader) ProgressBar(max int) {
	for {
		<-time.After(1 * time.Second)
		fmt.Printf("\rDownloading: %d/%d", d.progress, max)
	}
}

func main() {
	gov := govdata.NewClient()
	resourceID := "06e65b06-3120-4713-8003-7905a83f95f5"
	resource, err := gov.ResourceShow(context.Background(), resourceID)
	if err != nil {
		logger.Fatal(err)
	}

	downloader := downloader{
		revisions: make(chan govdata.Revision),
	}

	go downloader.ProgressBar(len(resource.Revisions))
	downloader.Download(resource)
	logger.Info("Amount of revisions: %d", len(resource.Revisions))
}
