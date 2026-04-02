package service

import (
	"project-keuangan-keluarga/repository"
)

type CategoryService interface {
}

type repoCategory struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &repoCategory{repo: repo}
}
