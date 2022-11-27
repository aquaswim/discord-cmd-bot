package main

import (
	"discord-cmd-bot/config"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	token      string
	configPath string
)

func init() {
	flag.StringVar(&configPath, "c", "", "Bot Config")
	token = os.Getenv("DISCORD_TOKEN")
	if token == "" {
		flag.StringVar(&token, "t", "", "Bot Token")
	}
	flag.Parse()
	if token == "" {
		panic("Need to provide token in -t or env variable DISCORD_TOKEN")
	}
}

func main() {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("error loading config,", err)
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	cmdRunner := NewCommandRunner(cfg)
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if strings.HasPrefix(m.Content, "%run ") {
			dcCommand := strings.TrimLeft(m.Content, "%run ")
			if cmdRunner.HasCommand(dcCommand) {
				output, err := cmdRunner.RunCommand(dcCommand)
				if err != nil {
					_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("result of %s is error: ```%s```", dcCommand, err))
					if err != nil {
						fmt.Println("send message error", err)
					}
					return
				}
				//_, err = s.ChannelMessageSendComplex(m.ChannelID, fmt.Sprintf("result of %s ```%s```", dcCommand, output))
				_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title:       fmt.Sprintf("result of command %s", dcCommand),
					Description: fmt.Sprintf("```%s```", output),
				})
				if err != nil {
					fmt.Println("send message error", err)
				}
				return
			}
		}
	})

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
