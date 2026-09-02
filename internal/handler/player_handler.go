package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/dancer-elite/internal/pkg/response"
	"github.com/welcomemonth/dancer-elite/internal/service"
)

// PlayerHandler 选手处理器
type PlayerHandler struct {
	svc *service.PlayerService
}

// NewPlayerHandler 创建选手处理器
func NewPlayerHandler(svc *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

type playerRequest struct {
	RealName    string `json:"real_name"`
	Gender      string `json:"gender"`
	Phone       string `json:"phone"`
	Institution string `json:"institution"`
	Teacher     string `json:"teacher"`
	AgeGroup    string `json:"age_group"`
}

// List 获取选手列表
func (h *PlayerHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	name := c.Query("name")
	institution := c.Query("institution")
	ageGroup := c.Query("age_group")

	list, total, err := h.svc.List(c.Request.Context(), page, pageSize, name, institution, ageGroup)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.PageOK(c, list, total, page, pageSize)
}

// Get 获取选手详情
func (h *PlayerHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	player, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "选手不存在")
		return
	}
	response.OK(c, player)
}

// Update 更新选手资料
func (h *PlayerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req playerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]any{
		"real_name":   req.RealName,
		"gender":      req.Gender,
		"phone":       req.Phone,
		"institution": req.Institution,
		"teacher":     req.Teacher,
		"age_group":   req.AgeGroup,
	}

	if err := h.svc.Update(c.Request.Context(), id, updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

// Delete 删除选手
func (h *PlayerHandler) Delete(c *gin.Context) {
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
