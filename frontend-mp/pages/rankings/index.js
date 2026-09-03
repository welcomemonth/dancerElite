const { LEVELS, DANCE_STYLES, ADVANCE_QUOTA } = require('../../utils/leaderboard.js');

const api = require('../../request/index.js');

Page({
  data: {
    levels: LEVELS,
    styles: DANCE_STYLES, //就先这样 暂时直接写死都有哪些舞种和级别，后续再扩展
    level: 'U11',
    style: '民族民间舞',
    keyword: '',
    rows: [],
    empty: false,
    total: 0,
    season: null
  },

  async onLoad() {
    await getApp().getSeason().then((season) => {
      this.setData({ season });
    });
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

  async apply() {
    const { level, style, season, keyword } = this.data;
    if (!season || !season.id) return;

    let list = [];
    try {
      list = await api.ranking.leaderboard({
        seasonId: season.id,
        ageGroup: level,
        danceType: style
      });
    } catch (err) {
      console.error('[年度积分榜] 调用失败', err);
      this.setData({ rows: [], empty: true, total: 0 });
      return;
    }

    // 晋级名单：直通选手无条件晋级（不占名额）+ 按名次取满 ADVANCE_QUOTA 个非直通选手
    const advancingIds = new Set();
    let quota = 0;
    list.forEach((d) => {
      if (d.is_direct) {
        advancingIds.add(d.player_id);
        return;
      }
      if (quota < ADVANCE_QUOTA) {
        advancingIds.add(d.player_id);
        quota += 1;
      }
    });

    const kw = keyword.trim();
    const displayed = kw
      ? list.filter((d) => (d.name || '').indexOf(kw) !== -1)
      : list;

    let crossedBoundary = false;
    const rows = displayed.map((d, idx) => {
      const isAdvancing = advancingIds.has(d.player_id);

      const prevDancer = idx > 0 ? displayed[idx - 1] : null;
      const prevAdvancing = prevDancer ? advancingIds.has(prevDancer.player_id) : false;
      const showAdvancingLabel = isAdvancing && (!prevDancer || !prevAdvancing);

      let showDivider = false;
      if (!isAdvancing && prevDancer && prevAdvancing && !crossedBoundary) {
        showDivider = true;
        crossedBoundary = true;
      }

      // 前三名奖牌圆标
      let medalBg = '';
      let medalColor = '';
      if (d.rank === 1) {
        medalBg = 'linear-gradient(135deg,#FFD700,#FFA500)';
        medalColor = '#5a3000';
      } else if (d.rank === 2) {
        medalBg = 'linear-gradient(135deg,#C0C0C0,#9E9E9E)';
        medalColor = '#ffffff';
      } else if (d.rank === 3) {
        medalBg = 'linear-gradient(135deg,#CD7F32,#8B4513)';
        medalColor = '#ffffff';
      }

      // 名次变化（previous_rank === 0 表示新上榜）
      let changeType;
      let changeText;
      if (d.previous_rank === 0) {
        changeType = 'new';
        changeText = 'NEW';
      } else if (d.rank_change === 0) {
        changeType = 'same';
        changeText = '—';
      } else if (d.rank_change > 0) {
        changeType = 'up';
        changeText = '▲' + d.rank_change;
      } else {
        changeType = 'down';
        changeText = '▼' + Math.abs(d.rank_change);
      }

      // 头像描边：直通金 > 晋级粉 > 普通灰
      let avatarBorder = isAdvancing ? '#ec4899' : '#e2e8f0';
      if (d.is_direct) avatarBorder = '#f59e0b';

      // 积分颜色
      let scoreColor = isAdvancing ? '#be185d' : '#64748b';
      if (d.is_direct) scoreColor = '#b45309';

      return {
        ...d, // 保留后端原始字段（rank/name/total_points/is_direct/avatar 等）
        avatarText: d.avatar ? '' : (d.name || '').slice(-1),
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
