package main

import (
	"errors"
	"fmt"
	"path"

	"github.com/altid/libs/fs"
	"github.com/bwmarrin/discordgo"
)

// Returns the chan ID and guild ID or an error
func getChanID(s *server, bufname string) (string, error) {
	for _, g := range s.guilds {

		// Same guild by path
		if path.Dir(bufname) != g.Name {
			continue
		}

		ch := path.Base(bufname)
		// make sure chan exists
		for _, c := range g.Channels {
			if c.Name != ch {
				continue
			}

			return c.ID, nil
		}
	}
	// Group
	uc, _ := s.dg.UserChannels()
	for _, c := range uc {
		switch c.Type {
		case discordgo.ChannelTypeDM:
			for _, user := range c.Recipients {
				if user.Username == path.Base(bufname) {
					return c.ID, nil
				}
			}
		case discordgo.ChannelTypeGroupDM:
			if c.Name == path.Base(bufname) {
				return c.ID, nil
			}
		}
	}
	return "", errors.New("No such guild/channel")
}

func errorWrite(c *fs.Control, err error) {
	ew, _ := c.ErrorWriter()
	defer ew.Close()

	fmt.Fprintf(ew, "discordfs: %v\n", err)
}
