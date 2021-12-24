package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"astuart.co/go-sse"
	"unifiedpush.org/go/np2p_dbus/utils"
)

type Message struct {
	Token   string
	Message string
}

func Sync(deviceId string, ans chan Message) (err error) {
	events := make(chan *sse.Event)

	go func() {
		for {
			state, opened := <-events
			if state == nil {
				fmt.Println("MSGCONT")
				continue
			}

			out := Message{}
			b, _ := io.ReadAll(state.Data)

			if string(b) == `{"type":"start"}` {
				continue
			}

			err = json.Unmarshal(b, &out)

			m, _ := base64.StdEncoding.DecodeString(out.Message)
			out.Message = string(m)

			utils.Log.Debugln("MESSAGE", string(b), out.Token, out.Message)
			ans <- out

			if !opened {
				break
			}

		}
	}()

	sse.GetReq = request
	//TODO does the colon split in the library interfere with JSON???
	err = sse.Notify(NCAppPathPrefix+"/device/"+deviceId, events)
	return

}
