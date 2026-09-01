const api = require('./request/index.js');
App({
  onLaunch() {
    wx.login({
      success: (res) => {
        api.auth.login(res.code)
          .then((user) => {
            this.globalData.user = user;
            // 如果还没绑定手机号，跳转到注册页（register 已降级为非 tab 页，用 navigateTo）
            // if (user && !user.phone) {
            //   const pages = getCurrentPages();
            //   const currentPage = pages[pages.length - 1];
            //   if (currentPage && currentPage.route !== 'pages/register/register') {
            //     wx.navigateTo({ url: '/pages/register/register' });
            //   }
            // }
          })
          .catch((err) => {
            wx.showToast({ title: err.toString(), icon: 'none' });
          });
      },
      fail: () => {
        wx.showToast({ title: '微信登录失败', icon: 'none' });
      }
    });
  },
  globalData: {
    user: null
  }
});
