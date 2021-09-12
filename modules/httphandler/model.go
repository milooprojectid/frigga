package httphandler

import "frigga/modules/common"

// BroadcastMessage ...
type BroadcastMessage struct {
	Body  string              `json:"body" validate:"required,min=3"`
	Album *[]common.AlbumItem `json:"album,omitempty" validate:"min=2"`
	Type  string              `json:"type" validate:"required,min=4"`
	Mode  string              `json:"mode" validate:"required,min=4"`
}
