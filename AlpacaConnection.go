package alpaca_game_ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)



type Gamestate struct {
	MyTurn        bool                      `json:"my_turn"`
	OtherPlayers  []map[string]PlayerStat  `json:"other_players"`
	Hand          []Cart                `json:"hand"`
	DiscardedCard Cart                  `json:"discarded_card"`
	Score         int                   `json:"score"`
	Coins         []Coin                `json:"coins"`
	CardpileLeft  int                   `json:"discarded_card"`
}

type PlayerStat struct {
	PlayerName string `json:"player_name"`
	CardCount  int    `json:"hand_cards"`
	Coins      []Coin `json:"coins"`
	Score      int    `json:"score"`
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
	Card string `json:"card"`
}

type GameTurnFunc func(gamestate *Gamestate) GameAction


type AlpacaConnection struct {
	Url        string
	CallbackIP string
	CallbackPort int
	CallbackHandle string
	CallbackFunc GameTurnFunc
	sessionID  string
	state      int
	stateStr   string
}

func (c *AlpacaConnection) Login(name string) {

	paraMeterIn := make(map[string]string)
	paraMeterIn["name"] = name
	paraMeterIn["callbackUrl"] = "http://" + c.CallbackIP+":"+strconv.Itoa(c.CallbackPort)+ c.CallbackHandle


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
	fmt.Println(string(buf))
	idMap := make(map[string]string)
	err = json.Unmarshal(buf, &idMap)
	if err != nil {
		panic(err)
	}

	fmt.Println(idMap)

	id := idMap["player_id"]

	c.sessionID = id

}

func (c *AlpacaConnection) callback(w http.ResponseWriter, r *http.Request) {

	game := c.Pull()

	action := c.CallbackFunc(game)

	c.Push(action)

	fmt.Println(action)

}

func (c *AlpacaConnection) ReceiveCallbacks(finished chan bool) {
	http.HandleFunc(c.CallbackHandle, c.callback)
	serve := http.ListenAndServe(":"+strconv.Itoa(c.CallbackPort), nil)
	if serve != nil { panic(serve) }
	finished <- true
}

func (c *AlpacaConnection) Pull() *Gamestate {

	game := new(Gamestate)
	resp, err := http.Get(c.Url + "/alpaca" + "?id=" + c.sessionID)

	if err != nil {
		c.state = -1
		c.stateStr = "Connection login Failes"
		return nil
	}

	fmt.Printf("Status: %s\n", resp.Status)
	c.state = resp.StatusCode
	c.stateStr = resp.Status
	if resp.StatusCode != 200 {
		return nil
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &game)
	if err != nil {
		panic(err)
	}

	return game
}

func (c *AlpacaConnection) Push(action GameAction) bool {

	postJson, err := json.Marshal(action)
	if err != nil { panic(err) }
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
		fmt.Println(resp.Status)
		fmt.Println(resp.Header)
		fmt.Println(resp.Body)
		return false
	}
	return true
}