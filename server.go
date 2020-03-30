package main

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/altid/libs/fs"
	"github.com/altid/libs/markup"
	"github.com/bwmarrin/discordgo"
)

var workdir = path.Join(*mtpt, *srv)

type server struct {
	c      *fs.Control
	dg     *discordgo.Session
	guilds []*discordgo.Guild
}

func (s *server) Run(c *fs.Control, cmd *fs.Command) error {
	switch cmd.Name {
	case "open":
		g, err := s.dg.State.Guild(strings.Join(cmd.Args, " "))
		if err != nil {
			return err
		}

		return s.dg.State.GuildAdd(g)
	case "close":
		g, err := s.dg.State.Guild(strings.Join(cmd.Args, " "))
		if err != nil {
			return err
		}

		return s.dg.State.GuildRemove(g)
	default:
		return errors.New("command not supported")
	}
}

func (s *server) Quit() {
	s.dg.Close()
}

func (s *server) Handle(bufname string, l *markup.Lexer) error {
	var m strings.Builder

	id, err := msgID(s.dg, bufname)
	if err != nil {
		return err
	}

	for {
		i := l.Next()
		switch i.ItemType {
		case markup.EOF:
			// So I get in a channel or guild name, split into one+ paths
			_, err = s.dg.ChannelMessageSend(id, m.String())
			return err
		// TODO(halfwit) We want to allow markup as well
		case markup.ErrorText:
		case markup.URLLink, markup.URLText, markup.ImagePath, markup.ImageLink, markup.ImageText:
		case markup.ColorText, markup.ColorTextBold:
		case markup.BoldText:
		case markup.EmphasisText:
		case markup.UnderlineText:
		default:
			m.Write(i.Data)
		}
	}
}

func msgID(dg *discordgo.Session, bufname string) (string, error) {

	for _, guild := range dg.State.Guilds {
		if guild.Name != path.Dir(bufname) {
			continue
		}

		for _, ch := range guild.Channels {
			if ch.Name != path.Base(bufname) {
				continue
			}

			return ch.ID, nil
		}
	}

	for _, pm := range dg.State.PrivateChannels {
		switch pm.Type {
		case discordgo.ChannelTypeDM:
			for _, id := range pm.Recipients {
				if id.ID == dg.State.SessionID {
					continue
				}

				if id.String() == bufname {
					return pm.ID, nil
				}
			}
		case discordgo.ChannelTypeGroupDM:
			if pm.Name == bufname {
				return pm.ID, nil
			}
		}

	}

	return "", fmt.Errorf("couldn't find channel for %s", bufname)
}
