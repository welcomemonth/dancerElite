const app = getApp();

Page({
  data: {
    user: {},
    avatarText: 'A'
  },

  onShow() {
    const user = (app.globalData && app.globalData.user) || {};
    const name = user.name || '';
    this.setData({
      user,
      avatarText: name ? name.charAt(0).toUpperCase() : 'A'
    });
  },

  onEditProfile() {
    wx.showToast({ title: '功能开发中', icon: 'none' });
  },

  onMyRegistrations() {
    wx.showToast({ title: '功能开发中', icon: 'none' });
  },

  onMyOrders() {
    wx.showToast({ title: '功能开发中', icon: 'none' });
  },

  onAbout() {
    wx.showModal({
      title: '关于远山平台',
      content: '远山平台是一站式的内容、活动与用户运营工作台。',
      showCancel: false
    });
  },

  onContact() {
    wx.showToast({ title: '客服微信：yuanshan-cs', icon: 'none', duration: 2500 });
  }
});
