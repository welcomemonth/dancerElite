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

    // 进入小程序时拉取当前激活赛季（Promise 缓存，全小程序共享一次请求）
    this.getSeason().catch(() => {});
  },

  /**
   * 获取当前激活赛季，返回 Promise。
   * 首次调用发请求，之后直接复用缓存的 Promise，页面无需关心时序：
   *   getApp().getSeason().then(season => this.setData({ season }))
   * 失败时清空缓存，下次调用自动重试。
   */
  getSeason() {
    if (!this.seasonReady) {
      this.seasonReady = api.season.active()
        .then((season) => {
          this.globalData.season = season;
          return season;
        })
        .catch((err) => {
          this.seasonReady = null; // 允许下次调用重试
          console.error('[激活赛季] 获取失败', err);
          throw err;
        });
    }
    return this.seasonReady;
  },

  globalData: {
    user: null,
    season: null
  }
});
