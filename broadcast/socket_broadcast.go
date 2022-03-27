package broadcast

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func SocketBroadcast(ids []string, event string, payload map[string]interface{}) error {
	res, err := json.Marshal(map[string]interface{}{
		"event":   event,
		"payload": payload,
	})

	if err != nil {
		return err
	}

	for _, id := range ids {
		go socketBroadcast(id, res)
	}
	return nil
}

func socketBroadcast(id string, payload []byte) {
	resp, err := SendPost(os.Getenv("INTERNAL_SERVICES_URL")+"/realtime/internal/notify/"+id, payload)
	if err != nil {
		resp.Body.Close()
	}
}

func SendPost(url string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer " + os.Getenv("INTERNAL_TOKEN")},
	}
	return http.DefaultClient.Do(req)
}
