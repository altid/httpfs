package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/altid/libs/config"
	"github.com/altid/libs/config/types"
	"github.com/altid/libs/fs"
)

var (
	mtpt     = flag.String("p", "/tmp/altid", "Path for filesystem")
	srv      = flag.String("s", "http", "Name to use for service")
	homepage = flag.String("home", "google.com", "default homepage")
	setup    = flag.Bool("conf", false, "Run configuration")
	debug    = flag.Bool("d", false, "enable debug logging")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	conf := &struct {
		Log types.Logdir
	}{"none"}

	if *setup {
		if e := config.Create(conf, *srv, "", *debug); e != nil {
			log.Fatal(e)
		}

		os.Exit(0)
	}

	if e := config.Marshal(conf, *srv, "", *debug); e != nil {
		log.Fatal(e)
	}

	ctx, cancel := context.WithCancel(context.Background())
	b, err := newBrowser(cancel)
	if err != nil {
		log.Fatal(err)
	}

	c, err := fs.CreateCtlFile(ctx, b, string(conf.Log), *mtpt, *srv, "document", *debug)
	if err != nil {
		log.Fatal(err)
	}

	c.SetCommands(Commands...)

	// initialize homepage
	c.CreateBuffer(*homepage, "document")
	fetchSite(c, *homepage, "https://"+*homepage)

	defer c.Cleanup()

	if e := c.Listen(); e != nil {
		log.Fatal(e)
	}
}
