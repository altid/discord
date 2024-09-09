package session

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"atld.ca/libs/markup"
	"altd.ca/libs/services/commander"
	"altd.ca/libs/services/controller"
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
	Client   *discordgo.Session
	ctrl     controller.Controller
	Defaults *Defaults
	Verbose  bool
	debug    func(ctlItem, ...any)
}

type Defaults struct {
	Address string       `altid:"address,prompt:IP address of Discord server"`
	Auth    string	     `altid:"auth,prompt:Authentication method to use:,pick:factotum|password"`
	SSL     string       `altid:"ssl,prompt:SSL mode,pick:none|simple|certificate"`
	User    string       `altid:"user,no_prompt"`
	Logdir  string       `altid:"logdir,no_prompt"`
	TLSCert string       `altid:"tlscert,no_prompt"`
	TLSKey  string       `altid:"tlskey,no_prompt"`
}

func (s *Session) Parse(ctx context.Context) {
	s.debug = func(ctlItem, ...any) {}
	s.ctx, s.cancel = context.WithCancel(ctx)

	if s.Verbose {
		s.debug = ctlLogging
	}
}

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
	// TODO: oauth2 instead, uses a Bearer token
	// -conf will spit out a link, click it, authorize it
	client, err := discordgo.New("Bearer " + s.Defaults.Auth)
	if err != nil {
		return err
	}
	// Top Level
	client.AddHandler(s.ready)
	client.AddHandler(s.userUpdate)
	client.AddHandler(s.resumed)
	client.AddHandler(s.webhooksUpdate)
	client.AddHandler(s.rateLimit) //discordgo
	client.AddHandler(s.connect) //discordgo
	client.AddHandler(s.disconnect) //discordgo
	// Messages
	client.AddHandler(s.msgCreate)
	client.AddHandler(s.msgUpdate)
	client.AddHandler(s.msgDelete)
	//client.AddHandler(s.msgReactionAdd)
	//client.AddHandler(s.msgRemove)
	//client.AddHandler(s.msgRemoveAll)
	//client.AddHandler(s.msgDeleteBulk)
	// Channels
	//client.AddHandler(s.chanPinsUpdate)
	client.AddHandler(s.chanCreate)
	client.AddHandler(s.chanUpdate)
	client.AddHandler(s.chanDelete)
	// Threads
	//client.AddHandler(s.threadCreate)
	//client.AddHandler(s.threadUpdate)
	//client.AddHandler(s.threadDelete)
	//client.AddHandler(s.threadListSync)
	//client.AddHandler(s.threadMemberUpdate)
	//client.AddHandler(s.threadMembersUpdate)
	// Guilds
	client.AddHandler(s.guildDelete)
	client.AddHandler(s.guildUpdate)
	client.AddHandler(s.guildMemNew)
	client.AddHandler(s.guildMemBye)
	client.AddHandler(s.guildMemUpd)
	//client.AddHandler(s.guildMembersChunk)
	//client.AddHandler(s.guildRoleCreate)
	//client.AddHandlers(s.guildRoleAdd)
	//client.AddHandlers(s.guildRoleDelete)
	//client.AddHandlers(s.guildEmojisUpdate)
	// Are these needed for DM, etc?
	// InviteCreate
	// InviteDelete
	
	s.debug(ctlSucceed, "registered client")
	s.ctrl = c
	s.Client = client
	if e := s.Client.Open(); e != nil {
log.Println(e)
		return e
	}
log.Println("REturning from start")
	return nil
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
			// Write our message to the buffer and the network, since they don't come back
			msg := m.String()
			_, err = s.Client.ChannelMessageSend(cid, msg)
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
	l := log.New(os.Stdout, "discord ", 0)

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
		l.Printf("command name=\"%s\" heading=\"%d\" args=\"%s\" from=\"%s\"", m.Name, m.Heading, m.Args, m.From)
	case ctlMsg:
		m := args[0].(*commander.Command)
		line := strings.Join(m.Args, " ")
		l.Printf("%s: data=\"%s\"\n", m.Name, line)
	case ctlErr:
		l.Printf("error: err=\"%v\"\n", args[0])
	}
}
