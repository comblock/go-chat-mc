package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/mattn/go-colorable"
	GMMAuth "github.com/maxsupermanhd/go-mc-ms-auth"
)

var address string
var client *bot.Client
var player *basic.Player

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "> ",
		HistoryFile:            "/tmp/readline-multiline",
		DisableAutoSaveHistory: false,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	log.SetOutput(colorable.NewColorableStdout())

	mauth, err := GMMAuth.GetMCcredentials("auth.cache", "aeb3912d-6ddf-463d-bcd7-77d699ea01be")

	if err != nil {
		log.Fatal(err)
	}

	client = bot.NewClient()
	client.Auth = mauth
	player = basic.NewPlayer(client, basic.DefaultSettings)

	basic.EventsListener{
		GameStart:  onGameStart,
		ChatMsg:    onChatMsg,
		Disconnect: onDisconnect,
		Death:      onDeath,
	}.Attach(client)

	fmt.Print("Enter the server adress: ")
	fmt.Scan(&address)

	if err := client.JoinServer(address); err != nil {
		log.Fatalf("Login failed: %v\n", err)
	}

	log.Println("Sucessfully logged in")

	go func() {
		for {
			if err := client.HandleGame(); err == nil {
				panic("HandleGame should never return nil")
			} else {

				if err2 := new(bot.PacketHandlerError); errors.As(err, err2) {
					if err := new(bot.DisconnectErr); errors.As(err2, err) {
						log.Print("Disconnect: ", err)
						return
					} else {
						log.Print(err2)
					}
				} else {
					log.Fatal(err)
				}
			}
		}
	}()

	time.Sleep(500 * time.Millisecond)

	for {
		var message packet.String

		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		message = packet.String(line)

		client.Conn.WritePacket(packet.Marshal(packetid.ChatServerbound, message))
	}
}

func onDeath() error {
	log.Println("Died")
	go func() {
		err := player.Respawn()
		if err != nil {
			log.Print(err)
		} else {
			log.Print("Respawned")
		}
	}()
	return nil
}

func onGameStart() error {
	log.Println("Game start")
	return nil
}

func onChatMsg(message chat.Message, pos byte, uuid uuid.UUID) error {
	fmt.Print("\r")
	log.Println(message)
	return nil
}

func onDisconnect(reason chat.Message) error {
	return bot.DisconnectErr(reason)
}
