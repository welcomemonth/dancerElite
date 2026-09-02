package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/dancer-elite/internal/pkg/response"
	"github.com/welcomemonth/dancer-elite/internal/service"
)

// AnnualRankingHandler 年度积分榜处理器
type AnnualRankingHandler struct {
	svc *service.AnnualRankingService
}

// NewAnnualRankingHandler 创建年度积分榜处理器
func NewAnnualRankingHandler(svc *service.AnnualRankingService) *AnnualRankingHandler {
	return &AnnualRankingHandler{svc: svc}
}

// List 获取年度积分榜列表
func (h *AnnualRankingHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	seasonID, _ := strconv.ParseInt(c.DefaultQuery("season_id", "0"), 10, 64)
	ageGroup := c.Query("age_group")
	danceType := c.Query("dance_type")

	list, total, err := h.svc.List(c.Request.Context(), page, pageSize, seasonID, ageGroup, danceType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.PageOK(c, list, total, page, pageSize)
}

// Get 获取单条榜单记录
func (h *AnnualRankingHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	item, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "榜单记录不存在")
		return
	}
	response.OK(c, item)
}

// Update 更新榜单记录
func (h *AnnualRankingHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		TotalPoints *int `json:"total_points"`
		Rank        *int `json:"rank"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]any{}
	if req.TotalPoints != nil {
		updates["total_points"] = *req.TotalPoints
	}
	if req.Rank != nil {
		updates["rank"] = *req.Rank
	}
	if len(updates) == 0 {
		response.BadRequest(c, "没有需要更新的字段")
		return
	}

	if err := h.svc.Update(c.Request.Context(), id, updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

// Delete 删除榜单记录
func (h *AnnualRankingHandler) Delete(c *gin.Context) {
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

// Recalculate 重算某赛季的年度积分榜
func (h *AnnualRankingHandler) Recalculate(c *gin.Context) {
	var req struct {
		SeasonID int64 `json:"season_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "缺少 season_id")
		return
	}

	count, err := h.svc.RecalculateSeason(c.Request.Context(), req.SeasonID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"count": count})
}
