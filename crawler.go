package crawler

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const maxGoRoutine = 150

var basePath string
var tags = map[string]string{
	"link":   "href",
	"script": "src",
	"input":  "src",
	"img":    "src",
	"area":   "href",
}
var cache map[url.URL]string
var localFile map[string]bool
var isCrawled map[url.URL]bool
var isDL map[url.URL]bool
var jobs chan url.URL
var wg *sync.WaitGroup

// Website struct
type Website struct {
	URL             *url.URL
	Links, Elements map[url.URL]string
	HTML            string
}

// Start to crawle with URL
func Start(url *url.URL, maxDepth int, dir string) error {
	return start(url, maxDepth, dir)
}

func start(url1 *url.URL, maxDepth int, dir string) error {
	initilize(dir)
	website := &Website{Links: make(map[url.URL]string)}
	website.Links[*url1] = ""
	nextWebsites := []*Website{website}
	isLast := false
	for i := 0; i < maxDepth; i++ {
		if i == maxDepth-1 {
			isLast = true
		}
		websites := make([]*Website, len(nextWebsites))
		copy(websites, nextWebsites)
		nextWebsites = []*Website{}
		for _, website := range websites {
			for link := range website.Links {
				if _, ok := isCrawled[link]; ok {
					continue
				}
				next, err := crawl(&link, isLast)
				if err != nil {
					fmt.Println(err)
				}
				nextWebsites = append(nextWebsites, next)
				isCrawled[link] = true
			}
		}
	}
	close(jobs)
	wg.Wait()
	return nil
}

func initilize(dir string) {
	basePath = dir
	cache = make(map[url.URL]string)
	localFile = make(map[string]bool)
	isCrawled = make(map[url.URL]bool)
	isDL = make(map[url.URL]bool)
	jobs = make(chan url.URL, 1000)
	wg = new(sync.WaitGroup)
	removeBaseDir()
	makeBaseDir()
	for i := 1; i <= maxGoRoutine; i++ {
		go downloader(i, jobs)
	}
}

func removeBaseDir() {
	if _, err := os.Stat(basePath); !os.IsNotExist(err) {
		os.RemoveAll(basePath)
	}
}
func makeBaseDir() {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.Mkdir(basePath, 0755)
	}
}

func crawl(url *url.URL, isLast bool) (*Website, error) {
	res, err := download(url)
	if err != nil {
		return nil, err
	}
	website, err := parse(res, isLast)
	if err != nil {
		return nil, err
	}

	for url := range website.Elements {
		if _, ok := isDL[url]; ok {
			continue
		}
		wg.Add(1)
		if len(jobs) == cap(jobs) {
			log.Fatal("jobs channel is full")
		}
		jobs <- url
		isDL[url] = true
	}

	local, ok := cache[*url]
	if !ok {
		local = generateName(url, ".html")
		cache[*url] = local
	}
	saveStr(website.HTML, local)
	return website, nil
}

func generateName(url *url.URL, suffix string) string {
	name := path.Base(url.Path)
	if name == "." || name == "/" {
		name = strings.Split(url.Host, ".")[0]
		if name == "www" {
			name = strings.Split(url.Host, ".")[1]
		}
	}
	ext := path.Ext(name)
	if ext == "" {
		ext = ".html"
		name += ext
	}
	_, err := os.Stat(name)
	if _, ok := localFile[name]; !os.IsNotExist(err) || ok {
		time := time.Now()
		const layout = "_20060102_150405"
		name = strings.Replace(name, ext, time.Format(layout)+ext, 1)
	}
	localFile[name] = false
	return name
}
