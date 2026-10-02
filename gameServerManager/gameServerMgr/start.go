package gameservermgr

import (
	"context"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/BrunoBiz/Narwhal/gameServerManager/logger"
)

func (gameServer *GameServer) start() ReturnValue {
	slog.Log(context.Background(), logger.LevelFile, "[Narwhal Start]")

	// Starts a new TMUX session running the server start shell
	cmd := exec.Command("tmux", "new", "-d", "-s", gameServer.config.TMUXSessionName, gameServer.config.GameStartFilePath)
	cmd.Dir = gameServer.config.GameServerDir
	tmuxStart, err := cmd.CombinedOutput()

	// Server already running
	if err != nil && strings.Contains(string(tmuxStart), "duplicate session:") {
		return newReturnValue("start", cmd.String(), string(tmuxStart), false, true, "Server is already running", err)
	}

	// Any other error - Server failed to start
	if err != nil {
		return newReturnValue("start", cmd.String(), string(tmuxStart), false, false, "Unable to start", err)
	}

	// Server started successfully
	return newReturnValue("start", cmd.String(), string(tmuxStart), true, false, "Server started", err)
}
