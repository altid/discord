package main

import "github.com/altid/libs/fs"

var Commands = []*fs.Command{
	&fs.Command{
		Name:        "action",
		Alias:       []string{"me", "act"},
		Description: "Send an emote to the server",
		Heading:     fs.ActionGroup,
	},
	&fs.Command{
		Name:        "dance",
		Description: "Dance!",
		Heading:     fs.ActionGroup,
	},
	&fs.Command{
		Name:        "nick",
		Description: "Set new nickname",
		Args:        []string{"<name>"},
		Heading:     fs.DefaultGroup,
	},
	&fs.Command{
		Name:        "edit",
		Alias:       []string{"s"},
		Description: "Edit a previous message with regex (will edit the most recent message that matches)",
		Heading:     fs.DefaultGroup,
	},
	&fs.Command{
		Name:        "create",
		Description: "Create a channel within the current guild",
		Heading:     fs.DefaultGroup,
		Args:        []string{"<name>"},
	},
	&fs.Command{
		Name:    "msg",
		Heading: fs.DefaultGroup,
		Description: "Send a message to user",
		Args:    []string{"<name>", "<msg>"},
		Alias:   []string{"query", "m", "q"},
	},
}
