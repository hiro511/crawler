package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"
)

func parse(res *http.Response, isLast bool) (*Website, error) {
	website := &Website{
		URL:      res.Request.URL,
		Links:    make(map[url.URL]string),
		Elements: make(map[url.URL]string),
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return website, err
	}

	if !isLast {
		doc.Find("a").Each(func(_ int, s *goquery.Selection) {
			rawurl, ok := s.Attr("href")
			if !ok {
				return
			}
			url, err := url.Parse(rawurl)
			if err != nil || url.Scheme == "android-app" ||
				url.Scheme == "javascript" || url.Scheme == "mailto" {
				fmt.Println("invalid url: " + rawurl)
				return
			}
			absURL := website.URL.ResolveReference(url)
			local, ok := cache[*absURL]
			if !ok {
				local = generateName(absURL, "")
				cache[*absURL] = local
			}
			s.SetAttr("href", local)
			website.Links[*absURL] = local
		})
	}

	for tag := range tags {
		doc.Find(tag).Each(func(_ int, s *goquery.Selection) {
			attr := tags[tag]
			rawurl, ok := s.Attr(attr)
			if !ok {
				return
			}
			url, err := url.Parse(rawurl)
			if err != nil || url.Scheme == "android-app" ||
				url.Scheme == "javascript" || url.Scheme == "mailto" {
				fmt.Println("invalid url: " + rawurl)
				return
			}
			absURL := website.URL.ResolveReference(url)
			local, ok := cache[*absURL]
			if !ok {
				local = generateName(absURL, "")
				cache[*absURL] = local
			}
			s.SetAttr(attr, local)
			if path.Ext(local) == ".html" && tag != "script" {
				website.Links[*absURL] = local
			} else {
				website.Elements[*absURL] = local
			}
		})
	}

	if website.HTML, err = doc.Html(); err != nil {
		return nil, err
	}
	return website, nil
}
