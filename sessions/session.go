package sessions

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)


type Session struct{
	subscribers *ws.conn[] // array of subscriber websocket connections that we write to
	chatid uuid
	destructionTime time.Time
}

type destructionTimeJsonResponseObject struct {
	TimeToDestruction time.Time `json: time_to_destruction`
}

const sessionManager = make(map[uuid]Session)

func getDestructionTimeForChat(id uuid) json {
	session, ok := sessionManager[id]
	if ok != nil {
		return nil
	}

	return destructionTimeJsonResponseObject{
		TimeToDestruction: session.destructionTime,
	}
}

func CreateSession() uuid {
	keys = sessionManager.Keys()

}