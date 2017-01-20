A Crawler Package for Golang
=====
This is a Go package for crawling website resource.

Installation
-----

```
go get -u github.com/hiro511/crawler
```

Example
-----

```
package main

import (
	"fmt"
	"net/url"

	"github.com/hiro511/crawler"
)

func main() {
	url, _ := url.Parse("http://example.com")
	err := crawler.Start(url, 1, "./test/")
	if err != nil {
		fmt.Println(err)
	}
}
```
