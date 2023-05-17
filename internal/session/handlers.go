package session

import (
	"fmt"
	"runtime"

	"github.com/altid/libs/markup"
	"github.com/bwmarrin/discordgo"
)

func (s *Session) ready(ds *discordgo.Session, event *discordgo.Ready) {
	sysname := fmt.Sprintf("Discordfs on %s", runtime.GOOS)
	usd := discordgo.UpdateStatusData{
		AFK: false,
		Status: sysname,
	}
	ds.UpdateStatusComplex(usd)
}

func (s *Session) msgCreate(ds *discordgo.Session, event *discordgo.MessageCreate) {
	c, err := ds.State.Channel(event.ChannelID)
	if err != nil {
		s.debug(ctlErr, err)
		return
	}

	s.debug(ctlSucceed, "msg Callback")
	// TODO: We could look for this in the recipients
	name := c.Name

	g, err := ds.State.Guild(event.GuildID)
	if err == nil {
		name = fmt.Sprintf("%s-%s", g.Name, c.Name)
	}

	s.debug(ctlJoin, name)
	if !s.ctrl.HasBuffer(name) {
		s.chanCreate(ds, &discordgo.ChannelCreate{Channel: c})
	}

	w, err := s.ctrl.FeedWriter(name)
	if err != nil {
		s.debug(ctlErr, err)
		return
	}

	feed := markup.NewCleaner(w)
	defer feed.Close()

	feed.WritefEscaped("%%[%s](blue): %s\n", event.Author.Username, event.Message.Content)
}

func (s *Session) msgUpdate(ds *discordgo.Session, event *discordgo.MessageUpdate) {
	// Show edits
}

func (s *Session) msgDelete(ds *discordgo.Session, event *discordgo.MessageDelete) {
	// Show deletions
}

func (s *Session) chanPins(ds *discordgo.Session, event *discordgo.ChannelPinsUpdate) {
	// Pins, eventually
}

// Use our config to flag out if we care about these events
// This event has a _lot_ of useful parts, and will be much cleaner to target than
// The original way of iterating channels + such
func (s *Session) guildCreate(ds *discordgo.Session, event *discordgo.GuildCreate) {

}

func (s *Session) chanCreate(ds *discordgo.Session, event *discordgo.ChannelCreate) {
	var name string
	switch event.Type {
	case discordgo.ChannelTypeGuildText:
		g, err := ds.State.Guild(event.GuildID)
		if err != nil {
			s.debug(ctlErr, err)
			return
		}
		name = fmt.Sprintf("%s-%s", g.Name, event.Name)
	case discordgo.ChannelTypeDM, discordgo.ChannelTypeGroupDM:
		// Use the state just in case
		c, _ := s.Client.State.Channel(event.Channel.ID)
		name = c.Name
	case discordgo.ChannelTypeGuildVoice:
		return
	}
	if e := s.ctrl.CreateBuffer(name); e != nil {
		s.debug(ctlErr, e)
		return
	}
	if tw, err := s.ctrl.TitleWriter(name); err == nil {
		fmt.Fprintf(tw, "%s\n", event.Channel.Topic)
	}
	
	s.debug(ctlSucceed, "creating buffer", name)
}

func (s *Session) chanUpdate(ds *discordgo.Session, event *discordgo.ChannelUpdate) {
	// We have members here, etc
}

func (s *Session) chanDelete(ds *discordgo.Session, event *discordgo.ChannelDelete) {
	s.ctrl.DeleteBuffer(event.Name)
}

func (s *Session) guildDelete(ds *discordgo.Session, event *discordgo.GuildDelete) {
	s.ctrl.DeleteBuffer(event.Name)
}

func (s *Session) guildUpdate(ds *discordgo.Session, event *discordgo.GuildUpdate) {
	// Guild changed - log to named server
}

func (s *Session) guildMemNew(ds *discordgo.Session, event *discordgo.GuildMemberAdd) {
	// Update nicklist
}

func (s *Session) guildMemBye(ds *discordgo.Session, event *discordgo.GuildMemberRemove) {
	// Update nicklist
}

func (s *Session) guildMemUpd(ds *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	// Update nicklist
}

func (s *Session) userUpdate(ds *discordgo.Session, event *discordgo.UserUpdate) {
	// Probably can ignore this, outside of nick logging
}
