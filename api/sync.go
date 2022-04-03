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
	// ignore in public api
	Type           string
	Token          string
	Message        string
	DecodedMessage []byte
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

			err = json.Unmarshal(b, &out)
			if out.Type == "start" {
				continue
			} else if out.Type == "warning" {
				utils.Log.Infoln("WARNING", out.Message)
				continue
			} else if out.Type == "" {
				continue
			}

			m, _ := base64.StdEncoding.DecodeString(out.Message)
			out.DecodedMessage = m

			utils.Log.Debugln("MESSAGE", string(b), out.Token, out.DecodedMessage)
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
