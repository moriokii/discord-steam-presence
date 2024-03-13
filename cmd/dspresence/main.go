package main

import (
	"errors"
	"os"
	"time"

	discordclient "github.com/moriokii/discord-steam-presence/pkg/discord_client"
	"github.com/moriokii/discord-steam-presence/pkg/parser"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("steam_id", 0)
	viper.SetDefault("client_id", "02")
	viper.SetDefault("interval", 5)
}

func main() {
	if err := loadConfig("./config.toml"); err != nil {
		panic(err)
	}

	client := discordclient.New(viper.GetString("client_id"))
	client.Handshake()
	sid := viper.GetInt("steam_id") - 76561197960265728
	go changesListener(client, sid)
	client.Wait()
}

func loadConfig(filename string) error {
	viper.SetConfigFile(filename)

	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := viper.WriteConfig(); err != nil {
				return err
			}
			panic("config in ./config.toml")
		}
	}

	return nil
}

func changesListener(client *discordclient.Client, sid int) {
	ticker := time.NewTicker(viper.GetDuration("interval") * time.Second)
	for {
		select {
		case <-ticker.C:
			a, b, c := parser.Parse(sid)
			if err := client.SetActivity(discordclient.Activity{
				Details: a,
				State:   b,
				Assets: discordclient.Assets{
					LargeImageID: c,
				},
			}); err != nil {
				panic(err)
			}
		}
	}
}
