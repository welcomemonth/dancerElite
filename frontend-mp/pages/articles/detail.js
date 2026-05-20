const api = require('../../utils/api.js');

Page({
  data: {
    article: null,
    loading: true,
    liked: false,
    favorited: false
  },

  onLike() {
    this.setData({ liked: !this.data.liked });
    wx.showToast({
      title: this.data.liked ? '已点赞' : '已取消',
      icon: 'none'
    });
  },

  onFavorite() {
    this.setData({ favorited: !this.data.favorited });
    wx.showToast({
      title: this.data.favorited ? '已收藏' : '已取消收藏',
      icon: 'none'
    });
  },

  onShareAppMessage() {
    const a = this.data.article || {};
    return {
      title: a.title || '远山平台',
      path: `/pages/articles/detail?id=${a.id || ''}`
    };
  },

  onLoad(options) {
    const { id } = options;
    if (id) {
      this.loadArticle(id);
    } else {
      wx.showToast({ 
        title: '文章ID不存在', 
        icon: 'none' 
      });
      setTimeout(() => {
        wx.navigateBack();
      }, 1500);
    }
  },

  // 加载文章详情
  loadArticle(id) {
    this.setData({ loading: true });
    
    api.getArticleDetail(id)
      .then((article) => {
        this.setData({ 
          article: article,
          loading: false 
        });
        
        // 设置页面标题
        wx.setNavigationBarTitle({
          title: article.title || '文章详情'
        });
      })
      .catch((err) => {
        this.setData({ loading: false });
        wx.showToast({ 
          title: '加载失败', 
          icon: 'none' 
        });
        setTimeout(() => {
          wx.navigateBack();
        }, 1500);
      });
  },

  // 预览图片
  previewImage(e) {
    const current = e.currentTarget.dataset.src;
    const images = this.data.article.images || [current];
    
    wx.previewImage({
      current: current,
      urls: images
    });
  }
}); 