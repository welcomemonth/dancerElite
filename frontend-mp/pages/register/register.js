const api = require('../../request/index.js');
const app = getApp();

Page({
  data: {
    phone: '',
    name: '',
    focused: '',
    submitting: false,
    registered: false,
    user: {}
  },

  onShow() {
    const user = (app.globalData && app.globalData.user) || {};
    if (user.phone) {
      this.setData({
        registered: true,
        user
      });
    } else {
      this.setData({ registered: false });
    }
  },

  onGoHome() {
    wx.switchTab({ url: '/pages/events/list' });
  },

  onFocus(e) {
    this.setData({ focused: e.currentTarget.dataset.field });
  },

  onBlur() {
    this.setData({ focused: '' });
  },

  onInputPhone(e) {
    this.setData({ phone: e.detail.value });
  },

  onInputName(e) {
    this.setData({ name: e.detail.value });
  },

  onSubmit() {
    const { phone, name } = this.data;
    if (!phone || phone.length !== 11) {
      wx.showToast({ title: '请输入正确的手机号', icon: 'none' });
      return;
    }
    if (!name) {
      wx.showToast({ title: '请输入昵称', icon: 'none' });
      return;
    }
    const openid = app.globalData.user && app.globalData.user.openid;
    if (!openid) {
      wx.showToast({ title: '登录态丢失，请重启', icon: 'none' });
      return;
    }

    this.setData({ submitting: true });
    api.auth.register({ phone, openid, name })
      .then((user) => {
        app.globalData.user = { ...app.globalData.user, ...user };
        wx.showToast({ title: '注册成功', icon: 'success' });
        setTimeout(() => {
          // 赛事列表是 tabBar 首页，使用 switchTab
          wx.switchTab({ url: '/pages/events/list' });
        }, 600);
      })
      .catch((err) => {
        wx.showToast({ title: (err && err.toString()) || '注册失败', icon: 'none' });
      })
      .finally(() => {
        this.setData({ submitting: false });
      });
  }
});
