package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func downloader(id int, jobs <-chan url.URL) {
	for url := range jobs {
		saveWithURL(url, cache[url])
	}
}

func download(url *url.URL) (*http.Response, error) {
	fmt.Println("downloading: " + url.String())
	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func save(content []byte, toFile string) error {
	file, err := os.OpenFile(basePath+toFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(content)
	return nil
}

func saveStr(str string, toFile string) error {
	return save([]byte(str), toFile)
}

func saveWithURL(url url.URL, toFile string) {
	defer wg.Done()
	res, err := download(&url)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = save(body, toFile); err != nil {
		fmt.Println(err)
		return
	}
}
