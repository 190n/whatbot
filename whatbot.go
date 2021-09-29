package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := flag.String("t", "", "discord token")
	guildId := flag.String("g", "", "ID of the guild to modify")
	iconsDir := flag.String("i", "icons", "directory to choose an icon from")
	flag.Parse()
	if *token == "" || *guildId == "" || *iconsDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	discord, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("error with discord initialization: %v", err)
	}
	
	if err := discord.Open(); err != nil {
		log.Fatalf("error opening discord connection: %v", err)
	}

	fmt.Println("connected to discord")

	files, err := ioutil.ReadDir(*iconsDir)
	if err != nil {
		log.Fatalf("error reading `%s` directory: %v", *iconsDir, err)
	}

	rand.Seed(time.Now().UnixNano())
	filename := filepath.Join(*iconsDir, files[rand.Intn(len(files))].Name())
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading `%s`: %v", filename, err)
	}

	fmt.Printf("using `%s`\n", filename)

	mime := http.DetectContentType(bytes)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(bytes))

	_, err = discord.GuildEdit(*guildId, discordgo.GuildParams{
		Icon: dataURL,
	})
	if err != nil {
		log.Fatalf("error setting icon: %v", err)
	}
}
