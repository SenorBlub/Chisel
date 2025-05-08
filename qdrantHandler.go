package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func SendChunkToQdrant(chunk Chunk, collection string) error {
	if len(chunk.Vector) == 0 {
		return fmt.Errorf("chunk vector is empty")
	}

	if collection == "" {
		collection = "Memory"
	}

	url := fmt.Sprintf("http://192.168.178.136:30333/collections/%s/points", collection)

	point := map[string]interface{}{
		"id":     uuid.New().String(),
		"vector": chunk.Vector,
		"payload": map[string]interface{}{
			"text":      chunk.Text,
			"origin":    chunk.Origin,
			"timestamp": chunk.Timestamp.Format(time.RFC3339),
			"tags":      chunk.Tags,
			"metadata":  chunk.Metadata,
		},
	}

	payload := map[string]interface{}{
		"points": []interface{}{point},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("qdrant insert failed: %s", resp.Status)
	}

	return nil
}

func DeletePointFromQdrant(pointID string, collection string) error {
	url := fmt.Sprintf("http://192.168.178.136:30333/collections/%s/points/delete", collection)

	bodyData := map[string]interface{}{
		"points": []string{pointID},
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		return fmt.Errorf("failed to marshal delete payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("qdrant delete failed: %s", resp.Status)
	}

	return nil
}
