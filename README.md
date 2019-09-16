****Alpaca Game AI****

This Project is Part of the Alpaca-Bot-Challenge where different Teams 
try to program the best AI for a simple game. Each AI receives the Gamedata from the Server. 
And then sends a turn to the Server.

For Training a new Network look at the AlpacaGameAI package. 

**Usage with Server**

 `./alpaca-game-ai -myip=[MYIP] -url=[SERVER-URL] -pCnt=[CNT] -gen=best_0.601917-89-319-winner`
 
*MYIP*: Your IP Address in the Network

*SERVER-URL*: URL of the server 

*CNT*: How many players are playing

**Example with Network:**

 `./alpaca-game-ai -myip=192.168.1.25 -url="http://192.168.1.2 -pCnt=4 -gen=best_0.601917-89-319-winner`

**Example with BaseAI:**

 `./alpaca-game-ai -myip=192.168.1.25 -url="http://192.168.1.2 -pCnt=4 -baseline`

**Usage for Testing**

This will test the Genome against the BaseAI. The game will run for 50000 rounds.
the Result of the Simulation is an Array of Scores. The last Score is the Score of the Genome

 `./alpaca-game-ai -test -pCnt=4 -gen=best_0.601917-89-319-winner`

