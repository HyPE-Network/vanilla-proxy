package modal

import (
	"encoding/json"
	"math/rand"
	"vanilla-proxy/log"
	"vanilla-proxy/proxy/player"
	"vanilla-proxy/proxy/player/form"
	"vanilla-proxy/utils"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ModalForm struct {
	title, body string
	Buttons     [2]string
	operation   func(bool, bool)
}

func CreateModalForm(title string, body string) ModalForm {
	return ModalForm{
		title:   title,
		body:    body,
		Buttons: [2]string{},
	}
}

func (m ModalForm) Encode() (packet.Packet, uint32, error) {
	id := rand.Uint32()

	data, err := m.MarshalJSON()
	if err != nil {
		return nil, 0, err
	}

	pk := &packet.ModalFormRequest{
		FormID:   id,
		FormData: data,
	}

	return pk, id, nil
}

// MarshalJSON ...
func (m ModalForm) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    form.ModalType,
		"title":   m.title,
		"content": m.body,
		"button1": m.Buttons[0],
		"button2": m.Buttons[1],
	})
}

func (m ModalForm) GetType() string {
	return form.ModalType
}

func (m ModalForm) Do(closed bool, data any) {
	m.operation(closed, data.(bool))
}

func (m *ModalForm) SetFirstButton(text string) {
	m.Buttons[0] = text
}

func (m *ModalForm) SetSecondButton(text string) {
	m.Buttons[1] = text
}

func (m *ModalForm) SetTitle(text string) {
	m.title = text
}

func (m *ModalForm) SetBody(text ...any) {
	m.body = utils.Format(text)
}

func (m *ModalForm) Send(player *player.Player) {
	pk, _, err := m.Encode()
	if err != nil {
		log.Logger.Errorln(err)
		return
	}

	player.DataPacket(pk)
}

func (m *ModalForm) SendFunc(player *player.Player, f func(bool, bool)) {
	pk, id, err := m.Encode()
	if err != nil {
		log.Logger.Errorln(err)
		return
	}

	m.operation = f
	player.PlayerData.Forms[id] = m
	player.DataPacket(pk)
}
