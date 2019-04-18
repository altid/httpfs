package main

import (
	fs "github.com/altid/fslib"
)

// google.ca | google.ca[1] | google.ca [2]
type browser map[string]string

// TODO(halfwit): We want to handle back/forward, so we'll need a list for each browser entry
type buffer struct {
	history []string
	current int
}

func (b browser) Open(c *fs.Control, request string) error {
	req := "https://" + request
	uri, err := getUri(req)
	if err != nil {
		return err
	}
	// make sure we create new one like google.ca/search[1] if a buffer already exists
	b[uri] = req
	c.CreateBuffer(uri, "document")
	return fetchSite(c, uri, req)
}

func (b browser) Close(c *fs.Control, uri string) error {
	delete(b, uri)
	return c.DeleteBuffer(uri, "document")
}

func (b browser) Link(c *fs.Control, from, request string) error {
	// find uri based on from
	uri, err := getUri(request)
	if err != nil {
		return err
	}
	if b[from] != uri {
		b.Close(c, from)
		return b.Open(c, request)
	}
	return fetchSite(c, uri, request)
}

func (b browser) Default(c *fs.Control, cmd, from, msg string) error {
	switch cmd {
	case "reload":
		fetchSite(c, from, b[from])
	}
	return nil
}

