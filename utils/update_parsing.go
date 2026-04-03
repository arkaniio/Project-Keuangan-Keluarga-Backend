package utils

import (
	"errors"
	"fmt"
	"project-keuangan-keluarga/model"
	"strings"

	"github.com/google/uuid"
)

func UpdateToolsCategory(payload model.UpdatePayloadCategory, id uuid.UUID) (string, error) {

	var settings []string
	var args []interface{}
	argsId := 1

	if payload.Name != nil {
		settings = append(settings, fmt.Sprintf("name=$%d", argsId))
		args = append(args, *payload.Name)
		argsId++
	}
	if payload.Type != nil {
		if *payload.Type != "expense" && *payload.Type != "income" {
			return "", errors.New("Failed to validate the type!")
		}
		settings = append(settings, fmt.Sprintf("type=$%d", argsId))
		args = append(args, *payload.Type)
		argsId++
	}

	query := fmt.Sprintf("UPDATE categories SET %s WHERE id = $%d", strings.Join(settings, ", "), argsId)
	args = append(args, id)

	return query, nil

}
