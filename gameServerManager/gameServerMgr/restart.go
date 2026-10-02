package gameservermgr

import (
	"context"
	"log/slog"

	"github.com/BrunoBiz/Narwhal/gameServerManager/logger"
)

func (gameServer *GameServer) restart() ReturnValue {
	slog.Log(context.Background(), logger.LevelFile, "[Narwhal Restart]")

	// First checks if the server is running
	returnDetails := gameServer.OptionSwitch("details", true)

	// SERVER RUNNING -> Stops it
	if returnDetails.ServerOnline {
		returnStop := gameServer.OptionSwitch("stop", true)

		if !returnStop.Success { // Cant stop server
			return returnStop
		}
	}

	// Start server
	returnStart := gameServer.OptionSwitch("start", false) // The last logs are already printed by the main restart function call
	return returnStart
}
