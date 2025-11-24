package controller

import (
	"encoding/json"
	"github.com/nicitapa/cooking-service/internal/models"
	"github.com/nicitapa/cooking-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	router *gin.Engine
	svc    *service.RecipeService
}

func NewController(router *gin.Engine, svc *service.RecipeService) *Controller {
	return &Controller{router: router, svc: svc}
}

// @Summary Ping
// @Description Health-check endpoint
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func (ctrl *Controller) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ping": "pong"})
}

// @Summary List recipes
// @Description Get all recipes
// @Tags recipes
// @Produce json
// @Success 200 {array} models.Recipe
// @Router /recipes [get]
func (ctrl *Controller) GetAll(c *gin.Context) {
	recs, err := ctrl.svc.GetAll(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("get all recipes failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, recs)
}

// @Summary Get recipe
// @Description Get recipe by ID
// @Tags recipes
// @Produce json
// @Param id query int true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Failure 404 {object} map[string]string
// @Router /recipes/get [get]
func (ctrl *Controller) GetByID(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	rec, err := ctrl.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("get recipe by id failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if rec == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}
	c.JSON(http.StatusOK, rec)
}

// @Summary Create recipe
// @Description Create a new recipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body models.Recipe true "Recipe"
// @Success 201 {object} models.Recipe
// @Router /recipes/create [post]
func (ctrl *Controller) Create(c *gin.Context) {
	var m models.Recipe
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if err := ctrl.svc.Create(c.Request.Context(), &m); err != nil {
		log.Error().Err(err).Msg("create recipe failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, m)
}

// @Summary Update recipe
// @Description Update an existing recipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body models.Recipe true "Recipe"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /recipes/update [put]
func (ctrl *Controller) Update(c *gin.Context) {
	var m models.Recipe
	if err := json.NewDecoder(c.Request.Body).Decode(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if m.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	if err := ctrl.svc.Update(c.Request.Context(), &m); err != nil {
		log.Error().Err(err).Msg("update recipe failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// @Summary Delete recipe
// @Description Delete by ID
// @Tags recipes
// @Produce json
// @Param id query int true "Recipe ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /recipes/delete [delete]
func (ctrl *Controller) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := ctrl.svc.Delete(c.Request.Context(), id); err != nil {
		log.Error().Err(err).Msg("delete recipe failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// @Summary Search recipes
// @Description Search by ingredients and tags
// @Tags recipes
// @Accept json
// @Produce json
// @Param search body models.SearchRequest true "Search"
// @Success 200 {array} models.Recipe
// @Router /recipes/search [post]
func (ctrl *Controller) SearchByIngredientsAndTags(c *gin.Context) {
	var payload struct {
		Ingredients []string `json:"ingredients"`
		Tags        []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	recs, err := ctrl.svc.Search(c.Request.Context(), payload.Ingredients, payload.Tags)
	if err != nil {
		log.Error().Err(err).Msg("search recipes failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, recs)
}
