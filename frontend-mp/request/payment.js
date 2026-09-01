// 支付模块（需要登录态）
const { get, post } = require('./http');

/**
 * 创建支付订单，返回 { payment, pay_params }
 */
function create(registrationId) {
  return post('/payments/create', { registration_id: registrationId }, { auth: true });
}

/**
 * 查询支付状态（按订单号轮询）
 */
function query(orderNo) {
  return get('/payments/query', { order_no: orderNo }, { auth: true });
}

module.exports = { create, query };
