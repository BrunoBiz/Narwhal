package gameservermgr

import (
	"context"
	"example/Narwhal/gameServerManager/logger"
	"log/slog"
	"os/exec"
	"strings"
)

func (gameServer *GameServer) details() ReturnValue {
	slog.Log(context.Background(), logger.LevelFile, "[Narwhal Details]")
	/*
		NOT a health check, this only validates that the TMUX session is active, the game may have crashed,
		is frozen or an error may have ocurred and it will still return as ONLINE if the session is still on!
	*/

	// Command to check if the tmux session exists
	cmd := exec.Command("tmux", "ls")
	tmuxSession, err := cmd.CombinedOutput()

	if err != nil {
		if strings.Contains(string(tmuxSession), "no server running") {
			// If there is no server running, even tho it returns an error, will not handle it as such
			return newReturnValue("details", cmd.String(), string(tmuxSession), true, false, "No server running", nil)
		} else {
			return newReturnValue("details", cmd.String(), string(tmuxSession), false, false, "tmux ls - Script failed to run", err)
		}
	}

	if strings.Contains(string(tmuxSession), gameServer.config.TMUXSessionName) {
		// tmux session exists
		return newReturnValue("details", cmd.String(), string(tmuxSession), true, true, "Server running", nil)
	} else {
		// tmux does not exists
		return newReturnValue("details", cmd.String(), string(tmuxSession), true, false, "Server offline", nil)
	}
}
