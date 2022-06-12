package util

import (
	"encoding/json"
	"fmt"
)

func MarshalPayloadToJson(payload any) (string, error) {

	if payload == nil {
		return "", nil
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("cannot serialize event payload: %w", err)
	}
	return string(payloadBytes), nil
}
