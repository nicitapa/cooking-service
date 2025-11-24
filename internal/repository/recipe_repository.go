package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/nicitapa/cooking-service/internal/models"
)

type RecipeRepository struct {
	DB *sqlx.DB
}

func NewRecipeRepository(db *sqlx.DB) *RecipeRepository {
	return &RecipeRepository{DB: db}
}

func (r *RecipeRepository) GetAll(ctx context.Context) ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.DB.SelectContext(ctx, &recipes, `
		SELECT id, title, description, instructions, image_url, category_id
		FROM recipes
		ORDER BY id DESC
	`)
	return recipes, err
}

func (r *RecipeRepository) GetByID(ctx context.Context, id int64) (*models.Recipe, error) {
	var rec models.Recipe
	err := r.DB.GetContext(ctx, &rec, `
		SELECT id, title, description, instructions, image_url, category_id
		FROM recipes
		WHERE id = $1
	`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *models.Recipe) error {
	return r.DB.QueryRowxContext(ctx, `
		INSERT INTO recipes (title, description, instructions, image_url, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, recipe.Title, recipe.Description, recipe.Instructions, recipe.ImageURL, recipe.CategoryID).Scan(&recipe.ID)
}

func (r *RecipeRepository) Update(ctx context.Context, recipe *models.Recipe) error {
	res, err := r.DB.ExecContext(ctx, `
		UPDATE recipes
		SET title = $1, description = $2, instructions = $3, image_url = $4, category_id = $5
		WHERE id = $6
	`, recipe.Title, recipe.Description, recipe.Instructions, recipe.ImageURL, recipe.CategoryID, recipe.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *RecipeRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.DB.ExecContext(ctx, `
		DELETE FROM recipes WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *RecipeRepository) FindByIngredientsAndTags(ctx context.Context, ingredients []string, tags []string) ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.DB.SelectContext(ctx, &recipes, `
		SELECT DISTINCT r.id, r.title, r.description, r.instructions, r.image_url, r.category_id
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		LEFT JOIN recipe_tags rt ON r.id = rt.recipe_id
		LEFT JOIN tags t ON rt.tag_id = t.id
		WHERE ($1::text[] IS NULL OR i.name = ANY($1))
		  AND ($2::text[] IS NULL OR t.name = ANY($2))
		ORDER BY r.id DESC
	`, nullableArray(ingredients), nullableArray(tags))
	return recipes, err
}

func nullableArray(arr []string) interface{} {
	if len(arr) == 0 {
		return nil
	}
	return pq.Array(arr)
}
