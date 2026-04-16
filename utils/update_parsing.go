package utils

import (
	"errors"
	"fmt"
	"project-keuangan-keluarga/model"
	"strings"

	"github.com/google/uuid"
)

// fieldMapping defines a column name and its pointer value for dynamic update queries.
type fieldMapping struct {
	Column string
	Value  interface{} // the dereferenced value (nil if pointer is nil)
	IsSet  bool        // true if the pointer field is not nil
}

func buildUpdateQuery(table string, fields []fieldMapping, user_id uuid.UUID) (string, []interface{}, error) {
	var settings []string
	var args []interface{}
	argIdx := 1

	for _, f := range fields {
		if f.IsSet {
			settings = append(settings, fmt.Sprintf("%s=$%d", f.Column, argIdx))
			args = append(args, f.Value)
			argIdx++
		}
	}

	if len(settings) == 0 {
		return "", nil, errors.New("no fields to update")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE user_id = $%d", table, strings.Join(settings, ", "), argIdx)
	args = append(args, user_id)

	return query, args, nil
}

func buildUpdateQueryWithId(table string, fields []fieldMapping, id uuid.UUID) (string, []interface{}, error) {

	var settings []string
	var args []interface{}
	argsId := 1

	for _, f := range fields {
		if f.IsSet {
			settings = append(settings, fmt.Sprintf("%s=$%d", f.Column, argsId))
			args = append(args, f.Value)
			argsId++
		}
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table, strings.Join(settings, ", "), argsId)
	args = append(args, id)

	return query, args, nil

}

func UpdateToolsCategory(payload model.UpdatePayloadCategory, id uuid.UUID) (string, []interface{}, error) {
	// Validate type if provided
	if payload.Type != nil {
		if *payload.Type != "expense" && *payload.Type != "income" {
			return "", nil, errors.New("Failed to validate the type!")
		}
	}

	fields := []fieldMapping{
		{
			Column: "name", Value: valOrNil(payload.Name), IsSet: payload.Name != nil},
		{
			Column: "type", Value: valOrNil(payload.Type), IsSet: payload.Type != nil},
	}

	return buildUpdateQuery("categories", fields, id)
}

func UpdateToolsFamilie(payload model.UpdateFamilie, user_id uuid.UUID) (string, []interface{}, error) {

	field_map := []fieldMapping{
		{
			Column: "name", Value: valOrNil(payload.Name), IsSet: payload.Name != nil,
		},
		{
			Column: "created_by", Value: valOrNil(payload.Created_By), IsSet: payload.Created_By != nil,
		},
	}

	return buildUpdateQuery("families", field_map, user_id)

}

func UpdateToolsTransactions(payload model.UpdatePayloadTransaction, id uuid.UUID) (string, []interface{}, error) {
	// Validate type if provided
	if payload.Type != nil {
		if *payload.Type != "expense" && *payload.Type != "income" {
			return "", nil, errors.New("Failed to validate the type payload! invalid payload!")
		}
	}

	fields := []fieldMapping{
		{
			Column: "type", Value: valOrNil(payload.Type), IsSet: payload.Type != nil},
		{
			Column: "amount", Value: valOrNil(payload.Amount), IsSet: payload.Amount != nil},
		{
			Column: "category_id", Value: valOrNil(payload.CategoryId), IsSet: payload.CategoryId != nil},
		{
			Column: "description", Value: valOrNil(payload.Description), IsSet: payload.Description != nil},
		{
			Column: "date", Value: valOrNil(payload.Date), IsSet: payload.Date != nil},
	}

	return buildUpdateQuery("transactions", fields, id)
}

func UpdateToolsBudget(payload model.UpdatePayloadBudget, id uuid.UUID) (string, []interface{}, error) {

	field := []fieldMapping{
		{
			Column: "category_id", Value: valOrNil(payload.Category_Id), IsSet: payload.Category_Id != nil,
		},
		{
			Column: "limit_amount", Value: valOrNil(payload.Limit_amount), IsSet: payload.Limit_amount != nil,
		},
		{
			Column: "period", Value: valOrNil(payload.Period), IsSet: payload.Period != nil,
		},
		{
			Column: "start_date", Value: valOrNil(payload.StartDate), IsSet: payload.StartDate != nil,
		},
		{
			Column: "end_date", Value: valOrNil(payload.EndDate), IsSet: payload.EndDate != nil,
		},
		{
			Column: "is_active", Value: valOrNil(payload.IsActive), IsSet: payload.IsActive != nil,
		},
	}

	return buildUpdateQuery("budgets", field, id)

}

func UpdateToolsGoals(payload model.PayloadUpdateGoals, user_id uuid.UUID) (string, []interface{}, error) {

	if payload.Current_amount != nil && payload.Target_amount != nil {
		if *payload.Current_amount >= *payload.Target_amount {
			*payload.Status = "completed"
		} else {
			*payload.Status = "active"
		}
	}

	field := []fieldMapping{
		{
			Column: "name", Value: valOrNil(payload.Name), IsSet: payload.Name != nil,
		},
		{
			Column: "target_amount", Value: valOrNil(payload.Target_amount), IsSet: payload.Target_amount != nil,
		},
		{
			Column: "current_amount", Value: valOrNil(payload.Current_amount), IsSet: payload.Current_amount != nil,
		},
		{
			Column: "start_date", Value: valOrNil(payload.Start_date), IsSet: payload.Start_date != nil,
		},
		{
			Column: "target_date", Value: valOrNil(payload.Target_date), IsSet: payload.Target_date != nil,
		},
		{
			Column: "status", Value: valOrNil(payload.Status), IsSet: payload.Status != nil,
		},
	}

	return buildUpdateQuery("goals", field, user_id)

}

func UpdateToolsUser(payload model.UpdatePayloadUser, id uuid.UUID) (string, []interface{}, error) {

	if payload.Email != nil {
		if err := IsValidEmail(*payload.Email); err != nil {
			return "", nil, errors.New("Failed to get and validate the email!")
		}
	}

	if payload.Password != nil {
		hash_password, err := HashPassword(*payload.Password)
		if err != nil {
			return "", nil, errors.New("Failed to hash the password!")
		}
		payload.Password = &hash_password
	}

	field := []fieldMapping{
		{
			Column: "username", Value: valOrNil(payload.Username), IsSet: payload.Username != nil,
		},
		{
			Column: "email", Value: valOrNil(payload.Email), IsSet: payload.Email != nil,
		},
		{
			Column: "password", Value: valOrNil(payload.Password), IsSet: payload.Password != nil,
		},
		{
			Column: "profile_img", Value: valOrNil(payload.Profile_img), IsSet: payload.Profile_img != nil,
		},
	}

	return buildUpdateQuery("users", field, id)

}

func UpdateToolsFamilyMember(payload model.UpdateFamilyMember, user_id uuid.UUID) (string, []interface{}, error) {

	field_map := []fieldMapping{
		{
			Column: "family_id", Value: valOrNil(payload.FamilyId), IsSet: payload.FamilyId != nil,
		},
		{
			Column: "role", Value: valOrNil(payload.Role), IsSet: payload.Role != nil,
		},
	}

	return buildUpdateQuery("family_members", field_map, user_id)

}

// valOrNil safely dereferences any pointer, returning the value or nil.
func valOrNil[T any](p *T) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
