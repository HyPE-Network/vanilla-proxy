package simple

import (
	"encoding/json"
	"math/rand"
	"vanilla-proxy/log"
	"vanilla-proxy/proxy/player/form"
	"vanilla-proxy/proxy/player/human"
	"vanilla-proxy/utils"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SimpleForm struct {
	title, body string
	Buttons     []form.Button
	operation   func(bool, uint)
}

func CreateSimpleForm(title string, body string) SimpleForm {
	return SimpleForm{
		title:   title,
		body:    body,
		Buttons: []form.Button{},
	}
}

func (m SimpleForm) Encode() (packet.Packet, uint32, error) {
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

func (m SimpleForm) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    form.SimpleType,
		"title":   m.title,
		"content": m.body,
		"buttons": m.Buttons,
	})
}

func (m SimpleForm) GetType() string {
	return form.SimpleType
}

func (m SimpleForm) Do(closed bool, data any) {
	m.operation(closed, data.(uint))
}

func (m *SimpleForm) SetTitle(text string) {
	m.title = text
}

func (m *SimpleForm) SetBody(text ...any) {
	m.body = utils.Format(text)
}

func (m *SimpleForm) AddButton(text string) {
	m.Buttons = append(m.Buttons, form.Button{Text: text})
}

func (m *SimpleForm) AddImageButton(text string, image form.ImageType) {
	m.Buttons = append(m.Buttons, form.Button{Text: text, Image: image})
}

func (m *SimpleForm) Send(player human.Human) {
	pk, _, err := m.Encode()
	if err != nil {
		log.Logger.Errorln(err)
		return
	}

	player.DataPacket(pk)
}

func (m *SimpleForm) SendFunc(player human.Human, f func(bool, uint)) {
	pk, id, err := m.Encode()
	if err != nil {
		log.Logger.Errorln(err)
		return
	}

	m.operation = f
	player.GetData().Forms[id] = m
	player.DataPacket(pk)
}
