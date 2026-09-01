const { LEVELS, DANCE_STYLES, ADVANCE_QUOTA, getDancers } = require('../../utils/leaderboard.js');

Page({
  data: {
    levels: LEVELS,
    styles: DANCE_STYLES,
    level: 'U13',
    style: '古典舞',
    keyword: '',
    rows: [],
    empty: false,
    total: 0
  },

  onLoad() {
    this.apply();
  },

  onLevel(e) {
    this.setData({ level: e.currentTarget.dataset.l }, () => this.apply());
  },

  onStyle(e) {
    this.setData({ style: e.currentTarget.dataset.s }, () => this.apply());
  },

  onSearch(e) {
    this.setData({ keyword: e.detail.value }, () => this.apply());
  },

  apply() {
    const sorted = getDancers(this.data.level, this.data.style);
    const keyword = this.data.keyword.trim();

    // 晋级名单：直通选手无条件晋级 + 按名次取满 ADVANCE_QUOTA 个非直通选手
    const advancingIds = new Set();
    let quota = 0;
    sorted.forEach((d) => {
      if (d.isDirect) {
        advancingIds.add(d.id);
        return;
      }
      if (quota < ADVANCE_QUOTA) {
        advancingIds.add(d.id);
        quota += 1;
      }
    });

    const displayed = keyword
      ? sorted.filter((d) => d.name.indexOf(keyword) !== -1)
      : sorted;

    let crossedBoundary = false;
    const rows = displayed.map((dancer, idx) => {
      const isAdvancing = advancingIds.has(dancer.id);
      const globalRank = sorted.findIndex((d) => d.id === dancer.id) + 1;

      const prevDancer = idx > 0 ? displayed[idx - 1] : null;
      const prevAdvancing = prevDancer ? advancingIds.has(prevDancer.id) : false;
      const showAdvancingLabel = isAdvancing && (!prevDancer || !prevAdvancing);

      let showDivider = false;
      if (!isAdvancing && prevDancer && prevAdvancing && !crossedBoundary) {
        showDivider = true;
        crossedBoundary = true;
      }

      // 前三名奖牌圆标
      let medalBg = '';
      let medalColor = '';
      if (globalRank === 1) {
        medalBg = 'linear-gradient(135deg,#FFD700,#FFA500)';
        medalColor = '#5a3000';
      } else if (globalRank === 2) {
        medalBg = 'linear-gradient(135deg,#C0C0C0,#9E9E9E)';
        medalColor = '#ffffff';
      } else if (globalRank === 3) {
        medalBg = 'linear-gradient(135deg,#CD7F32,#8B4513)';
        medalColor = '#ffffff';
      }

      // 名次变化
      const rc = dancer.rankChange;
      let changeType;
      let changeText;
      if (rc === 999) {
        changeType = 'new';
        changeText = 'NEW';
      } else if (rc === 0) {
        changeType = 'same';
        changeText = '—';
      } else if (rc > 0) {
        changeType = 'up';
        changeText = '▲' + rc;
      } else {
        changeType = 'down';
        changeText = '▼' + Math.abs(rc);
      }

      // 头像描边：直通金 > 晋级粉 > 普通灰
      let avatarBorder = isAdvancing ? '#ec4899' : '#e2e8f0';
      if (dancer.isDirect) avatarBorder = '#f59e0b';

      // 积分颜色
      let scoreColor = isAdvancing ? '#be185d' : '#64748b';
      if (dancer.isDirect) scoreColor = '#b45309';

      return {
        id: dancer.id,
        name: dancer.name,
        avatar: dancer.avatar,
        isDirect: dancer.isDirect,
        directReason: dancer.directReason || '',
        totalScore: dancer.totalScore,
        globalRank,
        isAdvancing,
        showAdvancingLabel,
        showDivider,
        medalBg,
        medalColor,
        changeType,
        changeText,
        avatarBorder,
        scoreColor
      };
    });

    this.setData({ rows, empty: displayed.length === 0, total: displayed.length });
  }
});
