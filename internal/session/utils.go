package session

import (
	"errors"
	"fmt"
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

// Returns the name of the channel, defaulting to "guest" if it is unsuccessful
func getName(s *Session, event *discordgo.MessageCreate) string {
	// We have a message, we need all the data
		g, err := s.Client.State.Guild(event.GuildID)
	if err == nil {
		c, _ := s.Client.State.Channel(event.ChannelID)
				return fmt.Sprintf("%s-%s", g.Name, c.Name)
	} else {
		// Loop through DMs, check if we have a good channel ID
		for _, item := range s.Client.State.PrivateChannels {
			if item.ID == event.Message.ChannelID {
				if item.Name != "" {
										return item.Name
				}
				if item.Topic != "" {
										return item.Topic
				}
				break
			}
		}
	}
	// Use a global lookup
	if event.Author.ID != s.Client.State.User.ID {
		user, err := s.Client.User(event.Author.ID)
		if err == nil {
						return user.Username
		}
	}
		// Fall through, we have no other names here that match; could be not in state yet
	c, _ := s.Client.State.Channel(event.Message.ChannelID)
		if len(c.Recipients) == 1 {
				for _, user := range c.Recipients {
			if user.ID != s.Client.State.User.ID {
								return user.Username
			}
		}
	}
	if c.Name != "" {
				return c.Name
	}
		return c.Topic
}

/*
func errorWrite(c *controller.Control, err error) {
	ew, _ := c.ErrorWriter()
	defer ew.Close()

	fmt.Fprintf(ew, "discordfs: %v\n", err)
}
*/