package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"dont-genshin-bot/utils"
	"github.com/bwmarrin/discordgo"
)

func Start() {
	utils.LoadSettings()
	client, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}
	client.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	if err = client.Open(); err != nil {
		log.Fatalln("Error opening connection,", err)
	}
	defer client.Close()

	client.AddHandler(SettingHandler)
	client.AddHandler(DontGenshinHandler)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
