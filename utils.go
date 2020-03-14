package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/altid/libs/fs"
	"github.com/altid/libs/html"
	"github.com/altid/libs/markup"
)

func getUri(request string) (string, error) {
	u, err := url.Parse(request)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

type sometype struct {
	c   *fs.Control
	s   *fs.WriteCloser
	uri string
	url string
}

func (p *sometype) Img(link string) error {
	u, err := url.Parse(p.url)
	u.Path = path.Dir(u.Path) + "/" + link

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	name, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("Unable to parse URL for image %v", err)
	}

	img, err := p.c.ImageWriter(p.uri, name.Path)
	if err != nil {
		return err
	}

	defer img.Close()
	_, err = io.Copy(img, resp.Body)

	return err
}

// Called for each line in a <nav>
func (p *sometype) Nav(u *markup.Url) error {
	fmt.Fprintf(p.s, "%s\n", u)
	return nil
}

func fetchSite(c *fs.Control, uri, url string) error {
	// TODO: TLS for https, try to upgrade http as well
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	m, err := c.MainWriter(uri, "document")
	if err != nil {
		return err
	}

	s, err := c.NavWriter(uri)
	if err != nil {
		return err
	}

	p := &sometype{
		c:   c,
		uri: uri,
		url: url,
		s:   s,
	}
	defer p.s.Close()

	body, err := html.NewHTMLCleaner(m, p)
	if err != nil {
		return err
	}

	defer body.Close()

	if err := body.Parse(resp.Body); err != io.EOF {
		return err
	}

	return nil
}
