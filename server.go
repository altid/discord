package main

import (
	"fmt"
	"path"
	"strings"

	cm "github.com/altid/cleanmark"
	"github.com/altid/libs/fs"
	"github.com/bwmarrin/discordgo"
)

var workdir = path.Join(*mtpt, *srv)

type server struct {
	c      *fs.Control
	dg     *discordgo.Session
	guilds []*discordgo.Guild
}

// TODO: Open and Close both need to also handle PMs
// An Open call on a hidden (from the discordfs directory) should just do a create
// if we're already connected to a given channel
func (s *server) Open(c *fs.Control, name string) error {
	g, err := s.dg.State.Guild(name)
	if err != nil {
		return err
	}
	return s.dg.State.GuildAdd(g)
}

func (s *server) Close(c *fs.Control, name string) error {
	g, err := s.dg.State.Guild(name)
	if err != nil {
		return err
	}
	return s.dg.State.GuildRemove(g)
}

func (s *server) Link(c *fs.Control, from, name string) error {
	return fmt.Errorf("link command not supported, please use open/close\n")
}

func (s *server) Default(c *fs.Control, cmd, from, m string) error {
	// TODO(halfwit) nick + edit + create(guild/channel) + msg + me
	// Create PM session
	// Send PM through Handle
	return fmt.Errorf("Unknown command %s", cmd)
}

// input is always sent down raw to the server
func (s *server) Handle(bufname string, l *cm.Lexer) error {
	var m strings.Builder
	for {
		i := l.Next()
		switch i.ItemType {
		case cm.EOF:
			cid, err := getChanID(s, bufname)
			if err != nil {
				return err
			}
			_, err = s.dg.ChannelMessageSend(cid, m.String())
			return err
		case cm.ErrorText:
		case cm.UrlLink, cm.UrlText, cm.ImagePath, cm.ImageLink, cm.ImageText:
		case cm.ColorText, cm.ColorTextBold:
		case cm.BoldText:
		case cm.EmphasisText:
		case cm.UnderlineText:
		default:
			m.Write(i.Data)
		}
	}
	return fmt.Errorf("Unknown error parsing input encountered")
}
