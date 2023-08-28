package command

import (
	"strings"

	"github.com/HyPE-Network/vanilla-proxy/log"

	rcon "github.com/DEBANMC/valve-rcon"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type RconSender struct {
	Client rcon.Client
}

func (rs *RconSender) SendMessage(mess string) {
	rs.Client.Write(mess)
}

func InitRCON(commands map[CommandExecutor]protocol.Command, port int, pass string) {
	rconSrv := rcon.NewRCON("127.0.0.1", port, pass)
	rconSrv.SetBanList([]string{})

	rconSrv.OnCommand(func(commandMessage string, client rcon.Client) {
		sender := &RconSender{
			Client: client,
		}

		for executor, command := range commands {
			if command.Name == commandMessage {
				split := strings.Fields(commandMessage)
				err := executor.Execute(sender, append(split[:0], split[1:]...))
				if err != nil {
					sender.SendMessage(err.Error())
				}
			}
		}
	})

	log.Logger.Println("Rcon-server has been started")

	err := rconSrv.ListenAndServe()
	if err != nil {
		panic(err)
	}

	rconSrv.CloseOnProgramEnd()
}
