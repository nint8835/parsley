package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/nint8835/parsley"
)

var bot *discordgo.Session

func main() {
	bot, _ = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	bot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	parser := parsley.New("test!")
	parser.RegisterHandler(bot)
	parser.NewCommand("hello", "Greets something.", _HelloWorldCommand)

	bot.Open()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Close()
}

type _HelloWorldArgs struct {
	Target string `default:"world" description:"Target of the greeting."`
}

func _HelloWorldCommand(message *discordgo.MessageCreate, args _HelloWorldArgs) {
	bot.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Hello %s!", args.Target))
}
