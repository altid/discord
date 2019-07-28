package main

import (
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/altid/cleanmark"
	"github.com/altid/fslib"
	"github.com/bwmarrin/discordgo"
)

func (s *server) ready(ds *discordgo.Session, event *discordgo.Ready) {
	s.guilds = event.Guilds
	sysname := fmt.Sprintf("Discordfs on %s", runtime.GOOS)
	ds.UpdateStatus(0, sysname)
}
func (s *server) msgCreate(ds *discordgo.Session, event *discordgo.MessageCreate) {
	c, err := ds.State.Channel(event.ChannelID)
	if err != nil {
		log.Print(err)
		return
	}
	name := c.Name
	g, err := ds.State.Guild(event.GuildID)
	if err == nil {
		name = fmt.Sprintf("%s-%s", g.Name, c.Name)
	}
	if ! s.c.HasBuffer(name, "feed") {
		s.chanCreate(ds, &discordgo.ChannelCreate{c})
	}
	w := s.c.MainWriter(name, "feed")
	if w == nil {
		log.Printf("Unable to create feed entry for %s\n", name)
		return
	}
	feed := cleanmark.NewCleaner(w)
	defer feed.Close()
	feed.WritefEscaped("%s: %s\n", event.Author.Username, event.Content)
}

func (s *server) msgUpdate(ds *discordgo.Session, event *discordgo.MessageUpdate) {
	// Show edits
}

func (s *server) msgDelete(ds *discordgo.Session, event *discordgo.MessageDelete) {
	// Show deletions
}

func (s *server) chanPins(ds *discordgo.Session, event *discordgo.ChannelPinsUpdate) {
	// Pins, eventually
}

// Use our config to flag out if we care about these events
// This event has a _lot_ of useful parts, and will be much cleaner to target than
// The original way of iterating channels + such
func (s *server) guildCreate(ds *discordgo.Session, event *discordgo.GuildCreate) {

}

func (s *server) chanCreate(ds *discordgo.Session, event *discordgo.ChannelCreate) {
	var name string
	switch event.Type {
	case discordgo.ChannelTypeGuildText:
		g, err := ds.State.Guild(event.GuildID)
		if err != nil {
			log.Println(err)
			return
		}
		name = fmt.Sprintf("%s-%s", g.Name, event.Name)
	case discordgo.ChannelTypeDM:
		name = event.Name
	case discordgo.ChannelTypeGroupDM:
		// For now, grab the last message and get the channel name from that
		m, err := ds.State.Message(event.LastMessageID, event.ID)
		if err != nil {
			log.Println(err)
			return
		}
		c, _ := ds.State.Channel(m.ChannelID)
		name = fmt.Sprintf("%s", c.Name)
	case discordgo.ChannelTypeGuildVoice:
		return
	}
	err := s.c.CreateBuffer(name, "feed")
	if err != nil {		
		return
	}
	input, err := fslib.NewInput(s, workdir, name)
	if err != nil {
		log.Println(err)
		return
	}
	defer s.c.Event(path.Join(workdir, name, "input"))
	go input.Start()
}

func (s *server) chanUpdate(ds *discordgo.Session, event *discordgo.ChannelUpdate) {
	// 
}

func (s *server) chanDelete(ds *discordgo.Session, event *discordgo.ChannelDelete) {
	s.c.DeleteBuffer(event.Name, "feed")
}

func (s *server) guildDelete(ds *discordgo.Session, event *discordgo.GuildDelete) {
	s.c.DeleteBuffer(event.Name, "feed")
}

func (s *server) guildUpdate(ds *discordgo.Session, event *discordgo.GuildUpdate) {
	// Guild changed - log to named server
}

func (s *server) guildMemNew(ds *discordgo.Session, event *discordgo.GuildMemberAdd) {
	// Update nicklist
}

func (s *server) guildMemBye(ds *discordgo.Session, event *discordgo.GuildMemberRemove) {
	// Update nicklist
}

func (s *server) guildMemUpd(ds *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	// Update nicklist
}

func (s *server) userUpdate(ds *discordgo.Session, event *discordgo.UserUpdate) {
	// Probably can ignore this, outside of nick logging
}
