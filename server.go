package main

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/altid/libs/fs"
	"github.com/altid/libs/markup"
	"github.com/bwmarrin/discordgo"
)

var workdir = path.Join(*mtpt, *srv)

type server struct {
	cancel context.CancelFunc
	c      *fs.Control
	dg     *discordgo.Session
	guilds []*discordgo.Guild
}

func (s *server) Refresh(*fs.Control) error {
	return nil
}

func (s *server) Restart(*fs.Control) error {
	return nil
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
	return fmt.Errorf("link command not supported, please use open/close")
}

func (s *server) Default(c *fs.Control, cmd *fs.Command) error {
	return runCommand(s, cmd)
}

func (s *server) Quit() {
	s.dg.Close()
}

func (s *server) Handle(bufname string, l *markup.Lexer) error {
	var m strings.Builder
	for {
		i := l.Next()
		switch i.ItemType {
		case markup.EOF:
			cid, err := getChanID(s, bufname)
			if err != nil {
				return err
			}

			_, err = s.dg.ChannelMessageSend(cid, m.String())
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
