package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"mini-product-management/backend/internal/model"
	"mini-product-management/backend/internal/validation"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	db *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

func (h *ProductHandler) List(ctx *gin.Context) {
	search := strings.TrimSpace(ctx.Query("search"))
	if err := validation.Search(search); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := ctx.Query("category")
	if category != "" {
		if err := validation.Category(category); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 100"})
		return
	}

	query := `SELECT id, name, slug, COALESCE(description, ''), price, category, created_at, updated_at
		FROM products WHERE (? = '' OR name LIKE ?) AND (? = '' OR category = ?)
		ORDER BY id DESC LIMIT ?`
	rows, err := h.db.Query(query, search, "%"+search+"%", category, category, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Category, &p.CreatedAt, &p.UpdatedAt); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read products"})
			return
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read products"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": products, "count": len(products)})
}

func (h *ProductHandler) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if err := validation.Slug(slug); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var p model.Product
	err := h.db.QueryRow(`SELECT id, name, slug, COALESCE(description, ''), price, category, created_at, updated_at
		FROM products WHERE slug = ?`, slug).
		Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Category, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ProductHandler) Create(ctx *gin.Context) {
	var input model.ProductInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.TrimSpace(input.Slug)
	input.Category = strings.TrimSpace(input.Category)
	if err := validation.ProductInput(input.Name, input.Slug, input.Category, input.Price); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.db.Exec(`INSERT INTO products (name, slug, description, price, category) VALUES (?, ?, ?, ?, ?)`,
		input.Name, input.Slug, input.Description, input.Price, input.Category)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "slug already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "product created but could not read its ID"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "product created", "id": id})
}

func (h *ProductHandler) Update(ctx *gin.Context) {
	id, err := validation.PositiveInt("id", ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input model.ProductInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.TrimSpace(input.Slug)
	input.Category = strings.TrimSpace(input.Category)
	if err := validation.ProductInput(input.Name, input.Slug, input.Category, input.Price); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.db.Exec(`UPDATE products SET name=?, slug=?, description=?, price=?, category=? WHERE id=?`,
		input.Name, input.Slug, input.Description, input.Price, input.Category, id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "slug already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found or no changes made"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "product updated"})
}

func (h *ProductHandler) Delete(ctx *gin.Context) {
	id, err := validation.PositiveInt("id", ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.db.Exec(`DELETE FROM products WHERE id=?`, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *ProductHandler) Categories(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": []string{"electronics", "books", "clothing", "home"}})
}
