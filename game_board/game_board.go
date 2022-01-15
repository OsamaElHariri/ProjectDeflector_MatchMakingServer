package game_board

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func StartGame(playerIds []string) (string, error) {
	res, err := json.Marshal(map[string]interface{}{
		"playerIds": playerIds,
	})

	if err != nil {
		return "", err
	}

	resp, _ := http.Post("http://127.0.0.1:8080/game/game", "application/json", bytes.NewBuffer(res))
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
