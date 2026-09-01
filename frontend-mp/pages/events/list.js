const { COMPETITIONS, STATUS } = require('../../utils/mock.js');

Page({
  data: {
    keyword: '',
    filter: 'all',
    filterTabs: [
      { key: 'all', label: '全部' },
      { key: 'pre', label: '预宣传' },
      { key: 'open', label: '报名中' },
      { key: 'closed', label: '已结束' }
    ],
    list: [],
    total: 0
  },

  onLoad() {
    this.applyFilter();
  },

  onInput(e) {
    this.setData({ keyword: e.detail.value }, () => this.applyFilter());
  },

  onFilter(e) {
    this.setData({ filter: e.currentTarget.dataset.key }, () => this.applyFilter());
  },

  applyFilter() {
    const keyword = this.data.keyword.trim();
    const filter = this.data.filter;

    const list = COMPETITIONS
      .filter((c) => filter === 'all' || c.status === filter)
      .filter((c) => !keyword || c.name.indexOf(keyword) !== -1 || c.location.indexOf(keyword) !== -1)
      .map((c) => {
        const statusMeta = STATUS[c.status];
        const footerText =
          c.status === 'open' ? `${c.registered}人已报名` :
          c.status === 'closed' ? `共${c.registered}人参赛` :
          c.deadline;
        return Object.assign({}, c, { statusMeta, footerText });
      });

    this.setData({ list, total: list.length });
  },

  onEventTap(e) {
    const id = e.currentTarget.dataset.id;
    wx.navigateTo({ url: '/pages/events/detail?id=' + id });
  }
});
