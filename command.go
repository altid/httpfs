package main

import "github.com/altid/libs/fs"

var Commands = []*fs.Command{
	{
		Name:        "field",
		Description: "Write to named input field",
		Args:        []string{"<field ID>", "<Input to write>"},
		Heading:     fs.DefaultGroup,
	},
}
