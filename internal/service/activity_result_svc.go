package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/store"
	"github.com/welcomemonth/dancer-elite/internal/utils"
)

// ActivityResultService 赛事成绩管理服务
type ActivityResultService struct {
	st *store.Store
}

// NewActivityResultService 创建赛事成绩服务
func NewActivityResultService(st *store.Store) *ActivityResultService {
	return &ActivityResultService{st: st}
}

// ActivityResultItem 成绩列表项（含活动标题、选手姓名）
type ActivityResultItem struct {
	model.ActivityResult
	ActivityTitle string `json:"activity_title"`
	PlayerName    string `json:"player_name"`
}

// List 分页查询成绩，支持按活动/选手/赛季/舞种/年龄组过滤
func (s *ActivityResultService) List(ctx context.Context, page, pageSize int, activityID, playerID, seasonID int64, danceType, ageGroup string) ([]ActivityResultItem, int64, error) {
	list, total, err := s.st.ActivityResultRepo.List(ctx, page, pageSize,
		func(db *gorm.DB) *gorm.DB {
			if activityID > 0 {
				db = db.Where("activity_id = ?", activityID)
			}
			if playerID > 0 {
				db = db.Where("player_id = ?", playerID)
			}
			if seasonID > 0 {
				db = db.Where("season_id = ?", seasonID)
			}
			if danceType != "" {
				db = db.Where("dance_type = ?", danceType)
			}
			if ageGroup != "" {
				db = db.Where("age_group = ?", ageGroup)
			}
			return db
		},
		func(db *gorm.DB) *gorm.DB { return db.Preload("Activity").Preload("Player").Order("created_at desc") },
	)
	if err != nil {
		return nil, 0, err
	}

	items := make([]ActivityResultItem, 0, len(list))
	for _, r := range list {
		items = append(items, toActivityResultItem(r))
	}
	return items, total, nil
}

// Get 查询单条成绩（含活动、选手关联）
func (s *ActivityResultService) Get(ctx context.Context, id int64) (*ActivityResultItem, error) {
	list, err := s.st.ActivityResultRepo.FindAll(ctx,
		func(db *gorm.DB) *gorm.DB { return db.Where("id = ?", id).Preload("Activity").Preload("Player") },
	)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errcode.ErrNotFound
	}
	item := toActivityResultItem(list[0])
	return &item, nil
}

// Create 创建成绩：校验关联存在、去重（同一活动同一选手同一舞种仅一条）、缺省赛季
func (s *ActivityResultService) Create(ctx context.Context, result *model.ActivityResult) error {
	if result.SeasonID == 0 {
		if err := s.resolveSeason(ctx, result); err != nil {
			return err
		}
	}
	if err := s.validate(ctx, result); err != nil {
		return err
	}

	exists, _ := s.st.ActivityResultRepo.Exists(ctx,
		"activity_id = ? AND player_id = ? AND dance_type = ?",
		result.ActivityID, result.PlayerID, result.DanceType)
	if exists {
		return errcode.ErrActivityResultExists
	}
	return s.st.ActivityResultRepo.Create(ctx, result)
}

// Update 更新成绩
func (s *ActivityResultService) Update(ctx context.Context, id int64, updates map[string]any) error {
	return s.st.ActivityResultRepo.Update(ctx, id, updates)
}

// Delete 删除成绩
func (s *ActivityResultService) Delete(ctx context.Context, id int64) error {
	return s.st.ActivityResultRepo.Delete(ctx, id)
}

// resolveSeason 未显式传入赛季时，优先取所属活动的赛季，其次取当前生效赛季
func (s *ActivityResultService) resolveSeason(ctx context.Context, result *model.ActivityResult) error {
	if result.ActivityID > 0 {
		if act, err := s.st.ActivityRepo.GetByID(ctx, result.ActivityID); err == nil && act.SeasonID > 0 {
			result.SeasonID = act.SeasonID
			return nil
		}
	}
	season, err := s.st.SeasonRepo.GetActive(ctx)
	if err != nil {
		return errcode.ErrNoActiveSeason
	}
	result.SeasonID = season.ID
	return nil
}

// validate 校验必填字段及活动、选手是否存在，且活动必须已结束
func (s *ActivityResultService) validate(ctx context.Context, result *model.ActivityResult) error {
	if result.ActivityID == 0 || result.PlayerID == 0 || result.DanceType == "" || result.AgeGroup == "" {
		return errcode.ErrInvalidParam
	}
	act, err := s.st.ActivityRepo.GetByID(ctx, result.ActivityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return err
	}
	if act.Status != 4 {
		return errcode.ErrActivityNotEnded
	}
	if _, err := s.st.PlayerRepo.GetByID(ctx, result.PlayerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return err
	}
	return nil
}

