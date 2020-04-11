package service

// BroadcastMessage ...
type BroadcastMessage struct {
	Body string `json:"body" validate:"required,min=3"`
	Type string `json:"type" validate:"required,min=4"`
	Mode string `json:"mode" validate:"required,min=4"`
}
