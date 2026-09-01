package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/dancer-elite/internal/pkg/response"
	"github.com/welcomemonth/dancer-elite/internal/service"
)

// RankingHandler 排行榜 / 选手处理器（小程序端只读）
type RankingHandler struct {
	svc *service.RankingService
}

// NewRankingHandler 创建排行榜处理器
func NewRankingHandler(svc *service.RankingService) *RankingHandler {
	return &RankingHandler{svc: svc}
}

// ActiveSeason 当前激活赛季
func (h *RankingHandler) ActiveSeason(c *gin.Context) {
	season, err := h.svc.ActiveSeason(c.Request.Context())
	if err != nil {
		response.NotFound(c, "暂无激活赛季")
		return
	}
	response.OK(c, season)
}

// Leaderboard 年度积分排行榜
func (h *RankingHandler) Leaderboard(c *gin.Context) {
	seasonID, _ := strconv.ParseInt(c.Query("season_id"), 10, 64)
	ageGroup := c.Query("age_group")
	danceType := c.Query("dance_type")

	list, err := h.svc.Leaderboard(c.Request.Context(), seasonID, ageGroup, danceType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, list)
}

// OrgLeaderboard 机构排行榜
func (h *RankingHandler) OrgLeaderboard(c *gin.Context) {
	seasonID, _ := strconv.ParseInt(c.Query("season_id"), 10, 64)

	list, err := h.svc.OrgLeaderboard(c.Request.Context(), seasonID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, list)
}

// PlayerDetail 选手详情
func (h *RankingHandler) PlayerDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	detail, err := h.svc.PlayerDetail(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "选手不存在")
		return
	}
	response.OK(c, detail)
}
