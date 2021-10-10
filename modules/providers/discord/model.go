package discord

import (
	"fmt"
	"time"
)

// Name ...
const Name = "discord"

type InteractionType int

const (
	_ InteractionType = iota
	Ping
	ApplicationCommand
)

type InteractionResponseType int

const (
	_ InteractionResponseType = iota
	Pong
	Acknowledge
	ChannelMessage
	ChannelMessageWithSource
	AcknowledgeWithSource
)

type InteractionResponseFlags int64

const Ephemeral InteractionResponseFlags = 1 << 6

type Data struct {
	Options []ApplicationCommandInteractionDataOption `json:"options"`
	Name    string                                    `json:"name"`
	ID      string                                    `json:"id"`
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	PublicFlags   int64  `json:"public_flags"`
}

type RequestData struct {
	Type   InteractionType `json:"type"`
	Token  string          `json:"token"`
	Member struct {
		User         User      `json:"user"`
		Roles        []string  `json:"roles"`
		PremiumSince time.Time `json:"premium_since"`
		Permissions  string    `json:"permissions"`
		Pending      bool      `json:"pending"`
		Nick         string    `json:"nick"`
		Mute         bool      `json:"mute"`
		JoinedAt     time.Time `json:"joined_at"`
		IsPending    bool      `json:"is_pending"`
		Deaf         bool      `json:"deaf"`
	} `json:"member"`
	User          User   `json:"user"`
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	GuildID       string `json:"guild_id"`
	Data          Data   `json:"data"`
	ChannelID     string `json:"channel_id"`
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                    `json:"name"`
	Value   interface{}                               `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

func (data *Data) GetInlineCommand() string {
	command := "/" + data.Name

	if len(data.Options) != 0 {
		return command + " " + fmt.Sprintf("%v", data.Options[0].Value)
	}

	return command
}

func (r *RequestData) GetId() string {
	if r.GuildID != "" {
		return r.Member.User.ID
	}
	return r.User.ID
}
