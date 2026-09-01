// 统一入口：聚合所有业务模块，对外只暴露一个命名空间对象。
// 用法：const api = require('../../request/index.js');  api.auth.login() / api.activity.list() ...
const auth = require('./auth');
const activity = require('./activity');
const ranking = require('./ranking');
const registration = require('./registration');
const payment = require('./payment');

module.exports = {
  auth,
  activity,
  ranking,
  registration,
  payment
};
