package bots

import (
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/utils"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Bot struct {
	Address     string
	XUID        string
	DisplayName string
	Conn        *minecraft.Conn
	Connected   bool
}

func NewBot(config utils.Config) *Bot {
	bot := &Bot{
		Address:     config.Connection.RemoteAddress,
		XUID:        config.Bot.XUID,
		DisplayName: config.Bot.DisplayName,
	}

	if err := bot.spawn(); err != nil {
		log.Logger.Errorln("Bot:", err)
		return nil
	}

	return bot
}

func (bot *Bot) spawn() error {
	serverConn, err := minecraft.Dialer{
		KeepXBLIdentityData: true,
		IdentityData: login.IdentityData{
			XUID:        bot.XUID,
			DisplayName: bot.DisplayName,
		},
	}.Dial("raknet", bot.Address)
	if err != nil {
		return err
	}

	bot.Conn = serverConn

	if err := serverConn.DoSpawn(); err != nil {
		return err
	}

	log.Logger.Debugln("Bot spawned with username", bot.DisplayName)

	bot.Connected = true
	bot.SendCommand("gamemode c")
	bot.SendCommand("tp 10000 200 10000")

	go func() {
		for {
			_, err := serverConn.ReadPacket()
			if err != nil {
				break
			}
		}

		bot.Close()
	}()

	return nil
}

func (bot *Bot) SendCommand(command string) error {
	if bot.Connected {
		cpk := &packet.CommandRequest{
			CommandLine: "/" + command,
			CommandOrigin: protocol.CommandOrigin{
				Origin:         0,
				UUID:           uuid.New(),
				RequestID:      "",
				PlayerUniqueID: 0,
			},
			Internal: false,
		}

		return bot.Conn.WritePacket(cpk)
	} else {
		if err := bot.spawn(); err != nil {
			log.Logger.Errorln("Bot:", err)
			return err
		}

		return bot.SendCommand(command)
	}
}

func (bot *Bot) Close() {
	if err := bot.Conn.Close(); err != nil {
		log.Logger.Errorln("Bot:", err)
	}
	log.Logger.Debugln("Bot despawned")
	bot.Connected = false
}
