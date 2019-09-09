package alpaca_game_ai

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Gamestate struct {
	MyTurn        bool                    `json:"my_turn"`
	OtherPlayers  []map[string]PlayerStat `json:"other_players"`
	Hand          []Cart                  `json:"hand"`
	DiscardedCard Cart                    `json:"discarded_card"`
	Score         int                     `json:"score"`
	Coins         []Coin                  `json:"coins"`
	CardpileLeft  int                     `json:"cardpile_cards"`
	PlayersLeft   int                     `json:"players_left"`
}

type PlayerStat struct {
	PlayerName string `json:"player_name"`
	CardCount  int    `json:"hand_cards"`
	Coins      []Coin `json:"coins"`
	Score      int    `json:"score"`
	InGame     bool
}

type Cart struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Coin struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type GameAction struct {
	Action string `json:"action"`
	//DROP CARD
	//DRAW CARD
	//LEAVE ROUND
	Card  string `json:"card"`
	Reset bool
}

type GameTurnFunc func(gamestate *Gamestate, trunNr int) *GameAction

type AlpacaConnection struct {
	Url            string
	CallbackIP     string
	CallbackPort   int
	CallbackHandle string
	CallbackFunc   GameTurnFunc
	callbackCnt    int
	sessionID      string
	state          int
	stateStr       string
	LastScore      int
	GameEnd        chan bool
}

func (c *AlpacaConnection) Login(name string) {
	paraMeterIn := make(map[string]string)
	paraMeterIn["name"] = name
	paraMeterIn["callbackUrl"] = "http://" + c.CallbackIP + ":" + strconv.Itoa(c.CallbackPort) + c.CallbackHandle

	// Make this JSON
	postJson, err := json.Marshal(paraMeterIn)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Sending JSON string '%s'\n", string(postJson))

	//Send login JSON
	resp, err := http.Post(c.Url+"/join", "application/json", bytes.NewBuffer(postJson))

	if err != nil {
		c.state = -1
		c.stateStr = "Connection login Failes"
		panic("Connection login Failes")
	}

	c.state = resp.StatusCode
	c.stateStr = resp.Status

	if resp.StatusCode != 200 {
		panic("Login Rejected")
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(buf))
	idMap := make(map[string]string)
	err = json.Unmarshal(buf, &idMap)
	if err != nil {
		panic(err)
	}

	//fmt.Println(idMap)

	id := idMap["player_id"]

	c.sessionID = id

	//fmt.Println("Connected")

}

func (c *AlpacaConnection) callback(w http.ResponseWriter, r *http.Request) {
	game := c.Pull()

	if game == nil {
		return
	}

	c.callbackCnt++
	action := c.CallbackFunc(game, c.callbackCnt)

	if action != nil {
		c.Push(*action)
	}

}

func (c *AlpacaConnection) ReceiveCallbacks() {
	http.HandleFunc(c.CallbackHandle, c.callback)
}

func (c *AlpacaConnection) Pull() *Gamestate {

	game := new(Gamestate)
	resp, err := http.Get(c.Url + "/alpaca" + "?id=" + c.sessionID)

	if err != nil {
		c.state = -1
		c.stateStr = "Connection login Failes"
		return nil
	}

	//fmt.Printf("Status: %s\n", resp.Status)
	c.state = resp.StatusCode
	c.stateStr = resp.Status
	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			c.GameEnd <- true
		}
		return nil
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &game)
	//fmt.Println(game.DiscardedCard)
	if err != nil {
		panic(err)
	}

	c.LastScore = game.Score
	return game
}

func (c *AlpacaConnection) Reset() {
	c.callbackCnt = 0
}

func (c *AlpacaConnection) Push(action GameAction) bool {

	postJson, err := json.Marshal(action)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Sending JSON string '%s'\n", string(postJson))
	postContent := bytes.NewBuffer(postJson)
	resp, err := http.Post(c.Url+"/alpaca"+"?id="+c.sessionID, "application/json", postContent)

	if err != nil {
		c.state = -1
		c.stateStr = "POST ERROR"
		return false
	}

	c.state = resp.StatusCode
	c.stateStr = resp.Status

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			c.GameEnd <- true
		}
		//fmt.Println(resp.Status)
		//fmt.Println(resp.Header)
		//fmt.Println(resp.Body)
		return false
	}
	return true
}
