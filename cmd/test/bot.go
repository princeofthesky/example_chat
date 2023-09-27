package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gorilla/websocket"
	"github.com/princeofthesky/example_chat/token"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var rooms = []string{}

var userStatus = map[int]bool{}
var mu sync.Mutex

func createMessageId(room string, sender int) string {
	return room + strconv.Itoa(sender) + strconv.FormatInt(time.Now().UnixMicro(), 10)
}
func createMessageText(room string, sender int, rand int) string {
	text := room + strconv.Itoa(sender) + strconv.Itoa(rand)
	md5byte := md5.Sum([]byte(text))
	return "Msg:" + hex.EncodeToString(md5byte[:])
}

func createUser(userId int) {
	println("start user ", userId)
	mu.Lock()
	userStatus[userId] = true
	mu.Unlock()
	defer func() {
		mu.Lock()
		userStatus[userId] = false
		mu.Unlock()
		println("stop users ", userId)
	}()
	userToken, _ := token.GenerateToken(uint(userId))
	err := token.TokenStringisValid(userToken)
	if err != nil {
		println("verify token error", err.Error())
	}
	md5byte := md5.Sum([]byte(strconv.Itoa(userId)))
	md5UserIdHex := hex.EncodeToString(md5byte[:])
	md5byte = md5.Sum([]byte(strconv.Itoa(userId)))
	md5TokenHex := hex.EncodeToString(md5byte[:])
	url := "ws://localhost:8066/ws/" + md5UserIdHex +
		"/" + md5TokenHex + "?token=" + userToken
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		println("error when connect to server:", err.Error())
		return
	}
	stop := false
	defer c.Close()
	// send subcribe
	go func() {
		for !stop {
			_, _, err := c.ReadMessage()
			if err != nil {
				println("error when read message ", "userId", userId, "err", err.Error())
				stop = true
				return
			}
			//println(" user Id ", userId, " received Msg", string(payload))
		}
	}()

	time.Sleep(1 * time.Second)
	room := rooms[userId%len(rooms)]
	//for _, room := range rooms {
	err = c.WriteMessage(websocket.TextMessage, []byte("{\"id\":\""+createMessageId(room, userId)+
		"\",\"channel\":\""+room+
		"\",\"command\":\"0\",\"created_at\":\"\",\"sender\":\"\"}"))
	if err != nil {
		println("err when try subscribe room", err.Error())
		stop = true
		return
	}
	//}
	time.Sleep(5 * time.Second)
	sleepTime := rand.Int() % 5

	for !stop {
		if userId < 100 {
			//room := rooms[rand.Int()%len(rooms)]
			messageText := createMessageText(room, userId, sleepTime)
			err := c.WriteMessage(websocket.TextMessage, []byte("{\"id\":\""+createMessageId(room, userId)+
				"\",\"content\":\""+messageText+
				"\",\"channel\":\""+room+"\",\"command\":\"2\",\"created_at\":\"\",\"sender\":\"\"}"))
			if err != nil {
				println("err when send message from ", userId, " to room ", room, err.Error())
				break
			}
		} else {
			messageText := createMessageText(room, userId, sleepTime)
			err := c.WriteMessage(websocket.TextMessage, []byte("{\"id\":\""+createMessageId(room, userId)+
				"\",\"content\":\""+messageText+
				"\",\"channel\":\""+room+"\",\"command\":\"-1\",\"created_at\":\"\",\"sender\":\"\"}"))
			if err != nil {
				println("err when send message from ", userId, " to room ", room, err.Error())
				break
			}
		}
		//println("success send message ", messageText, "from user", userId, "room", room)

		sleepTime = 15
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

}
func makeRoom(max int) {
	for i := 0; i < max; i++ {
		md5byte := md5.Sum([]byte(strconv.Itoa(i)))
		md5RoomHex := hex.EncodeToString(md5byte[:])
		rooms = append(rooms, md5RoomHex)
	}
}
func main() {
	os.Setenv("TOKEN_HOUR_LIFESPAN", "1000000")
	os.Setenv("JWT_SECRET_KEY", "Ftb5nhj8fc6TfXbP")
	makeRoom(3)
	maxUser := 2000

	for {
		count := 0
		for i := 0; i < maxUser; i++ {
			status := false
			mu.Lock()
			status = userStatus[i]
			mu.Unlock()
			if status == false {
				go createUser(i)
				//time.Sleep(2 * time.Millisecond)
				count++
			}
		}
		println("start again ", count, " user")
		time.Sleep(1 * time.Second)
	}
}
