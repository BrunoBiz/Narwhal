The configuration file has to be named mgr.env and needs to be placed in `/home/gameserver/mgr.env`

| Variable                 	| Required 	| Description                                                                              	| Example                                  	|
|--------------------------	|----------	|------------------------------------------------------------------------------------------	|------------------------------------------	|
| TMUX_SESSION_NAME        	| Yes      	| How the TMUX session where the game server is running will be named                      	| "minecraft"                              	|
| GAME_START_PATH          	| Yes      	| Path to the executable or shell script that will start the game server                   	| "'/opt/minecraft/Server/ServerStart.sh'" 	|
| GAME_DIR                 	| Yes      	| Directory where the game server is installed                                             	| "/opt/minecraft/Server/"                 	|
| SERVER_STOP_TIMEOUT      	| Yes      	| How long, in seconds, Narwhal will wait for the server to stop before returning an error 	| 300                                      	|
| CUSTOM_SHUTDOWN_SEQUENCE 	| Yes      	| Custom shutdown sequence                                                                 	| 0                                        	|


CUSTOM_SHUTDOWN_SEQUENCE

Some servers may require custom shutdown sequences; hence, this parameter exists. Below are the current acceptable values (more may be added in the future):

- 0 -> Default behavior, will send a single `shutdown` command to the current running program;
- 1 -> Confirm shutdown, will send a single `shutdown` command and await for the program to request a restart, in this case will them send a `N` to represent a no;

Example:\
TMUX_SESSION_NAME="minecraft"\
GAME_START_PATH="'/opt/minecraft/Server/ServerStart.sh'"\
GAME_DIR="/opt/minecraft/Server/"\
SERVER_STOP_TIMEOUT=300 # 5 minutes\
CUSTOM_SHUTDOWN_SEQUENCE=1
