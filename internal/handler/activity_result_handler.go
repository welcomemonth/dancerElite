package handler

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/pkg/response"
	"github.com/welcomemonth/dancer-elite/internal/service"
	"github.com/welcomemonth/dancer-elite/internal/utils"
)

// ActivityResultHandler 赛事成绩处理器
type ActivityResultHandler struct {
	svc *service.ActivityResultService
}

// NewActivityResultHandler 创建赛事成绩处理器
func NewActivityResultHandler(svc *service.ActivityResultService) *ActivityResultHandler {
	return &ActivityResultHandler{svc: svc}
}

type activityResultRequest struct {
	ActivityID     int64  `json:"activity_id" binding:"required"`
	PlayerID       int64  `json:"player_id" binding:"required"`
	SeasonID       int64  `json:"season_id"`
	RegistrationID *int64 `json:"registration_id"`
	DanceType      string `json:"dance_type" binding:"required"`
	AgeGroup       string `json:"age_group" binding:"required"`
	Rank           int    `json:"rank" binding:"required"`
	Points         int    `json:"points"`
	Award          string `json:"award"`
	ParticipantNum int    `json:"participant_num"`
}

// List 获取成绩列表
func (h *ActivityResultHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	activityID, _ := strconv.ParseInt(c.DefaultQuery("activity_id", "0"), 10, 64)
	playerID, _ := strconv.ParseInt(c.DefaultQuery("player_id", "0"), 10, 64)
	seasonID, _ := strconv.ParseInt(c.DefaultQuery("season_id", "0"), 10, 64)
	danceType := c.Query("dance_type")
	ageGroup := c.Query("age_group")

	list, total, err := h.svc.List(c.Request.Context(), page, pageSize, activityID, playerID, seasonID, danceType, ageGroup)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.PageOK(c, list, total, page, pageSize)
}

// Get 获取成绩详情
func (h *ActivityResultHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	item, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "成绩不存在")
		return
	}
	response.OK(c, item)
}

// Create 创建成绩
func (h *ActivityResultHandler) Create(c *gin.Context) {
	var req activityResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result := &model.ActivityResult{
		ActivityID:     req.ActivityID,
		PlayerID:       req.PlayerID,
		SeasonID:       req.SeasonID,
		RegistrationID: req.RegistrationID,
		DanceType:      req.DanceType,
		AgeGroup:       req.AgeGroup,
		Rank:           req.Rank,
		Points:         req.Points,
		Award:          req.Award,
		ParticipantNum: req.ParticipantNum,
	}

	if err := h.svc.Create(c.Request.Context(), result); err != nil {
		h.handleCreateError(c, err)
		return
	}
	response.Created(c, result)
}

// Update 更新成绩
func (h *ActivityResultHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req activityResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]any{
		"activity_id":     req.ActivityID,
		"player_id":       req.PlayerID,
		"season_id":       req.SeasonID,
		"registration_id": req.RegistrationID,
		"dance_type":      req.DanceType,
		"age_group":       req.AgeGroup,
		"rank":            req.Rank,
		"points":          req.Points,
		"award":           req.Award,
		"participant_num": req.ParticipantNum,
	}

	if err := h.svc.Update(c.Request.Context(), id, updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"id": id})
}

// Delete 删除成绩
func (h *ActivityResultHandler) Delete(c *gin.Context) {
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

// handleCreateError 创建/导入成绩的错误映射，给出可读提示
func (h *ActivityResultHandler) handleCreateError(c *gin.Context, err error) {
	switch err {
	case errcode.ErrNotFound:
		response.BadRequest(c, "活动或选手不存在")
	case errcode.ErrActivityResultExists:
		response.BadRequest(c, err.Error())
	case errcode.ErrNoActiveSeason:
		response.BadRequest(c, err.Error())
	case errcode.ErrActivityNotEnded:
		response.BadRequest(c, err.Error())
	default:
		response.ServerError(c, err.Error())
	}
}

// Import 从 Excel 批量导入成绩（某活动某级别×舞种成绩单）
func (h *ActivityResultHandler) Import(c *gin.Context) {
	activityID, _ := strconv.ParseInt(c.PostForm("activity_id"), 10, 64)
	ageGroup := strings.TrimSpace(c.PostForm("age_group"))
	danceType := strings.TrimSpace(c.PostForm("dance_type"))
	if activityID == 0 || ageGroup == "" || danceType == "" {
		response.BadRequest(c, "缺少 activity_id / age_group / dance_type")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的成绩表文件")
		return
	}
	defer file.Close()
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".xlsx") {
		response.BadRequest(c, "仅支持 .xlsx 文件")
		return
	}

	tmp, err := os.CreateTemp("", "results-*.xlsx")
	if err != nil {
		response.ServerError(c, "创建临时文件失败")
		return
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := io.Copy(tmp, file); err != nil {
		tmp.Close()
		response.ServerError(c, "读取文件失败")
		return
	}
	tmp.Close()

	records, err := utils.ParseExcelToStruct(tmpName)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// TODO: 在此实现成绩入库的业务逻辑
	// 1. 校验活动状态为「已结束」（ActivityRepo.GetByID 判断 status == 4）
	// 2. 按身份证脱敏 / 姓名匹配选手（utils.MaskIDCard + PlayerRepo.FindAll）
	// 3. 逐条创建 model.ActivityResult（activity_id / season_id / player_id / dance_type / age_group / rank / points / award），并去重
	// 4. 返回导入结果 {imported, skipped, errors}

	response.OK(c, gin.H{
		"total":      len(records),
		"age_group":  ageGroup,
		"dance_type": danceType,
		"records":    records,
	})
}
