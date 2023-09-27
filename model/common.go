package model

import (
	"time"
)

const (
	commandSubscribe       = "0"
	commandUnsubscribe     = "1"
	commandChat            = "2"
	commandNotify          = "3"
	commandNotifyPublic    = "4"
	commandConfigUpdate    = "96"
	commandCloseConnection = "97"
	commandcheckConnection = "68"
	DebugMod               = false
	PongWait               = 60 * time.Second
	MaxMessageSize         = 51200000
	PingPeriod             = (PongWait * 9) / 10
	TrackingUser           = "trackingUser"
)
const (
	publish_channel             = "list_room_available"                // available channel
	user_channel                = "user_channel_%s"                    //user subscribe channel list
	list_user_available         = "list_system_user_available"         // list available user
	list_conversation_available = "list_system_conversation_available" // list available conversation
	conversation_mesage         = "conversation_message_%s"
)

var List_public_channel = []string{
	"public70",
}

var public_channel = "public_channel_skyx"
var private_channel = "private_channel_skyx"

const maxRangeHash = 100

func hashString(id string) uint64 {
	hashBytes := []byte(id)
	hashValue := uint64(hashBytes[0])<<56 | uint64(hashBytes[1])<<48 |
		uint64(hashBytes[2])<<40 | uint64(hashBytes[3])<<32 |
		uint64(hashBytes[4])<<24 | uint64(hashBytes[5])<<16 |
		uint64(hashBytes[6])<<8 | uint64(hashBytes[7])
	return hashValue
}

func hashToRange(id string) int {
	hashValue := hashString(id)
	rangeValue := int(hashValue % maxRangeHash)
	return rangeValue
}
