package utils

import (
	"errors"
	"io"
)

// errImportNotImplemented 占位错误：成绩导入的 Excel 解析逻辑尚未实现。
// 实现 parseActivityResultExcel 后，可删除此变量并返回真实解析结果。
var errImportNotImplemented = errors.New("成绩导入解析逻辑待实现")

// activityResultRow 一条解析完成的成绩记录。
// 由 parseActivityResultExcel 从 Excel 中解析并补齐后返回，
// 上层（UploadResults）直接据此写入 activity_results 表。
//
// 注意：这里的字段已经是“可直接落库”的最终形态——
//   - PlayerID 已根据 Excel 中的选手标识（身份证号/姓名等）解析或创建得到；
//   - Points   已按积分规则（后台计算）填好。
type activityResultRow struct {
	PlayerID       int64  // 选手 ID（解析/匹配选手后得到）
	RegistrationID *int64 // 可选：关联的报名记录 ID
	DanceType      string // 舞种
	AgeGroup       string // 年龄组（U11/U13/U15 等）
	Rank           int    // 本场名次
	Points         int    // 本场积分（后台计算后写入）
	Award          string // 奖项（冠军/亚军/季军/特金奖 等）
	ParticipantNum int    // 该榜单本场参赛人数（便于后续追溯）
}

// parseActivityResultExcel 解析活动成绩 Excel 文件。
//
// TODO: 在此实现你自己的 Excel 解析方案（当前为占位，直接返回 errImportNotImplemented）。
//
// 建议的 Excel 列（示例，可自行调整，请与上传模板保持一致）：
//
//	序号 | 姓名 | 身份证 | 积分赛 | 组别 | 剧目 | 指导老师 | 所属机构 | 成绩 | 本站积分
//
// （具体列名、列顺序、是否含表头、支持 .xlsx/.xls 等，均由你决定）
//
// 实现步骤建议：
//  1. 在 go.mod 中引入 Excel 解析库（如 github.com/xuri/excelize/v2），
//     依据 filename 扩展名或文件内容识别 .xlsx / .xls。
//  2. 逐行读取，跳过表头，按列映射到 activityResultRow 字段。
//  3. 对每行解析/匹配选手：按身份证号（或姓名+机构）查找已有 Player，
//     不存在则创建，得到 PlayerID（可在此函数内直接调用 store 完成）。
//  4. 按积分规则计算 Points（如 冠军=100、亚军=90 …… 可后台配置）。
//  5. 返回解析完成的 rows；若无任何有效数据，返回空切片 + nil。
func parseActivityResultExcel(file io.Reader, filename string) ([]activityResultRow, error) {
	_ = file     // 占位：上传的文件内容
	_ = filename // 占位：原始文件名（用于识别扩展名 / 内容）

	// TODO: 实现 Excel 解析逻辑
	return nil, errImportNotImplemented
}
