package game_lobby

import (
	"encoding/json"
	"io/ioutil"
	"projectdeflector/matchmaking/broadcast"
)

func StartGame(playerIds []string) (string, error) {
	res, err := json.Marshal(map[string]interface{}{
		"playerIds": playerIds,
	})

	if err != nil {
		return "", err
	}

	resp, err := broadcast.SendPost("http://127.0.0.1:8080/game/internal/game", res)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data map[string]string
	json.Unmarshal(body, &data)

	return data["gameId"], nil
}
