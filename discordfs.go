package discordfs

import (
	"context"

	"github.com/altid/discordfs/internal/commands"
	"github.com/altid/discordfs/internal/session"
	"github.com/altid/libs/config"
	"github.com/altid/libs/mdns"
	"github.com/altid/libs/service"
	"github.com/altid/libs/service/listener"
	"github.com/altid/libs/store"
)

type Discordfs struct {
	run		func() error
	session *session.Session
	name	string
	addr	string
	debug	bool
	mdns	*mdns.Entry
	ctx		context.Context
}

var defaults *session.Defaults = &session.Defaults{
	Address:	"discordapp.com",
	Auth:		"password",
	SSL:		"",
	User:		"",
	Logdir:		"",
	TLSCert:    "",
	TLSKey:		"",
}

func CreateConfig(srv string, debug bool) error {
	return config.Create(defaults, srv, "", debug)
}

func Register(ldir bool, addr, srv string, debug bool) (*Discordfs, error) {
	if e := config.Marshal(defaults, srv, "", debug); e != nil {
		return nil, e
	}
	l, err := tolisten(defaults, addr, debug)
	if err != nil {
		return nil, err
	}
	s := tostore(defaults, ldir, debug)
	session := &session.Session{
		Defaults: defaults,
		Verbose:  debug,
	}

	session.Parse()
	ctx := context.Background()

	d := &Discordfs{
		session:	session,
		ctx:		ctx,
		name:		srv,
		addr:		addr,
		debug:		debug,
	}

	c := service.New(srv, addr, debug)
	c.WithListener(l)
	c.WithStore(s)
	c.WithContext(ctx)
	c.WithCallbacks(session)
	c.WithRunner(session)

	c.SetCommands(commands.Commands)
	d.run = c.Listen

	return d, nil
}

func (discord *Discordfs) Run() error {
	return discord.run()
}

func (discord *Discordfs) Broadcast() error {
	entry, err := mdns.ParseURL(discord.addr, discord.name)
	if err != nil {
		return err
	}
	if e := mdns.Register(entry); e != nil {
		return e
	}
	discord.mdns = entry
	return nil
}

func (discord *Discordfs) Cleanup() {
	if discord.mdns != nil {
		discord.mdns.Cleanup()
	}
	discord.session.Quit()
}

func (discord *Discordfs) Session() *session.Session {
	return discord.session
}

func tolisten(d *session.Defaults, addr string, debug bool) (listener.Listener, error) {
	//if ssh {
	//    return listener.NewListenSsh()
	//}

	if d.TLSKey == "none" && d.TLSCert == "none" {
		return listener.NewListen9p(addr, "", "", debug)
	}

	return listener.NewListen9p(addr, d.TLSCert, d.TLSKey, debug)
}

func tostore(d *session.Defaults, ldir, debug bool) store.Filer {
	if ldir {
		return store.NewLogstore(d.Logdir.String(), debug)
	}

	return store.NewRamstore(debug)
}
