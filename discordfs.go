package main

import (
	"flag"
	"log"
	"os"

	"github.com/altid/fslib"
	"github.com/bwmarrin/discordgo"
)

var (
	mtpt = flag.String("p", "/tmp/altid", "Path for filesystem")
	srv  = flag.String("s", "discord", "Name of service")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	config, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}
	s := &server{}
	dg, err := discordgo.New(config.user, config.pass)
	if err != nil {
		log.Fatal("Error initiating discord session %v", err)
	}
	dg.AddHandler(s.ready)
	dg.AddHandler(s.msgCreate)
	dg.AddHandler(s.msgUpdate)
	dg.AddHandler(s.msgDelete)
	dg.AddHandler(s.chanPins)
	dg.AddHandler(s.chanCreate)
	dg.AddHandler(s.chanUpdate)
	dg.AddHandler(s.chanDelete)
	dg.AddHandler(s.guildUpdate)
	dg.AddHandler(s.guildMemNew)
	dg.AddHandler(s.guildMemBye)
	dg.AddHandler(s.guildMemUpd)
	dg.AddHandler(s.userUpdate)
	ctrl, err := fslib.CreateCtrlFile(s, config.log, *mtpt, *srv, "feed")
	defer ctrl.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
	s.c  = ctrl
	s.dg = dg
	ctrl.CreateBuffer("server", "feed")
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()
	ctrl.Listen()
}
