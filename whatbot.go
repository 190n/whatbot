package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	var token, guildId, iconsDir string
	flag.StringVar(&token, "t", "", "discord token")
	flag.StringVar(&guildId, "g", "", "ID of the guild to modify")
	flag.StringVar(&iconsDir, "i", "icons", "directory to choose an icon from")
	flag.Parse()
	if token == "" || guildId == "" || iconsDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error with discord initialization:", err)
		os.Exit(2)
	}
	
	err = discord.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening discord connection:", err)
		os.Exit(2)
	}

	fmt.Println("connected to discord")

	files, err := ioutil.ReadDir(iconsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading `%s` directory: %v\n", iconsDir, err)
		os.Exit(2)
	}

	rand.Seed(time.Now().UnixNano())
	filename := filepath.Join(iconsDir, files[rand.Intn(len(files))].Name())
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading `%s`: %v\n", filename, err)
		os.Exit(2)
	}

	fmt.Printf("using `%s`\n", filename)

	mime := http.DetectContentType(bytes)
	dataUrl := fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(bytes))

	_, err = discord.GuildEdit(guildId, discordgo.GuildParams{
		Icon: dataUrl,
	})
	if err != nil {
		fmt.Println("error setting icon:", err)
		os.Exit(2)
	}
}
