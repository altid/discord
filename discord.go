package discord

import (
	"context"

	"altd.ca/services/discord/internal/commands"
	"altd.ca/services/discord/internal/session"
	"altd.ca/libs/config"
	"altd.ca/libs/services"
)

type Discord struct {
	run     func() error
	session *session.Session
	name    string
	debug   bool
	ctx     context.Context
}

var defaults *session.Defaults = &session.Defaults{
	Address: "discordapp.com",
	Auth:    "password",
	SSL:     "",
	User:    "",
	TLSCert: "",
	TLSKey:  "",
}

func CreateConfig(srv string, debug bool) error {
	return config.Create(defaults, srv, "", debug)
}

func Register(srv string, fg, debug bool) (*Discord, error) {
	if e := config.Marshal(defaults, srv, "", debug); e != nil {
		return nil, e
	}

	session := &session.Session{
		Defaults: defaults,
		Verbose:  debug,
	}

	ctx := context.Background()
	session.Parse(ctx)

	d := &Discord{
		session: session,
		ctx:     ctx,
		name:    srv,
		debug:   debug,
	}

	svc, err := service.Register(ctx, srv, fg)
	if err != nil {
		return nil, err
	}

	svc.SetCallbacks(session)
	svc.SetCommands(commands.Commands)
	d.run = svc.Listen
	return d, nil
}

func (discord *Discord) Run() error {
	return discord.run()
}

func (discord *Discord) Cleanup() {
	discord.session.Quit()
}

func (discord *Discord) Session() *session.Session {
	return discord.session
}
