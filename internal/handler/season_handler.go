package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/pkg/response"
	"github.com/welcomemonth/dancer-elite/internal/service"
)

// SeasonHandler 赛季处理器
type SeasonHandler struct {
	svc *service.SeasonService
}

// NewSeasonHandler 创建赛季处理器
func NewSeasonHandler(svc *service.SeasonService) *SeasonHandler {
	return &SeasonHandler{svc: svc}
}

type seasonRequest struct {
	Year      int        `json:"year" binding:"required"`
	Name      string     `json:"name" binding:"required"`
	Status    string     `json:"status"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

// List 获取赛季列表（后台）
func (h *SeasonHandler) List(c *gin.Context) {
	seasons, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, seasons)
}

// Get 获取单个赛季
func (h *SeasonHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	season, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, season)
}

// Create 创建赛季
func (h *SeasonHandler) Create(c *gin.Context) {
	var req seasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Status == "" {
		req.Status = "active"
	}

	season := &model.Season{
		Year:      req.Year,
		Name:      req.Name,
		Status:    req.Status,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if err := h.svc.Create(c.Request.Context(), season); err != nil {
		if err == errcode.ErrAlreadyExists {
			response.BadRequest(c, "该年份赛季已存在")
			return
		}
		response.ServerError(c, err.Error())
		return
	}
	response.Created(c, season)
}

// Update 更新赛季
func (h *SeasonHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req seasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]any{
		"year":       req.Year,
		"name":       req.Name,
		"status":     req.Status,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
	}

	if err := h.svc.Update(c.Request.Context(), id, updates); err != nil {
		if err == errcode.ErrAlreadyExists {
			response.BadRequest(c, "该年份赛季已存在")
			return
		}
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

// UpdateStatus 更新赛季状态（激活 / 归档）
func (h *SeasonHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

// Delete 删除赛季
func (h *SeasonHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.NoContent(c)
}
