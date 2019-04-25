package main

import (
	"fmt"
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

type sometype struct{
	c *fs.Control
	uri string
}

// These are included out of band from the parser 
// Optional handlers that are called from tokenizer
func (p *sometype) Img(link string) error {
	resp, err := http.Get("https://" + p.uri + link)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	name, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("Unable to parse URL for image %v", err)
	}
	img := p.c.ImageWriter(p.uri, name.Path)
	_, err = io.Copy(img, resp.Body)
	return err
}


// Called for each line in a <nav>
func (p *sometype) Nav(url string) error {
	fmt.Println(url)
	return nil
}

func fetchSite(c *fs.Control, uri, url string) error {
	// TODO: TLS for https, try to upgrade http as well
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	m := c.MainWriter(uri, "document")
	p := &sometype{c, uri}
	// p can be nil
	body := cleanmark.NewHTMLCleaner(m, p)
	defer body.Close()
	if err := body.Parse(resp.Body); err != io.EOF {
		return err
	}
	return nil
}
