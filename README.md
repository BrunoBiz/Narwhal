# Narwhal

Narwhal is a lightweight game server lifecycle manager written in Go.

It provides a consistent command-line interface for starting, stopping,
restarting, and retrieving the state of dedicated game servers running on Linux.

Narwhal was developed as part of a self-hosted Proxmox environment and is
designed to run alongside a game server inside an LXC container.

```bash
./Narwhal start
./Narwhal stop
./Narwhal restart
./Narwhal details
```

## Why Narwhal?

Different game servers and server distributions do not always behave the same
way.

They may use different startup scripts, live in different directories, require
different amounts of time to shut down, or require additional interaction
before they can terminate completely.

Narwhal places those implementation details behind a common CLI interface.

An external system therefore only needs to know:

```text
start
stop
restart
details
```

while Narwhal handles how those operations are performed for the configured
server.

This is particularly useful when Narwhal is invoked remotely by an orchestration
system such as [Orcha](https://github.com/BrunoBiz/Orcha).

## Features

- Start game servers inside detached `tmux` sessions
- Stop running game servers through configurable shutdown strategies
- Restart servers through a coordinated stop/start sequence
- Retrieve the current server state
- Configurable game server executable
- Configurable game server working directory
- Configurable shutdown timeout
- Support for custom shutdown sequences
- File and stdout logging using Go's `slog`
- Environment-based configuration using Viper
- Mock server for development and lifecycle testing
- PowerShell build and deployment scripts
- Designed for remote orchestration

## How It Works

Narwhal is deployed on the same Linux system as the game server it manages.

```text
             External Client
                    |
                    v
                  Orcha
                    |
                   SSH
                    |
                    v
          Proxmox LXC Container
                    |
                    v
                 Narwhal
                    |
                    v
                  tmux
                    |
                    v
              Game Server
```

Narwhal uses `tmux` as the runtime boundary around the game server.

When `start` is requested, Narwhal creates a detached tmux session and launches
the configured game server from its configured working directory.

Subsequent operations interact with or inspect that tmux session.

This allows Narwhal to provide the same lifecycle interface regardless of the
startup path or shutdown behavior of the configured game server.

## Commands

Narwhal expects one command-line operation.

### Start

```bash
./Narwhal start
```

Creates a new detached tmux session using the configured session name and starts
the game server from its configured working directory.

If the tmux session already exists, Narwhal reports that the server is already
running.

### Stop

```bash
./Narwhal stop
```

Sends the appropriate shutdown sequence to the running server and waits until
the server has stopped.

Narwhal will wait up to the configured `SERVER_STOP_TIMEOUT` before reporting
a timeout.

The exact shutdown behavior is controlled by
`CUSTOM_SHUTDOWN_SEQUENCE`.

### Restart

```bash
./Narwhal restart
```

Checks whether the server is currently running.

If it is running, Narwhal stops it first. Once the server is offline, Narwhal
starts it again.

If the server is already offline, Narwhal proceeds directly to startup.

### Details

```bash
./Narwhal details
```

Checks whether the configured tmux session currently exists and reports the
server as running or stopped.

> **Important:** `details` is not a game health check.
>
> Narwhal currently uses the presence of the tmux session to determine server
> state. A game server that is frozen, unhealthy, or otherwise malfunctioning
> may still be reported as running if its tmux session remains active.

## Configuration

Narwhal uses a configuration file named:

```text
mgr.env
```

The standard deployment expects it at:

```text
/home/gameserver/mgr.env
```

Example:

```env
TMUX_SESSION_NAME="minecraft"
GAME_START_PATH="'/opt/minecraft/Server/ServerStart.sh'"
GAME_DIR="/opt/minecraft/Server/"
SERVER_STOP_TIMEOUT=300
CUSTOM_SHUTDOWN_SEQUENCE=1
```

### Configuration Variables

| Variable | Required | Description |
| --- | --- | --- |
| `TMUX_SESSION_NAME` | Yes | Name of the tmux session in which the game server runs |
| `GAME_START_PATH` | Yes | Path to the executable or shell script used to start the game server |
| `GAME_DIR` | Yes | Working directory from which the game server is started |
| `SERVER_STOP_TIMEOUT` | Yes | Maximum number of seconds Narwhal waits for the server to stop |
| `CUSTOM_SHUTDOWN_SEQUENCE` | Yes | Shutdown strategy Narwhal should use |

See [`docs/configuration-reference.md`](docs/configuration-reference.md) for
the complete configuration reference.

## Shutdown Strategies

Some game servers require more than a single command to shut down cleanly.

Narwhal currently provides two shutdown sequences.

### Sequence 0 — Default

```env
CUSTOM_SHUTDOWN_SEQUENCE=0
```

Narwhal sends:

```text
shutdown
```

to the configured tmux session and then waits for the session to terminate.

```text
Narwhal
   |
   | shutdown
   v
Game Server
   |
   | shutting down...
   v
tmux session terminates
```

### Sequence 1 — Restart Confirmation

```env
CUSTOM_SHUTDOWN_SEQUENCE=1
```

This strategy is intended for server startup scripts that automatically restart
the server after shutdown unless the user explicitly declines.

Narwhal:

1. Sends `shutdown`
2. Monitors the tmux output for the restart prompt
3. Sends `n`
4. Waits for the tmux session to terminate

```text
Send shutdown
      |
      v
Server shuts down
      |
      v
Restart prompt appears
      |
      v
Send "n"
      |
      v
Server remains offline
```

Narwhal monitors the recent tmux terminal output using `capture-pane` while
waiting for the restart prompt.

Additional shutdown strategies can be added as required by other game servers.

## Logging

Narwhal uses Go's `log/slog` package.

Logs are written both to stdout and to files.

File logs are stored under:

```text
Log/
```

with filenames based on the current date:

```text
YYYY-MM-DD-narwhal.log
```

Narwhal also defines a custom `LevelFile` logging level for diagnostic messages
that should be written to the log file without appearing in normal stdout
output.

## Deployment

Narwhal is currently developed and tested using Debian 13 LXC containers
running under Proxmox VE.

The standard deployment uses a dedicated Linux user:

```text
gameserver
```

Narwhal and the managed game server are deployed under that user's environment.

A typical installation looks like:

```text
/home/gameserver/
├── Narwhal
├── mgr.env
└── <game-server files>
```

### Runtime Requirements

The target system requires:

- Linux
- `tmux`
- A configured game server
- Any runtime required by that game server, such as Java

## Deployment Script

The repository includes:

```text
deploy.ps1
```

The script:

1. Cross-compiles Narwhal for Linux AMD64
2. Copies the binary to the target machine using `scp`
3. Installs it as `/home/gameserver/Narwhal`
4. Sets the executable permission

Configure the deployment target before running it:

```powershell
$deployIP = "<container-ip>"
```

Then execute the script from PowerShell.

The complete LXC deployment procedure is documented in:

[`docs/deployment.md`](docs/deployment.md)

## Development

Narwhal is written in Go and currently targets Go 1.26.

To build the application from the repository root:

```bash
go build -C ./gameServerManager -o ../Narwhal
```

The repository also contains:

```text
mockServer/
```

This is a small executable used to simulate a game server during development.

It accepts input, responds to the `shutdown` command, waits briefly, and exits.
This allows Narwhal's start, stop, restart, timeout, and tmux behavior to be
tested without repeatedly launching a full game server.

A development deployment script is also included:

```text
build-dev.ps1
```

It builds and deploys both Narwhal and the mock server to a development
environment.

## Integration with Orcha

Narwhal can operate as a standalone CLI, but it was built to complement
[Orcha](https://github.com/BrunoBiz/Orcha), a Go-based Proxmox
orchestration API.

The two projects have separate responsibilities:

```text
                  Orcha
                    |
        Infrastructure / SSH
                    |
                    v
            Proxmox Container
                    |
                    v
                 Narwhal
                    |
          Game Server Lifecycle
                    |
                    v
              Game Server
```

**Orcha** handles infrastructure-level orchestration and remote communication.

**Narwhal** handles the lifecycle of the game server inside the target
environment.

This separation means Orcha does not need to know how a particular game server
starts or shuts down. It only needs to invoke Narwhal's common operations.

## Project Structure

```text
Narwhal/
├── gameServerManager/
│   ├── gameServerMgr/
│   │   ├── details.go
│   │   ├── gameServerMgr.go
│   │   ├── restart.go
│   │   ├── returnValue.go
│   │   ├── start.go
│   │   └── stop.go
│   ├── logger/
│   │   └── logger.go
│   ├── util/
│   │   └── config.go
│   └── main.go
│
├── mockServer/
│   └── main.go
│
├── docs/
│   ├── configuration-reference.md
│   └── deployment.md
│
├── build-dev.ps1
├── deploy.ps1
├── go.mod
└── LICENSE.txt
```

## Technologies

- Go
- tmux
- Viper
- slog
- PowerShell
- SSH / SCP
- Linux
- Proxmox LXC as the primary deployment environment

## Related Project

### Orcha

A REST API for retrieving information from a Proxmox environment and
orchestrating game servers running inside its LXC containers.

[github.com/BrunoBiz/Orcha](https://github.com/BrunoBiz/Orcha)

## License

Narwhal is licensed under the MIT License.
