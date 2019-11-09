package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/opencars/wanted/pkg/bom"
	"github.com/opencars/wanted/pkg/govdata"
)

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
	_, err = io.Copy(out, bom.NewReader(resp.Body))
	return err
}

type downloader struct {
	mp        sync.Mutex
	progress  int
	revisions chan govdata.Revision

	wg sync.WaitGroup
}

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
					log.Fatal(err)
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

func (d *downloader) ProgressBar(max int) {
	for {
		<-time.After(1 * time.Second)
		fmt.Printf("\r%d/%d", d.progress, max)
	}
}

func main() {
	gov := govdata.NewClient()
	resourceID := "06e65b06-3120-4713-8003-7905a83f95f5"
	resource, err := gov.ResourceShow(context.Background(), resourceID)
	if err != nil {
		log.Fatal(err)
	}

	downloader := downloader{
		revisions: make(chan govdata.Revision),
	}

	go downloader.ProgressBar(len(resource.Revisions))

	downloader.Download(resource)
	log.Println("Amount of revisions:", len(resource.Revisions))
}
