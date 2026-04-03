package discordutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func ValidateBase64MessageJSON(base64Str string) (discord.MessageCreate, error) {
	if base64Str == "" {
		return discord.MessageCreate{}, fmt.Errorf("base64 string está vazia")
	}

	jsonBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return discord.MessageCreate{}, fmt.Errorf("base64 inválido: %w", err)
	}

	var rawJSON map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &rawJSON); err != nil {
		return discord.MessageCreate{}, fmt.Errorf("JSON inválido: %w", err)
	}

	var msg discord.MessageCreate
	if err := json.Unmarshal(jsonBytes, &msg); err != nil {
		return discord.MessageCreate{}, fmt.Errorf("estrutura de mensagem inválida: %w", err)
	}

	return msg, nil
}

func EncodeMessage(msg discord.MessageCreate) (string, error) {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

func DecodeMessage(base64Str string) (discord.MessageCreate, error) {
	return ValidateBase64MessageJSON(base64Str)
}
