package conferrencing

import (
	"fmt"
	"os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/livekit/protocol/auth"
	_ "github.com/livekit/server-sdk-go"
	"telehealers.in/router/models"
	opn "telehealers.in/router/restapi/operations/conferrencing"
)

const (
	ROOM = "ROOM-TEST"
)

var (
	//Constant Messages
	LIVE_KIT_NOT_WORKING = "Error: Room creation error."
	API_KEY              = os.Getenv("LIVEKIT_API_KEY")
	API_SECRET           = os.Getenv("LIVEKIT_API_SECRET")
)

func GetAccessToken(params opn.GetRoomAccessTokenParams) middleware.Responder {
	token, err := getJoinToken(API_KEY, API_SECRET, params.Room, params.ID)
	if err != nil {
		fmt.Errorf("[CONFERRENCE]Unable to create access token:%v", err)
		return opn.NewGetRoomAccessTokenDefault(400).WithPayload(
			&models.Error{Message: &LIVE_KIT_NOT_WORKING})
	} else {
		return opn.NewGetRoomAccessTokenOK().WithPayload(token)
	}
}

func getJoinToken(apiKey, apiSecret, room, identity string) (string, error) {
	canPublish := true
	canSubscribe := true

	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         room,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscribe,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
