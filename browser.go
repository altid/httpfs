package main

import (
	"context"
	"errors"

	"github.com/altid/libs/fs"
)

// google.ca | google.ca[1] | google.ca [2]
type browser struct {
	buffers map[string]*buffer
	cancel  context.CancelFunc
}

// TODO(halfwit): We want to handle back/forward, so we'll need a list for each browser entry
type buffer struct {
	history []*tuple
	current int
}

type tuple struct {
	uri string
	req string
}

func newBrowser(cancel context.CancelFunc) (*browser, error) {
	b := &browser{
		buffers: make(map[string]*buffer),
		cancel:  cancel,
	}

	uri, err := getURI(*homepage)
	if err != nil {
		return nil, err
	}

	b.push(uri, *homepage)
	return b, nil
}

func (b browser) Run(c *fs.Control, cmd *fs.Command) error {
	switch cmd.Name {
	case "field":
		// This is denoted by an [>Tag](hint) 
	case "open":
		req := "https://" + cmd.Args[0]

		uri, err := getURI(req)
		if err != nil {
			return err
		}

		b.push(uri, req)
		c.CreateBuffer(uri, "document")
		return fetchSite(c, uri, req)
	case "close":
		b.pop(cmd.Args[0])
		return c.DeleteBuffer(cmd.Args[0], "document")
	case "link":
		// find uri based on from
		uri, err := getURI(cmd.Args[0])
		if err != nil {
			return err
		}

		if cmd.From == "" {
			return errors.New("Unable to swap buffers")
		}

		if b.uri(cmd.From) != uri {
			b.pop(cmd.Args[0])
			c.DeleteBuffer(cmd.Args[0], "document")

			// After delete, send open cmd
			cmd.Name = "open"
			return b.Run(c, cmd)
		}

		return fetchSite(c, uri, cmd.Args[0])
	//case "reload":
	//	fetchSite(c, from, b.uri(from))
	case "back":
	case "forward":
	}

	return nil
}

func (b *browser) Quit() {
	b.cancel()
}

func (b *browser) push(uri, req string) {
	t := &tuple{
		uri: uri,
		req: req,
	}

	if _, ok := b.buffers[uri]; !ok {
		// make new and append
		b.buffers[uri] = &buffer{}
	}

	b.buffers[uri].history = append(b.buffers[uri].history, t)
	b.buffers[uri].current = len(b.buffers[uri].history)
}

func (b *browser) pop(uri string) {

}

func (b *browser) uri(req string) string {
	for _, buff := range b.buffers {
		if buff.history[buff.current].req == req {
			return buff.history[buff.current].uri
		}
	}

	return ""
}
