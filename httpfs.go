package main

import (
	"flag"
	"log"
	"os"

	fs "github.com/altid/fslib"
)

var (
	mtpt = flag.String("p", "/tmp/altid", "Path for filesystem")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}
	var b browser
	b = make(map[string]string)
	logdir := fs.GetLogDir("https")
	c, err := fs.CreateCtrlFile(b, logdir, *mtpt, "http", "document")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Cleanup()
	c.Listen()
}
