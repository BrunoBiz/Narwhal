package gameservermgr

import (
	"context"
	"errors"
	"example/Narwhal/gameServerManager/logger"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

func (gameServer *GameServer) stop() ReturnValue {
	slog.Log(context.Background(), logger.LevelFile, "[Narwhal Stop]")

	var returnStop ReturnValue

	slog.Info("Stopping server...")

	// Some servers may require custom shutdown sequences
	/*
		CustomShutdownSequence
		0 - Default - Shutdown command
		1 - Shutdown command followed by a N(o) commnand
	*/
	switch gameServer.config.CustomShutdownSequence {
	case 0:
		returnStop = shutdownSequenceDefault(gameServer)
	case 1:
		returnStop = shutdownSequenceOne(gameServer)
	default:
		returnStop = newReturnValue("stop", "", "", false, false, "Invalid CustomShutdownSequence", nil)
	}

	return returnStop
}

func shutdownSequenceDefault(gameServer *GameServer) ReturnValue {
	var stopServerDetails ReturnValue
	var err error
	var cmd *exec.Cmd
	var tmuxStop []byte

	slog.Log(context.Background(), logger.LevelFile, "[Narwhal Stop - Shutdown sequence 0]")
	cmd = exec.Command("tmux", "send-keys", "-t", gameServer.config.TMUXSessionName, "shutdown", "ENTER")
	tmuxStop, err = cmd.CombinedOutput()

	// Server is already stopped
	if err != nil && strings.Contains(string(tmuxStop), "no server running on") {
		return newReturnValue("stop", cmd.String(), string(tmuxStop), false, false, "Server is already stopped", err)
	}

	// Any other error
	if err != nil {
		return newReturnValue("stop", cmd.String(), string(tmuxStop), false, false, "shutdown - Script failed to run (0)", err)
	}

	// max timeout wait
	slog.Debug("Max timeout")
	currentTime := time.Now()
	timeOut := currentTime.Add(time.Duration(gameServer.config.ServerStopTimeout) * time.Second)

	// Waits until the server is offline -  Start Timer
	for {
		fmt.Printf("\rTime elapsed: %s / %ds", time.Since(currentTime).Round(time.Second), gameServer.config.ServerStopTimeout)

		stopServerDetails = gameServer.OptionSwitch("details", false)
		if !stopServerDetails.ServerOnline {
			// Server stopped
			fmt.Print("\n\n")
			break
		}

		// Time out - server did not shutdown in time
		if time.Now().After(timeOut) {
			return newReturnValue("stop", "", "", false, false, "Server timed out", errors.New("Timeout"))
		}

		time.Sleep(time.Second)
	}

	return newReturnValue("stop", cmd.String(), string(tmuxStop), true, false, "Server stopped", nil)
}

func shutdownSequenceOne(gameServer *GameServer) ReturnValue {
	var stopServerDetails ReturnValue
	var err error
	var cmd, cmdNCommand *exec.Cmd
	var tmuxStop, tmuxCapturePane []byte

	slog.Log(context.Background(), logger.LevelFile, "[Narwhal Stop - Shutdown sequence 1]")
	slog.Debug("Shutdown - Sequence 1")

	/*
		After shutdown, server prompts user to restart automatically after 30s.
		Need to send a 'N' to stop the auto-restart and keep the server off.

	*/
	cmd = exec.Command("tmux", "send-keys", "-t", gameServer.config.TMUXSessionName, "shutdown", "ENTER")
	tmuxStop, err = cmd.CombinedOutput() // send-keys has no output, fire and forget

	// Server is already stopped
	if err != nil && strings.Contains(string(tmuxStop), "no server running on") {
		return newReturnValue("stop", cmd.String(), string(tmuxStop), false, false, "Server is already stopped", err)
	}

	// Any other error
	if err != nil {
		return newReturnValue("stop", cmd.String(), string(tmuxStop), false, false, "shutdown - Script failed to run (0)", err)
	}

	// Max timeout wait
	slog.Debug("Max timeout")
	currentTime := time.Now()
	timeOut := currentTime.Add(time.Duration(gameServer.config.ServerStopTimeout) * time.Second)

	// Start timer
	for {
		fmt.Printf("\rTime elapsed: %s / %ds", time.Since(currentTime).Round(time.Second), gameServer.config.ServerStopTimeout)

		// Will need to read the stdout and wait here until 'Server will re-start *automatically* in less than 30 seconds...' shows up
		for {
			// Captures the last 100 lines from the tmux terminal
			cmd = exec.Command("tmux", "capture-pane", "-S", "-100", "-E", "-", "-p", "-t", gameServer.config.TMUXSessionName)
			tmuxCapturePane, err = cmd.CombinedOutput()

			if err != nil {
				return newReturnValue("stop", cmd.String(), string(tmuxCapturePane), false, false, "shutdown - Script failed to run (1)", err)
			}

			// When pronpted to restart the server or not, uses send-keys to send a N command and complete the server shutdown
			if strings.Contains(string(tmuxCapturePane), "Server will re-start ") {
				cmdNCommand = exec.Command("tmux", "send-keys", "-t", gameServer.config.TMUXSessionName, "n", "ENTER")
				_, err = cmdNCommand.CombinedOutput()

				if err != nil {
					return newReturnValue("stop", cmdNCommand.String(), "", false, false, "shutdown - Script failed to run (2)", err)
				}

				break
			}

			// Time out - server did not shutdown in time
			// Since this is inside a loop, need to check if it went past the timeout threshold
			if time.Now().After(timeOut) {
				return newReturnValue("stop", "", "", false, false, "Server timed out", errors.New("Timeout"))
			}

			time.Sleep(time.Second)
		}

		stopServerDetails = gameServer.OptionSwitch("details", false)
		if !stopServerDetails.ServerOnline {
			// Server stopped
			fmt.Print("\n\n")
			break
		}

		// Time out - server did not shutdown in time
		if time.Now().After(timeOut) {
			return newReturnValue("stop", "", "", false, false, "Server timed out", errors.New("Timeout"))
		}

		time.Sleep(time.Second)
	}

	return newReturnValue("stop", cmd.String(), string(tmuxStop), true, false, "Server stopped", nil)
}
