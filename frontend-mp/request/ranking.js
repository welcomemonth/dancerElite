// 排行榜模块：排行 / 选手（公开只读）
const { get } = require('./http');

/**
 * 年度积分排行榜（按年龄组、舞种筛选）
 */
function leaderboard({ seasonId, ageGroup, danceType } = {}) {
  return get('/rankings', {
    season_id: seasonId,
    age_group: ageGroup,
    dance_type: danceType
  });
}

/**
 * 机构排行榜
 */
function orgLeaderboard({ seasonId } = {}) {
  return get('/rankings/organization', { season_id: seasonId });
}

/**
 * 选手详情 + 成绩
 */
function playerDetail(id) {
  return get(`/players/${id}`);
}

module.exports = { leaderboard, orgLeaderboard, playerDetail };
