package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	downloadDir = "/opt/www/wallpaper.theyan.gs/wallpapers/NASA/"
	nasaFeedURL = "https://www.nasa.gov/feeds/iotd-feed"
)

func downloadPic(url string, dir string) error {
	fileName := filepath.Base(url)
	fullName := filepath.Join(dir, fileName)
	if _, err := os.Open(fullName); os.IsNotExist(err) {
		client := http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("Failed to download file: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Download failed: %s", resp.Status)
		}
		f, err := os.OpenFile(fullName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return fmt.Errorf("Failed to create file: %w", err)
		}
		defer f.Close()
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return fmt.Errorf("Failed to copy file: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	} else {
		log.Printf("File %s already exists", fileName)
		return nil
	}
	return nil
}

func parseNASAFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse NASA feed: %w", err)
	}
	return feed, nil
}

func main() {
	feed, err := parseNASAFeed(nasaFeedURL)
	if err != nil {
		log.Printf("Failed to parse NASA feed: %v", err)
		return
	}
	log.Println(feed.Title)
	for _, item := range feed.Items {
		if len(item.Enclosures) > 0 {
			log.Println(item.Enclosures[0].URL)
			url := item.Enclosures[0].URL
			err := downloadPic(url, downloadDir)
			if err != nil {
				log.Printf("Failed to download picture: %v", err)
			}
		}
	}
}
