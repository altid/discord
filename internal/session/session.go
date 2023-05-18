package session

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/altid/libs/config/types"
	"github.com/altid/libs/markup"
	"github.com/altid/libs/service/commander"
	"github.com/altid/libs/service/controller"
	"github.com/bwmarrin/discordgo"
)

type ctlItem int

const (
	ctlJoin ctlItem = iota
	ctlPart
	ctlStart
	ctlEvent
	ctlMsg
	ctlCommand
	ctlInput
	ctlRun
	ctlSucceed
	ctlErr
)

type Session struct {
	ctx      context.Context
	cancel   context.CancelFunc
	Client *discordgo.Session
	ctrl     controller.Controller
	Defaults *Defaults
	Verbose  bool
	debug    func(ctlItem, ...any)
}

type Defaults struct {
	Address	string		 `altid:"address,prompt:IP address of Discord server"`
	Auth	types.Auth	 `altid:"auth,prompt:Authentication method to use:,pick:factotum|password"`
	SSL     string       `altid:"ssl,prompt:SSL mode,pick:none|simple|certificate"`
	User	string 		 `altid:"user,no_prompt"`
	Logdir	types.Logdir `altid:"logdir,no_prompt"`
	TLSCert string       `altid:"tlscert,no_prompt"`
	TLSKey  string       `altid:"tlskey,no_prompt"`

}

func (s *Session) Parse() {
	s.debug = func(ctlItem, ...any) {}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	if s.Verbose {
		s.debug = ctlLogging
	}
}

// Future, multiuser
func (s *Session) Connect(Username string) error {
	return nil
}

func (s *Session) Run(c controller.Controller, cmd *commander.Command) error {
	switch cmd.Name {
	case "open":
		g, err := s.Client.State.Guild(strings.Join(cmd.Args, " "))
		if err != nil {
			return err
		}

		return s.Client.State.GuildAdd(g)
	case "close":
		g, err := s.Client.State.Guild(strings.Join(cmd.Args, " "))
		if err != nil {
			return err
		}

		return s.Client.State.GuildRemove(g)
	default:
		return errors.New("command not supported")
	}
}

func (s *Session) Start(c controller.Controller) error {
	// We would like to do this any other way ideally
	// but this saves us many allocations on using a channel receiver
	s.ctrl = c
	// TODO: oauth2 token?
	client, err := discordgo.New(s.Defaults.Auth.String())
	if err != nil {
		return err
	}
	// Create a buttload of handlers here
	client.AddHandler(s.ready)
	client.AddHandler(s.msgCreate)
	client.AddHandler(s.msgUpdate)
	client.AddHandler(s.msgDelete)
	client.AddHandler(s.chanPins)
	client.AddHandler(s.chanCreate)
	client.AddHandler(s.chanUpdate)
	client.AddHandler(s.chanDelete)
	client.AddHandler(s.guildDelete)
	client.AddHandler(s.guildUpdate)
	client.AddHandler(s.guildMemNew)
	client.AddHandler(s.guildMemBye)
	client.AddHandler(s.guildMemUpd)
	client.AddHandler(s.userUpdate)
	s.debug(ctlSucceed, "registered client")

	s.Client = client
	if e := s.Client.Open(); e != nil {
		return e
	}

	for _, guild := range client.State.Guilds {
		for _, room := range guild.Channels {
			// We only really care about text rooms for our needs
			if room.Type == discordgo.ChannelTypeGuildText {
				name := fmt.Sprintf("%s-%s", guild.Name, room.Name)
				s.ctrl.CreateBuffer(name)
				if tw, e := s.ctrl.TitleWriter(name); e == nil {
					fmt.Fprintf(tw, "%s\n", room.Topic)
				}
			}
		}
	}

	return nil
}

func (s *Session) Listen(c controller.Controller) {
	if e := s.Start(c); e != nil {
		log.Fatal(e)
	}

	<-s.ctx.Done()
}

func (s *Session) Command(cmd *commander.Command) error {
	return s.Run(s.ctrl, cmd)
}

func (s *Session) Quit() {
	s.Client.Close()
	s.cancel()
}

func (s *Session) Handle(bufname string, l *markup.Lexer) error {
	var m strings.Builder
	for {
		i := l.Next()
		switch i.ItemType {
		case markup.EOF:
			cid, err := getChanID(s, bufname)
			if err != nil {
				return err
			}

			 _, err = s.Client.ChannelMessageSend(cid, m.String())
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

func ctlLogging(ctl ctlItem, args ...any) {
	l := log.New(os.Stdout, "discordfs ", 0)

	switch ctl {
	case ctlSucceed:
		l.Printf("%s succeeded\n", args[0])
	case ctlJoin:
		l.Printf("join: target=\"%s\"\n", args[0])
	case ctlStart:
		l.Printf("start: addr=\"%s\", port=%d\n", args[0], args[1])
	case ctlRun:
		l.Println("connected")
	case ctlPart:
		l.Printf("part: target=\"%s\"\n", args[0])
	case ctlEvent:
		l.Printf("event: data=\"%s\"\n", args[0])
	case ctlInput:
		l.Printf("input: data=\"%s\" bufname=\"%s\"", args[0], args[1])
	case ctlCommand:
		m := args[0].(*commander.Command)
		l.Printf("command name=\"%s\" heading=\"%d\" sender=\"%s\" args=\"%s\" from=\"%s\"", m.Name, m.Heading, m.Sender, m.Args, m.From)
	case ctlMsg:
		m := args[0].(*commander.Command)
		line := strings.Join(m.Args, " ")
		l.Printf("%s: data=\"%s\"\n", m.Name, line)
	case ctlErr:
		l.Printf("error: err=\"%v\"\n", args[0])
	}
}
