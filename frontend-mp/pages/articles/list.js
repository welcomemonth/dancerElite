const api = require('../../utils/api.js');
const app = getApp();

Page({
  data: {
    columnId: null,
    columnName: '',
    articles: [],
    page: 1,
    pageSize: 10,
    total: 0,
    loading: false,
    hasMore: true,
    noColumn: false
  },

  onLoad(options) {
    // 兼容 navigateTo 传参（保留旧逻辑）
    const { columnId, columnName } = options || {};
    if (columnId) {
      this.applyColumn(columnId, decodeURIComponent(columnName || ''));
    }
  },

  // tabBar 页面每次显示都会触发 onShow，用于消费 home 写入的 pendingColumn
  onShow() {
    const pending = app.globalData && app.globalData.pendingColumn;
    if (pending) {
      app.globalData.pendingColumn = null;
      this.applyColumn(pending.id, pending.name);
      return;
    }
    if (!this.data.columnId) {
      this.setData({ noColumn: true, loading: false });
      wx.setNavigationBarTitle({ title: '文章' });
    }
  },

  // 应用选中的栏目并加载第一页
  applyColumn(columnId, columnName) {
    this.setData({
      columnId,
      columnName,
      articles: [],
      page: 1,
      hasMore: true,
      noColumn: false
    });
    wx.setNavigationBarTitle({ title: columnName || '文章列表' });
    this.loadArticles();
  },

  onGoSelectColumn() {
    wx.switchTab({ url: '/pages/home/home' });
  },

  // 加载文章列表
  loadArticles() {
    if (this.data.loading || !this.data.hasMore) return;
    
    this.setData({ loading: true });
    
    api.getArticlesByColumn(this.data.columnId, this.data.page, this.data.pageSize)
      .then((res) => {
        const newArticles = res.list || [];
        this.setData({
          articles: this.data.page === 1 ? newArticles : this.data.articles.concat(newArticles),
          total: res.total,
          loading: false,
          hasMore: this.data.articles.length + newArticles.length < res.total
        });
      })
      .catch((err) => {
        this.setData({ loading: false });
        wx.showToast({ 
          title: '加载失败', 
          icon: 'none' 
        });
      });
  },

  // 点击文章，跳转到详情页
  onArticleTap(e) {
    const articleId = e.currentTarget.dataset.id;
    wx.navigateTo({
      url: `/pages/articles/detail?id=${articleId}`
    });
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.setData({ page: 1, hasMore: true });
    this.loadArticles();
    wx.stopPullDownRefresh();
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.hasMore) {
      this.setData({ page: this.data.page + 1 });
      this.loadArticles();
    }
  }
}); 