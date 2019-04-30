package main

import (
	"errors"
	"path"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Returns the chan ID and guild ID or an error
func getChanID(s *server, bufname string) (string, error) {
	name := path.Base(bufname)
	for _, g := range s.guilds {
		if ! strings.HasPrefix(name, g.Name) {
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

