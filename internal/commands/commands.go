package commands

import "altd.ca/libs/services/commander"

var Commands = []*commander.Command{
	{
		Name:        "action",
		Alias:       []string{"me", "act"},
		Description: "Send an emote to the server",
		Heading:     commander.ActionGroup,
	},
	{
		Name:        "dance",
		Description: "Dance!",
		Heading:     commander.ActionGroup,
	},
	{
		Name:        "nick",
		Description: "Set new nickname",
		Args:        []string{"<name>"},
		Heading:     commander.DefaultGroup,
	},
	{
		Name:        "edit",
		Alias:       []string{"s"},
		Description: "Edit a previous message with regex (will edit the most recent message that matches)",
		Heading:     commander.DefaultGroup,
	},
	{
		Name:        "create",
		Description: "Create a channel within the current guild",
		Heading:     commander.DefaultGroup,
		Args:        []string{"<name>"},
	},
	{
		Name:		 "msg",
		Heading: 	 commander.DefaultGroup,
		Description: "Send a message to user",
		Args:    	 []string{"<name>", "<msg>"},
		Alias:   	 []string{"query", "m", "q"},
	},
}
