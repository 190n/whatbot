package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("hello")
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("error with discord initialization:", err)
		return
	}
	
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening discord connection:", err)
		return
	}

	_, err = discord.GuildEdit("891527327144611901", discordgo.GuildParams{
		Icon: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVQI12NgyP4PAAHZAWuZrByGAAAAAElFTkSuQmCC",
	})
	if err != nil {
		fmt.Println("error setting icon:", err)
		return
	}
}
