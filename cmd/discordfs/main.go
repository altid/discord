package main

import (
	"flag"
	"log"
	"os"

	"github.com/altid/discordfs"
)

var (
	srv		= flag.String("s", "discord", "name of service")
	addr 	= flag.String("a", "127.0.0.1:12345", "listening address")
	mdns	= flag.Bool("m", false, "enable mDNS broadcast of service")
	debug 	= flag.Bool("d", false, "enable debug printing")
	ldir	= flag.Bool("l", false, "enable logging for main buffers")
	setup	= flag.Bool("conf", false, "run configuration setup")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	if *setup {
		if e := discordfs.CreateConfig(*srv, *debug); e != nil {
			log.Fatal(e)
		}
		os.Exit(0)
	}

	discord, err := discordfs.Register(*ldir, *addr, *srv, *debug)
	if err != nil {
		log.Fatal(err)
	}

	defer discord.Cleanup()
	if *mdns {
		if e := discord.Broadcast(); e != nil {
			log.Fatal(e)
		}
	}

	if e := discord.Run(); e != nil {
		log.Fatal(e)
	}

	os.Exit(0)
}