package main

import (
	"io"
	"net/http"
	"net/url"

	"github.com/altid/cleanmark"
	fs "github.com/altid/fslib"
)

func getUri(request string) (string, error) {
	u, err := url.Parse(request)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

func fetchSite(c *fs.Control, uri, url string) error {
	// maybe try tls dialing for everything we do eventually, but for now just tcp
	// TODO: fallback to http:// if this fails as well
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	m := c.MainWriter(uri, "document")
	body := cleanmark.NewHTMLCleaner(m)
	defer body.Close()
	if err := body.Parse(resp.Body); err != io.EOF {
		return err
	}
	return nil
	// We need to eventually parse out the title and navigation elements
	// read any <navigation>, make a sidebar
}
