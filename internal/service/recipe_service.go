package service

import (
	"context"
	"database/sql"

	"github.com/nicitapa/cooking-service/internal/models"
	"github.com/nicitapa/cooking-service/internal/repository"
	"github.com/rs/zerolog/log"
)

type RecipeService struct {
	repo *repository.RecipeRepository
}

func NewRecipeService(repo *repository.RecipeRepository) *RecipeService {
	return &RecipeService{repo: repo}
}

func (s *RecipeService) GetAll(ctx context.Context) ([]models.Recipe, error) {
	recipes, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all recipes")
		return nil, err
	}
	return recipes, nil
}

func (s *RecipeService) GetByID(ctx context.Context, id int64) (*models.Recipe, error) {
	rec, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Int64("id", id).Msg("failed to get recipe by id")
		return nil, err
	}
	return rec, nil
}

func (s *RecipeService) Create(ctx context.Context, recipe *models.Recipe) error {
	if err := s.repo.Create(ctx, recipe); err != nil {
		log.Error().Err(err).Msg("failed to create recipe")
		return err
	}
	return nil
}

func (s *RecipeService) Update(ctx context.Context, recipe *models.Recipe) error {
	if recipe.ID == 0 {
		return sql.ErrNoRows
	}
	if err := s.repo.Update(ctx, recipe); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		log.Error().Err(err).Int64("id", recipe.ID).Msg("failed to update recipe")
		return err
	}
	return nil
}

func (s *RecipeService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		log.Error().Err(err).Int64("id", id).Msg("failed to delete recipe")
		return err
	}
	return nil
}

func (s *RecipeService) Search(ctx context.Context, ingredients, tags []string) ([]models.Recipe, error) {
	recipes, err := s.repo.FindByIngredientsAndTags(ctx, ingredients, tags)
	if err != nil {
		log.Error().Err(err).Msg("failed to search recipes")
		return nil, err
	}
	return recipes, nil
}
