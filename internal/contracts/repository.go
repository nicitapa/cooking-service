package contracts

import (
	"context"
	"github.com/nicitapa/cooking-service/internal/models"
)

type RepositoryI interface {
	GetAll(ctx context.Context) ([]models.Recipe, error)
	GetByID(ctx context.Context, id int64) (*models.Recipe, error)
	Create(ctx context.Context, recipe *models.Recipe)
	Update(ctx context.Context, recipe *models.Recipe)
	Delete(ctx context.Context, id int64)
	Search(ctx context.Context, ingredients, tags []string)
}
