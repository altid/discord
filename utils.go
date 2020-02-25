package main

import (
	"errors"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/altid/libs/fs"
	"github.com/bwmarrin/discordgo"
)

// Returns the chan ID and guild ID or an error
func getChanID(s *server, bufname string) (string, error) {
	name := path.Base(bufname)
	for _, g := range s.guilds {
		if !strings.HasPrefix(name, g.Name) {
			continue
		}
		ch := strings.TrimPrefix(name, g.Name+"-")
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
				if user.Username == name {
					return c.ID, nil
				}
			}
		case discordgo.ChannelTypeGroupDM:
			if c.Name == name {
				return c.ID, nil
			}
		}
	}
	return "", errors.New("No such guild/channel")
}

func errorWrite(c *fs.Control, err error) {
	ew, err := c.ErrorWriter()
	if err != nil {
		log.Fatal(err)
	}

	defer ew.Close()

	fmt.Fprintf(ew, "discordfs: %s\n", err)
}
