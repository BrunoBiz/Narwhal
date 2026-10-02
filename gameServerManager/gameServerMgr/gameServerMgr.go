package gameservermgr

import (
	"context"
	"errors"
	"example/Narwhal/gameServerManager/logger"
	"example/Narwhal/gameServerManager/util"
	"log/slog"
)

type GameServer struct {
	config util.Config
}

func NewGameServer(config util.Config) *GameServer {
	gameServer := &GameServer{
		config: config,
	}

	return gameServer
}

func (gameServer *GameServer) OptionSwitch(option string, printLogs bool) ReturnValue {
	var returnSwitch ReturnValue

	slog.Log(context.Background(), logger.LevelFile, "[Starting Narwhal] - Validade option - "+option)

	switch option {
	case "start":
		returnSwitch = gameServer.start()
	case "stop":
		returnSwitch = gameServer.stop()
	case "restart":
		returnSwitch = gameServer.restart()
	case "details":
		returnSwitch = gameServer.details()
	default:
		returnSwitch = newReturnValue("INVALID", option, "", false, false, "Invalid or empty option", errors.New("OPCAO INVALIDA/NULA/VAZIA\rOPCAO: "+option))
	}

	if printLogs {
		returnSwitch.PrintLogs()
	}

	return returnSwitch
}
