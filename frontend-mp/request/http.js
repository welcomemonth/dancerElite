// 统一请求核心：唯一与 wx.request 打交道的地方。
// 集中处理 baseURL、请求头、token 注入、响应信封解包与 query 序列化，
// 业务模块无需关心这些细节，只需调用 get/post/put/del。

const BASE_URL = 'http://localhost:8080/api/mp';

const TOKEN_KEY = 'token';

// ---------- token 管理 ----------
function getToken() {
  return wx.getStorageSync(TOKEN_KEY) || '';
}

function setToken(token) {
  wx.setStorageSync(TOKEN_KEY, token);
}

function clearToken() {
  wx.removeStorageSync(TOKEN_KEY);
}

// ---------- query 序列化 ----------
function buildQuery(params = {}) {
  const qs = Object.keys(params)
    .filter((k) => params[k] !== '' && params[k] !== undefined && params[k] !== null)
    .map((k) => `${encodeURIComponent(k)}=${encodeURIComponent(params[k])}`)
    .join('&');
  return qs ? '?' + qs : '';
}

// ---------- 统一请求 ----------
// 成功（HTTP 2xx 且 code === 0）resolve 出 data，否则 reject 一个带 code/statusCode 的 Error。
// 401 时自动清除本地 token，方便上层引导重新登录。
function request(path, { method = 'GET', data, params, auth = false } = {}) {
  return new Promise((resolve, reject) => {
    const header = { 'content-type': 'application/json' };
    if (auth) {
      header.Authorization = 'Bearer ' + getToken();
    }

    wx.request({
      url: BASE_URL + path + buildQuery(params),
      method,
      header,
      data,
      success: (res) => {
        const body = res.data || {};
        if (res.statusCode >= 200 && res.statusCode < 300 && body.code === 0) {
          resolve(body.data);
          return;
        }
        if (res.statusCode === 401) {
          clearToken();
        }
        const err = new Error(body.message || '请求失败');
        err.code = body.code;
        err.statusCode = res.statusCode;
        reject(err);
      },
      fail: (err) => reject(err)
    });
  });
}

// ---------- 便捷方法 ----------
const get = (path, params = {}, opts = {}) => request(path, { ...opts, method: 'GET', params });
const post = (path, data = {}, opts = {}) => request(path, { ...opts, method: 'POST', data });
const put = (path, data = {}, opts = {}) => request(path, { ...opts, method: 'PUT', data });
const del = (path, opts = {}) => request(path, { ...opts, method: 'DELETE' });

module.exports = {
  get,
  post,
  put,
  del,
  request,
  getToken,
  setToken,
  clearToken,
  BASE_URL
};
