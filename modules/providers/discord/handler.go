package discord

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	c "frigga/modules/common"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kataras/iris"
)

// VerifySignature ...
func VerifySignature(ctx iris.Context) {
	hexEncodedDiscordPubkey := os.Getenv("DISCORD_PUBLIC_KEY")
	discordPubkey, _ := hex.DecodeString(hexEncodedDiscordPubkey)

	isVerified := verify(ctx, discordPubkey)
	if !isVerified {
		ctx.StatusCode(401)
		ctx.Text("invalid signature")
		return
	}

	var data RequestData
	if err := ctx.ReadJSON(&data); err != nil {
		ctx.StatusCode(500)
		ctx.Text("error decoding data")
		return
	}

	if data.Type == Ping {
		ctx.JSON(map[string]int{
			"type": 1,
		})
		return
	}

	ctx.Next()
}

func verify(ctx iris.Context, key ed25519.PublicKey) bool {
	var msg bytes.Buffer

	signature := ctx.GetHeader("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize || sig[63]&224 != 0 {
		return false
	}

	timestamp := ctx.GetHeader("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	msg.WriteString(timestamp)

	defer ctx.Request().Body.Close()
	var body bytes.Buffer

	// at the end of the function, copy the original body back into the request
	defer func() {
		ctx.Request().Body = ioutil.NopCloser(&body)
	}()

	// copy body into buffers
	_, err = io.Copy(&msg, io.TeeReader(ctx.Request().Body, &body))
	if err != nil {
		return false
	}

	return ed25519.Verify(key, msg.Bytes(), sig)
}

// EventAdapter ...
func EventAdapter(ctx iris.Context) ([]RequestData, error) {
	var request RequestData

	if err := ctx.ReadJSON(&request); err != nil {
		return []RequestData{
			request,
		}, err
	}

	return []RequestData{
		request,
	}, nil
}

// EventReplier ...
func EventReplier(message c.Message, token string, isFirst bool) error {
	client := &http.Client{}
	appId := os.Getenv("DISCORD_APP_ID")

	requestBody, _ := json.Marshal(map[string]interface{}{
		"content": message.Text,
	})
	buffer := bytes.NewBuffer(requestBody)

	var apiURL string
	var req *http.Request

	if isFirst {
		apiURL = fmt.Sprintf("https://discord.com/api/v8/webhooks/%s/%s/messages/@original", appId, token)
		req, _ = http.NewRequest(http.MethodPatch, apiURL, buffer)
	} else {
		apiURL = fmt.Sprintf("https://discord.com/api/v8/webhooks/%s/%s", appId, token)
		req, _ = http.NewRequest(http.MethodPost, apiURL, buffer)
	}

	req.Header.Set("Content-Type", "application/json")

	if _, err := client.Do(req); err != nil {
		return err
	}

	return nil
}

func ConvertDataToInlineCommand(data Data) string {
	command := "/" + data.Name

	if len(data.Options) != 0 {
		return command + " " + fmt.Sprintf("%v", data.Options[0].Value)
	}

	return command
}
