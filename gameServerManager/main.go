package main

import (
	"context"
	gameservermgr "example/Narwhal/gameServerManager/gameServerMgr"
	"example/Narwhal/gameServerManager/logger"
	"example/Narwhal/gameServerManager/util"
	"log/slog"
	"os"
)

func main() {
	// Create/Load the log file and set as the default logger for SLOG
	err := logger.LoadLogger()
	if err != nil {
		slog.Error("ERR - Cannot load/create log file: " + err.Error())
		return
	}

	config, err := util.LoadConfig("/home/gameserver")
	slog.Log(context.Background(), logger.LevelFile, "[Starting Narwhal] - Loading config...")
	if err != nil {
		slog.Log(context.Background(), logger.LevelFile, "[Starting Narwhal] - Could not load from config: "+err.Error())
		return
	}

	if len(os.Args) > 1 {
		gameServer := gameservermgr.NewGameServer(config)
		gameServer.OptionSwitch(os.Args[1], true)
	} else {
		slog.Log(context.Background(), logger.LevelFile, "[Starting Narwhal] - Empty os.Args")
		return
	}
}