// ImportResult 成绩批量导入结果
type ImportResult struct {
	Total    int      `json:"total"`
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors"`
}

// ImportResults 将 Excel 解析出的成绩记录批量入库（某活动某级别×舞种一张成绩单）。
// 规则：
//  1. 活动必须已结束（status==4）；
//  2. 名次按「本站积分」降序自动生成（1..N）；
//  3. 按身份证号精确匹配选手，不存在则先创建选手（同时创建离线用户以保持 1:1 绑定）；
//  4. 同一活动同一选手同一舞种去重（已存在则跳过）。
func (s *ActivityResultService) ImportResults(ctx context.Context, activityID int64, ageGroup, danceType string, records []utils.PlayerScoreRecord) (*ImportResult, error) {
	act, err := s.st.ActivityRepo.GetByID(ctx, activityID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	if act.Status != 4 {
		return nil, errcode.ErrActivityNotEnded
	}

	res := &ImportResult{Total: len(records), Errors: []string{}}

	// 按本站积分降序排序，生成名次
	sorted := make([]utils.PlayerScoreRecord, len(records))
	copy(sorted, records)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].CurrentPoint > sorted[j].CurrentPoint
	})

	participantNum := len(sorted)

	for i, rec := range sorted {
		rank := i + 1

		player, err := s.resolveOrCreatePlayer(ctx, rec, ageGroup)
		if err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("[%s] 选手处理失败: %v", rec.Name, err))
			continue
		}

		exists, _ := s.st.ActivityResultRepo.Exists(ctx,
			"activity_id = ? AND player_id = ? AND dance_type = ?",
			activityID, player.ID, danceType)
		if exists {
			res.Skipped++
			continue
		}

		result := &model.ActivityResult{
			ActivityID:     activityID,
			SeasonID:       act.SeasonID,
			PlayerID:       player.ID,
			DanceType:      danceType,
			AgeGroup:       ageGroup,
			Rank:           rank,
			Points:         rec.CurrentPoint,
			Award:          rec.Award,
			ParticipantNum: participantNum,
		}
		if err := s.st.ActivityResultRepo.Create(ctx, result); err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("[%s] 写入失败: %v", rec.Name, err))
			continue
		}
		res.Imported++
	}

	return res, nil
}

// resolveOrCreatePlayer 按身份证号匹配选手，不存在则创建（含离线用户）。
func (s *ActivityResultService) resolveOrCreatePlayer(ctx context.Context, rec utils.PlayerScoreRecord, ageGroup string) (*model.Player, error) {
	idCard := strings.TrimSpace(rec.IDCard)
	if idCard == "" {
		return nil, errors.New("身份证号为空")
	}

	players, err := s.st.PlayerRepo.FindAll(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("id_card = ?", idCard)
	})
	if err != nil {
		return nil, err
	}
	if len(players) > 0 {
		return &players[0], nil
	}

	// 选手不存在 直接创建选手，暂时不绑定User
	gender, year, month, day := parseIDCard(idCard)

	player := &model.Player{
		RealName:    rec.Name,
		Gender:      gender,
		IDCard:      idCard,
		Institution: rec.Organization,
		Teacher:     rec.Tutor,
		BirthYear:   year,
		BirthMonth:  month,
		BirthDay:    day,
		AgeGroup:    ageGroup,
	}
	if err := s.st.PlayerRepo.Create(ctx, player); err != nil {
		return nil, fmt.Errorf("创建选手失败: %w", err)
	}
	return player, nil
}

// parseIDCard 从 18 位身份证号解析性别与出生年月日；解析失败时给出安全默认值。
func parseIDCard(idCard string) (gender string, year, month, day int) {
	if len(idCard) == 18 {
		if y, err := strconv.Atoi(idCard[6:10]); err == nil {
			year = y
		}
		if m, err := strconv.Atoi(idCard[10:12]); err == nil {
			month = m
		}
		if d, err := strconv.Atoi(idCard[12:14]); err == nil {
			day = d
		}
		if n, err := strconv.Atoi(string(idCard[16])); err == nil {
			if n%2 == 1 {
				gender = "male"
			} else {
				gender = "female"
			}
		}
	}
	if gender == "" {
		gender = "unknown"
	}
	return
}

// toActivityResultItem 映射为列表项，并清空嵌套关联避免 JSON 重复输出
func toActivityResultItem(r model.ActivityResult) ActivityResultItem {
	item := ActivityResultItem{ActivityResult: r}
	if r.Activity != nil {
		item.ActivityTitle = r.Activity.Title
	}
	if r.Player != nil {
		item.PlayerName = r.Player.RealName
	}
	item.Activity = nil
	item.Player = nil
	return item
}
