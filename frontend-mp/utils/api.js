// API 调用封装
const API_BASE_URL = "http://localhost:8080/api/mp";

/**
 * 统一请求封装：解包后端 {code, message, data} 信封
 * 成功（HTTP 2xx 且 code === 0）resolve 出 data，否则 reject message
 */
function request(path, method = 'GET', data, needAuth = false) {
  return new Promise((resolve, reject) => {
    const header = { 'content-type': 'application/json' };
    if (needAuth) {
      header['Authorization'] = 'Bearer ' + (wx.getStorageSync('token') || '');
    }
    wx.request({
      url: API_BASE_URL + path,
      method,
      header,
      data,
      success: (res) => {
        const body = res.data || {};
        if (res.statusCode >= 200 && res.statusCode < 300 && body.code === 0) {
          resolve(body.data);
        } else {
          reject(new Error(body.message || '请求失败'));
        }
      },
      fail: (err) => reject(err)
    });
  });
}

// token 落盘 / 读取
function setToken(token) {
  wx.setStorageSync('token', token);
}
function getToken() {
  return wx.getStorageSync('token') || '';
}

/**
 * 小程序登录，后端返回 {token, user}，token 落盘后返回 user
 */
function login(code) {
  return request('/login', 'POST', { code }).then((data) => {
    if (data && data.token) setToken(data.token);
    return data && data.user;
  });
}

/**
 * 绑定手机号（注册），返回 {token, user}
 */
function register(phone, openid, name) {
  return request('/register', 'POST', { phone, openid, name }).then((data) => {
    if (data && data.token) setToken(data.token);
    return data && data.user;
  });
}

/**
 * 获取栏目列表
 */
function getColumns() {
  return request('/columns/');
}

/**
 * 根据栏目获取文章列表
 */
function getArticlesByColumn(columnId, page = 1, pageSize = 10) {
  return request(`/articles/column/${columnId}?page=${page}&page_size=${pageSize}`);
}

/**
 * 获取文章详情
 */
function getArticleDetail(id) {
  return request(`/articles/${id}`).then((article) => {
    if (article && article.created_at) {
      article.created_at = formatDate(article.created_at);
    }
    return article;
  });
}

/**
 * 获取当前赛季
 */
function getActiveSeason() {
  return request('/seasons/active');
}

/**
 * 获取排行榜（年度积分）
 * params: { season_id, age_group, dance_type }
 */
function getRankings(params = {}) {
  const qs = Object.keys(params)
    .filter((k) => params[k] !== '' && params[k] !== undefined && params[k] !== null)
    .map((k) => `${k}=${encodeURIComponent(params[k])}`)
    .join('&');
  return request('/rankings' + (qs ? '?' + qs : ''));
}

/**
 * 获取机构排行榜
 */
function getOrgRankings() {
  return request('/rankings/organization');
}

/**
 * 获取选手详情 + 成绩
 */
function getPlayerDetail(id) {
  return request(`/players/${id}`);
}

/**
 * 获取活动列表
 */
function getActivities(page = 1, pageSize = 10) {
  return request(`/activities/?page=${page}&page_size=${pageSize}`);
}

/**
 * 获取活动详情
 */
function getActivityDetail(id) {
  return request(`/activities/${id}`);
}

/**
 * 创建报名（需要登录态）
 * data: { activity_id, name, phone, id_card }
 */
function createRegistration(data) {
  return request('/registrations', 'POST', data, true);
}

/**
 * 我的报名（需要登录态）
 */
function getMyRegistrations(page = 1, pageSize = 10) {
  return request(`/registrations/mine?page=${page}&page_size=${pageSize}`, 'GET', null, true);
}

/**
 * 取消报名（需要登录态）
 */
function cancelRegistration(id) {
  return request(`/registrations/${id}/cancel`, 'PUT', {}, true);
}

/**
 * 创建支付订单（需要登录态），返回 {payment, pay_params}
 */
function createPayment(registrationId) {
  return request('/payments/create', 'POST', { registration_id: registrationId }, true);
}

/**
 * 格式化日期
 */
function formatDate(dateStr) {
  const date = new Date(dateStr);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

module.exports = {
  login,
  register,
  getColumns,
  getArticlesByColumn,
  getArticleDetail,
  getActiveSeason,
  getRankings,
  getOrgRankings,
  getPlayerDetail,
  getActivities,
  getActivityDetail,
  createRegistration,
  getMyRegistrations,
  cancelRegistration,
  createPayment,
  setToken,
  getToken
};
