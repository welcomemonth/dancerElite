// 赛季模块（公开只读）
const { get } = require('./http');

/**
 * 当前激活赛季：返回 { id, year, name, status, start_date, end_date }
 */
function active() {
  return get('/seasons/active');
}

module.exports = { active };
