package main

import "github.com/bwmarrin/discordgo"

func main() {
	parser := New(".")

	parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: "."}})
}
