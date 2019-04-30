package main

import (
	"fmt"
	"path"
	"strings"

	cm "github.com/altid/cleanmark"
	"github.com/altid/fslib"
	"github.com/bwmarrin/discordgo"
)

var workdir = path.Join(*mtpt, *srv)

type server struct {
	c *fslib.Control
	dg *discordgo.Session
	guilds []*discordgo.Guild
}

func (s *server) Open(c *fslib.Control, name string) error {
	g, err := s.dg.State.Guild(name)
	if err != nil {
		return err
	}
	return s.dg.State.GuildAdd(g)
}

func (s *server) Close(c *fslib.Control, name string) error {
	g, err := s.dg.State.Guild(name)
	if err != nil {
		return err
	}
	return s.dg.State.GuildRemove(g)
}

func (s *server) Link(c *fslib.Control, from, name string) error {
	return fmt.Errorf("link command not supported, please use open/close\n")
}

func (s *server) Default(c *fslib.Control, cmd, from, m string) error {
	// Nick + Edit + Create(guild/channel)
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
