// 认证模块：登录 / 注册
const { post, setToken } = require('./http');

/**
 * 微信登录：code 换取 { token, user }，token 落盘后返回 user
 */
function login(code) {
  return post('/login', { code }).then((data) => {
    if (data && data.token) setToken(data.token);
    return data && data.user;
  });
}

/**
 * 绑定手机号注册：{ phone, openid, name }，返回 user
 */
function register({ phone, openid, name }) {
  return post('/register', { phone, openid, name }).then((data) => {
    if (data && data.token) setToken(data.token);
    return data && data.user;
  });
}

module.exports = { login, register };
