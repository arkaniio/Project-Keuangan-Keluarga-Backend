package controller

import (
	"context"
	"net/http"
	"project-keuangan-keluarga/middleware"
	"project-keuangan-keluarga/model"
	"project-keuangan-keluarga/service"
	"project-keuangan-keluarga/utils"
	"time"

	"github.com/google/uuid"
)

type ControllerHandlerFamilyMember struct {
	FamilyMemberService service.FamilyMemberService
}

func NewControllerHandlerFamilyMember(familyMemberService service.FamilyMemberService) ControllerHandlerFamilyMember {
	return ControllerHandlerFamilyMember{FamilyMemberService: familyMemberService}
}

func (c *ControllerHandlerFamilyMember) CreateFamilyMember_Bp(w http.ResponseWriter, r *http.Request) {

	var payload model.PayloadFamilyMember
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the json!", err.Error())
		return
	}

	if err := utils.ValidatePayloads(payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the payload!", err.Error())
		return
	}

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware token id!", err.Error())
		return
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the uuid type!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	member, err := utils.ParsingPayloadFamilyMember(payload, middleware_token_id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the payload!", err.Error())
		return
	}

	if err := c.FamilyMemberService.CreateFamilyMember(ctx, member); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to create the new family member!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully created the new family member!", nil)

}

func (c *ControllerHandlerFamilyMember) UpdateFamilyMember_Bp(w http.ResponseWriter, r *http.Request) {

	var payload model.UpdateFamilyMember
	if err := utils.DecodeJson(r, &payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the family member payload!", err.Error())
		return
	}

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware token id!", err.Error())
		return
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the uuid type!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.FamilyMemberService.UpdateFamilyMember(ctx, middleware_token_id, payload); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to update the family member data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to update family member!", true)

}

func (c *ControllerHandlerFamilyMember) DeleteFamilyMember_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware token id!", err.Error())
		return
	}
	if middleware_token_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the uuid type!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	if err := c.FamilyMemberService.DeleteFamilyMember(ctx, middleware_token_id); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to delete the family member data!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully to delete the family member data!", true)

}

func (c *ControllerHandlerFamilyMember) GetAllFamilyMember_Bp(w http.ResponseWriter, r *http.Request) {

	// Parse pagination from query params
	allowedSorts := []string{"role", "joined_at"}
	params := utils.ParsePaginationParams(r, allowedSorts, "joined_at")

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	userId, err := middleware.GetTokenId(w, r)
	if err != nil {
		return
	}

	paginatedData, err := c.FamilyMemberService.GetAllFamilyMember(ctx, userId, params)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get family members!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully retrieved family members!", paginatedData)

}

func (c *ControllerHandlerFamilyMember) GetMyMembership_Bp(w http.ResponseWriter, r *http.Request) {

	middleware_token_id, err := middleware.GetTokenId(w, r)
	if err != nil {
		utils.ResponseError(w, http.StatusUnauthorized, "Failed to get the user id from token!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()

	member, err := c.FamilyMemberService.GetFamilyMemberByUserId(ctx, middleware_token_id)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to check family membership", err.Error())
		return
	}

	if member == nil {
		utils.ResponseError(w, http.StatusNotFound, "No family membership found for this user", "Member record is empty")
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Successfully retrieved membership!", member)
}
