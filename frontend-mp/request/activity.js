// 赛事模块
const { get } = require('./http');

/**
 * 赛事列表（公开，分页）
 */
function list({ page = 1, pageSize = 10 } = {}) {
  return get('/activities/', { page, page_size: pageSize });
}

/**
 * 赛事详情（公开）
 */
function detail(id) {
  return get(`/activities/${id}`);
}

module.exports = { list, detail };
