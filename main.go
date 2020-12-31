package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	genshinPlayer sync.Map

	token      = os.Getenv("TOKEN")
	guildId    = os.Getenv("GUILD_ID")
	channelId  = os.Getenv("CHANNEL_ID")
	userRoleId = os.Getenv("USER_ROLE_ID")
	banRoleId  = os.Getenv("BAN_ROLE_ID")
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}
	client.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	if err = client.Open(); err != nil {
		log.Fatalln("Error opening connection,", err)
	}
	defer client.Close()

	client.AddHandler(DontGenshin)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func DontGenshin(s *discordgo.Session, e *discordgo.PresenceUpdate) {
	_, hp := genshinPlayer.Load(e.User.ID)
	var np bool
	for _, game := range e.Activities {
		if game.Type == discordgo.GameTypeGame &&
			isGenshin(game) &&
			e.GuildID == guildId {
			np = true
			if !hp {
				ban(s, e.User.ID)
			}
			msg := fmt.Sprintf("<@%s> 正在玩狗屎", e.User.ID)
			SendMessage(s, msg)
		}
	}
	if hp && !np {
		unBan(s, e.User.ID)
	}
}

func SendMessage(s *discordgo.Session, message string) {
	_, _ = s.ChannelMessageSend(channelId, message)
}

func unBan(s *discordgo.Session, userId string) {
	_ = s.GuildMemberRoleAdd(guildId, userId, userRoleId)
	_ = s.GuildMemberRoleRemove(guildId, userId, banRoleId)
	genshinPlayer.Delete(userId)
}

func ban(s *discordgo.Session, userId string) {
	_ = s.GuildMemberRoleRemove(guildId, userId, userRoleId)
	_ = s.GuildMemberRoleAdd(guildId, userId, banRoleId)
	genshinPlayer.Store(userId, struct{}{})
}

func isGenshin(game *discordgo.Game) bool {
	switch game.Name {
	case "原神":
		return true
	case "Genshin Impact":
		return true
	}
	return false
}
