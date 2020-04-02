package main

import (
	"errors"
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
		case markup.ColorText, markup.ColorTextBold, markup.ColorTextStrong:
		case markup.BoldText:
		case markup.EmphasisText:
		case markup.StrongText:
		default:
			m.Write(i.Data)
		}
	}
}
