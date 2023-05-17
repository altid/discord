package session

import (
	"errors"
	"path"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Returns the chan ID and guild ID or an error
func getChanID(s *Session, bufname string) (string, error) {
	for _, g := range s.Client.State.Guilds {
		if ! strings.HasPrefix(bufname, g.Name) {
			continue
		}
	
		// make sure chan exists
		for _, c := range g.Channels {
			if ! strings.HasSuffix(bufname, c.Name) {
				continue
			}
			return c.ID, nil
		}
	}
	for _, c := range s.Client.State.PrivateChannels {
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

	return "", errors.New("no such guild/channel")
}
/*
func errorWrite(c *controller.Control, err error) {
	ew, _ := c.ErrorWriter()
	defer ew.Close()

	fmt.Fprintf(ew, "discordfs: %v\n", err)
}
*/