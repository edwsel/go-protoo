package main

import (
	"encoding/json"

	"github.com/cloudwebrtc/go-protoo/client"
	"github.com/cloudwebrtc/go-protoo/logger"
	"github.com/cloudwebrtc/go-protoo/peer"
	"github.com/cloudwebrtc/go-protoo/transport"
)

func JsonEncode(str string) map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		panic(err)
	}
	return data
}

type AcceptFunc func(data map[string]interface{})
type RejectFunc func(errorCode int, errorReason string)

func handleWebSocketOpen(transport *transport.WebSocketTransport) {
	logger.Infof("handleWebSocketOpen")

	peer := peer.NewPeer("aaa", transport)
	peer.On("close", func() {
		logger.Infof("peer close")
	})

	handleRequest := func(request map[string]interface{}, accept AcceptFunc, reject RejectFunc) {
		method := request["method"]
		logger.Infof("handleRequest =>  (%s) ", method)
		if method == "kick" {
			reject(486, "Busy Here")
		} else if method == "offer" {
			reject(500, "sdp error!")
		}
	}

	handleNotification := func(notification map[string]interface{}) {
		logger.Infof("handleNotification => %s", notification["method"])
	}

	handleClose := func() {
		logger.Infof("handleClose => peer (%s) ", peer.ID())
	}

	peer.On("request", handleRequest)
	peer.On("notification", handleNotification)
	peer.On("close", handleClose)

	peer.Request("login", JsonEncode(`{"username":"aaa","password":"XXXX"}`),
		func(result map[string]interface{}) {
			logger.Infof("login success: =>  %s", result)
		},
		func(code int, err string) {
			logger.Infof("login reject: %d => %s", code, err)
		})
	peer.Request("join", JsonEncode(`{"client":"aaa", "type":"sender"}`),
		func(result map[string]interface{}) {
			logger.Infof("join success: =>  %s", result)
		},
		func(code int, err string) {
			logger.Infof("join reject: %d => %s", code, err)
		})
	peer.Request("publish", JsonEncode(`{"type":"sender", "jsep":{"type":"offer", "sdp":"111111111111111"}}`),
		func(result map[string]interface{}) {
			logger.Infof("publish success: =>  %s", result)
		},
		func(code int, err string) {
			logger.Infof("publish reject: %d => %s", code, err)
		})

}

func main() {
	var ws_client = client.NewClient("wss://127.0.0.1:8443/ws?peer=aaa&room=room1", handleWebSocketOpen)
	ws_client.ReadMessage()
}