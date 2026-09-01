// 报名模块（需要登录态）
const { get, post, put } = require('./http');

/**
 * 创建报名
 * data: { activity_id, name, phone, id_card, extra_info }
 */
function create(data) {
  return post('/registrations', data, { auth: true });
}

/**
 * 取消报名
 */
function cancel(id) {
  return put(`/registrations/${id}/cancel`, {}, { auth: true });
}

/**
 * 我的报名列表（分页）
 */
function my({ page = 1, pageSize = 10 } = {}) {
  return get('/registrations/mine', { page, page_size: pageSize }, { auth: true });
}

module.exports = { create, cancel, my };
