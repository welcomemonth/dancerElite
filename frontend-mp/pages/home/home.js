const api = require('../../utils/api.js');
const app = getApp();

Page({
  data: {
    columns: [],
    loading: true
  },

  onLoad() {
    this.loadColumns();
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.setData({ loading: true });
    this.loadColumns(() => wx.stopPullDownRefresh());
  },

  // 加载栏目列表
  loadColumns(done) {
    api.getColumns()
      .then((columns) => {
        this.setData({
          columns: columns,
          loading: false
        });
        done && done();
      })
      .catch((err) => {
        this.setData({ loading: false });
        wx.showToast({
          title: '加载栏目失败',
          icon: 'none'
        });
        done && done();
      });
  },

  // 点击栏目，跳转到文章列表（articles/list 是 tabBar 页面，参数通过全局变量传递）
  onColumnTap(e) {
    const columnId = e.currentTarget.dataset.id;
    const columnName = e.currentTarget.dataset.name;
    app.globalData.pendingColumn = { id: columnId, name: columnName };
    wx.switchTab({ url: '/pages/articles/list' });
  }
}); 