package session

import (
	"log"
	"fmt"
	"runtime"

	"altd.ca/libs/markup"
	"github.com/bwmarrin/discordgo"
)

func (s *Session) ready(ds *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Ready event: version=%d user=%s\n", event.Version, event.User.Username)
	// There's guildIDs and privatechannels an event 

	// Set status
	sysname := fmt.Sprintf("alt/discord on %s", runtime.GOOS)
	usd := discordgo.UpdateStatusData{
		AFK: false,
		Status: sysname,
	}
	ds.UpdateStatusComplex(usd)
}

func (s *Session) userUpdate(ds *discordgo.Session, event *discordgo.UserUpdate) {
 	log.Printf("UserUpdate: %s\n", event.User.Username)
	// Probably can ignore this, outside of nick logging
}

func (s *Session) resumed(ds *discordgo.Session, event *discordgo.Resumed) {
	log.Println("Resumed")
}

func (s *Session) webhooksUpdate(ds *discordgo.Session, event *discordgo.WebhooksUpdate) {
	log.Printf("Webhooks Update: %s %s\n", event.GuildID, event.ChannelID)
}

func (s *Session) rateLimit(ds *discordgo.Session, event *discordgo.RateLimit) {
	log.Println("Rate limited")
}

func (s *Session) connect(ds *discordgo.Session, event *discordgo.Connect) {
	log.Println("Connected")
}

func (s *Session) disconnect(ds *discordgo.Session, event *discordgo.Disconnect) {
	log.Println("Disconnected")
}

// This seems like it could be greatly cleaned up
func (s *Session) msgCreate(ds *discordgo.Session, event *discordgo.MessageCreate) {
	log.Println("MsgCreate")
	c, err := s.Client.State.Channel(event.Message.ChannelID)
	if err != nil {
		s.debug(ctlErr, err)
		return
	}
	name := getName(s, event)
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
	msg := event.Message.ContentWithMentionsReplaced()
	for _, user := range event.Message.Mentions {
		if user.ID == s.Client.State.User.ID {
			feed.WritefEscaped("%%[%s](red): %s\n", event.Author.Username, msg)
			if e := s.ctrl.Notification(name, event.Message.Author.Username, msg); e != nil {
				s.debug(ctlErr, e)
			}
			return
		}
	}
	if event.Message.Author.Username == s.Client.State.User.Username {
		feed.WritefEscaped("%%[%s](blue): %s\n", event.Author.Username, msg)
	} else {
		feed.WritefEscaped("%%[%s](grey): %s\n", event.Author.Username, msg)
	}
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
		g, err := s.Client.State.Guild(event.GuildID)
		if err != nil {
			s.debug(ctlErr, err)
			return
		}
		name = fmt.Sprintf("%s-%s", g.Name, event.Name)
	case discordgo.ChannelTypeDM:
		// Single DM, find the channel name by finding other recipient
		for _, recipient := range event.Recipients {
			if recipient.ID != s.Client.State.User.ID {
				name = recipient.Username
			}
		}
	case discordgo.ChannelTypeGroupDM:
		// Group channel, can have a unique name with all the recipients
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
