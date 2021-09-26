package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error with discord initialization:", err)
		return
	}
	
	err = discord.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening discord connection:", err)
		return
	}

	fmt.Println("connected to discord")

	files, err := ioutil.ReadDir("icons")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading `icons` directory:", err)
	}

	rand.Seed(time.Now().UnixNano())
	filename := "icons/" + files[rand.Intn(len(files))].Name()
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading `%s`: %v\n", filename, err)
		return
	}

	fmt.Printf("using `%s`\n", filename)

	mime := http.DetectContentType(bytes)
	dataUrl := fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(bytes))

	fmt.Println("made data url")

	_, err = discord.GuildEdit("891527327144611901", discordgo.GuildParams{
		Icon: dataUrl,
	})
	if err != nil {
		fmt.Println("error setting icon:", err)
		return
	}
}
