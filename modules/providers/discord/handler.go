package discord

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
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

	defer ctx.Request().Body.Close()
	var data Data
	if err := json.NewDecoder(ctx.Request().Body).Decode(&data); err != nil {
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
