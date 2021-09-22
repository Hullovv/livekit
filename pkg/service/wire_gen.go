// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package service

import (
	"github.com/livekit/livekit-server/pkg/config"
	"github.com/livekit/livekit-server/pkg/routing"
)

// Injectors from wire.go:

func InitializeServer(conf *config.Config, currentNode routing.LocalNode) (*LivekitServer, error) {
	client, err := createRedisClient(conf)
	if err != nil {
		return nil, err
	}
	router := createRouter(client, currentNode)
	nodeSelector := CreateNodeSelector(conf)
	roomStore := createStore(client)
	roomAllocator := NewRoomAllocator(conf, router, nodeSelector, roomStore)
	roomService, err := NewRoomService(roomAllocator, roomStore, router)
	if err != nil {
		return nil, err
	}
	messageBus := createMessageBus(client)
	keyProvider, err := CreateKeyProvider(conf)
	if err != nil {
		return nil, err
	}
	notifier, err := CreateWebhookNotifier(conf, keyProvider)
	if err != nil {
		return nil, err
	}
	recordingService := NewRecordingService(messageBus, notifier)
	rtcService := NewRTCService(conf, roomAllocator, router, currentNode)
	localRoomManager, err := NewLocalRoomManager(roomStore, router, currentNode, nodeSelector, notifier, conf)
	if err != nil {
		return nil, err
	}
	server, err := NewTurnServer(conf, roomStore, currentNode)
	if err != nil {
		return nil, err
	}
	livekitServer, err := NewLivekitServer(conf, roomService, recordingService, rtcService, keyProvider, router, localRoomManager, server, currentNode)
	if err != nil {
		return nil, err
	}
	return livekitServer, nil
}

func InitializeRouter(conf *config.Config, currentNode routing.LocalNode) (routing.Router, error) {
	client, err := createRedisClient(conf)
	if err != nil {
		return nil, err
	}
	router := createRouter(client, currentNode)
	return router, nil
}
