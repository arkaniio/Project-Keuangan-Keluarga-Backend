package controller

import "project-keuangan-keluarga/service"

type ControllerHandlerCategory struct {
	CategoryService service.CategoryService
}

func NewControllerHandlerCategory(categoryService service.CategoryService) *ControllerHandlerCategory {
	return &ControllerHandlerCategory{CategoryService: categoryService}
}
