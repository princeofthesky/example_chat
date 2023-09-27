package repository

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/princeofthesky/example_chat/model"
	"github.com/princeofthesky/example_chat/trace_log"
	"os"
	"time"
)

func removeASocketInArray(src []*model.WsConnection, connId string) []*model.WsConnection {
	index := -1
	length := len(src)
	for i := 0; i < length; i++ {
		if src[i].Id == connId {
			index = i
		}
	}
	if index >= 0 && index < length && length > 0 {
		newArray := make([]*model.WsConnection, length-1)
		copy(newArray[:index], src[:index])
		copy(newArray[index:], src[index+1:])
		return newArray
	}
	return src
}
func handleReadData(socket *SocketLive, user *model.User, conn *model.WsConnection) {
	// fmt.Println("listtem message", conn.Name)
	defer func() {
		trace_log.Logger.Info().Str("c.Id", conn.Id).Msg("Stop loop wait read")
		socket.CloseAConnectionData(user, conn)
		//END
		conn.Stop = true
		conn.Ws.Close()
	}()

	con := conn.Ws

	//Lee disable timeout for testing only
	// fmt.Println("Check conn", conn.Id, conn.Expired)
	trace_log.Logger.Info().Str("conn.Id", conn.Id).Bool("conn.Expired", conn.Expired).Msg("Check conn expire ")
	connectiontesting := os.Getenv("APP_DISABLE_CONN_TIMEOUT")
	if conn.Expired == true {
		conn.Enable = false
		// conn.notifyCloseMsg <- true
		con.SetReadDeadline(time.Now().Add(10 * time.Second))
		// conn.notifyCloseMsg <- true
	} else {
		if connectiontesting == "1" {
			con.SetReadDeadline(time.Now().Add(30 * time.Minute))
		} else {
			con.SetReadDeadline(time.Now().Add(PongWait))
		}
	}

	// con.SetPingHandler(nil)
	// con.SetPongHandler(func(string) error { con.SetReadDeadline(time.Now().Add(PongWait)); return nil })
	for {
		_, item, err := con.ReadMessage()
		if err != nil {
			// fmt.Println("ReadData func", "error read socket", err.Error())
			trace_log.Logger.Error().Err(err).Str("conn.Id", conn.Id).Msg("Error when read data from connection")
			return
		}
		var msg model.Message
		err = json.Unmarshal(item, &msg)
		if err != nil {
			// fmt.Println("ReadData func", "parse json error", err.Error())
			trace_log.Logger.Info().Err(err).Str("item", string(item)).Msg("Error when parser message ")
			return
		}
		// fmt.Println("send message ", msg.Command, msg.Content)

		//Lee update sender to message
		msg.Sender = conn.Name

		switch msg.Command {
		case commandSubscribe:
			socket.Subscribe(user, msg.Channel, conn)

			var notifyMes = model.Message{
				Command: "",
				Content: "",
				Channel: msg.Channel,
			}

			conn.AddWaitingMsg(websocket.TextMessage, notifyMes.JSONEncode())

		case commandUnsubscribe:

			socket.Unsubscribe(user, msg.Channel, conn)
			var notifyMes = model.Message{
				Command: "",
				Content: "unsubscibe success",
				Channel: msg.Channel,
			}

			conn.AddWaitingMsg(websocket.TextMessage, notifyMes.JSONEncode())
		case commandChat:
			start := time.Now().UnixMilli()
			room := socket.GetRoom(msg.Channel)
			end := time.Now().UnixMilli() - start
			if end > 10 {
				println("GetRoom too long", end)
			}
			if room != nil {
				room.SendMsgToAll(&msg)
			} else {
				println("send to room nil ", msg.Channel)

			}
		}

	}

}
