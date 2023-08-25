package form

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	SimpleType string = "form"
	ModalType  string = "modal"
	CustomType string = "custom_form"

	ImagePath string = "path"
	ImageURL  string = "url"
)

type Form interface {
	Encode() (packet.Packet, uint32, error)
	MarshalJSON() ([]byte, error)
	GetType() string
	Do(bool, any)
}

func GetResponseData(response *packet.ModalFormResponse, formType string) (bool, any, error) {
	_, closed := response.CancelReason.Value()
	if closed {
		switch formType {
		case SimpleType:
			return true, uint(0), nil
		case ModalType:
			return true, false, nil
		case CustomType:
			return true, []string{}, nil
		}
	}

	data, ok := response.ResponseData.Value()
	if !ok {
		switch formType {
		case SimpleType:
			return true, 0, nil
		case ModalType:
			return true, false, nil
		case CustomType:
			return true, []string{}, nil
		}
	}

	switch formType {
	case SimpleType:
		var ind uint
		err := json.Unmarshal(data, &ind)
		if err != nil {
			return false, nil, err
		}
		return false, ind, nil
	case ModalType:
		var value bool
		if err := json.Unmarshal(data, &value); err != nil {
			return false, nil, fmt.Errorf("error parsing JSON as bool: %w", err)
		}
		return false, value, nil
	case CustomType:
		dec := json.NewDecoder(bytes.NewBuffer(data))
		dec.UseNumber()

		var data []string
		if err := dec.Decode(&data); err != nil {
			return false, nil, fmt.Errorf("error decoding JSON data to slice: %w", err)
		}

		return false, data, nil
	}

	return false, nil, fmt.Errorf("unknown type")
}
