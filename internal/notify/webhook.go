package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Payload struct {
	SchemaID    int64  `json:"schema_id"`
	Impact      string `json:"impact"`
	Summary     string `json:"summary"`
	VersionFrom int    `json:"version_from"`
	VersionTo   int    `json:"version_to"`
}

func Send(url string, payload Payload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	_, err = client.Do(req)
	return err
}
