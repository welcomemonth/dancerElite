const { COMPETITIONS, STATUS } = require('../../../utils/mock.js');

// 奖项配色（金/银/铜）
const AWARD_CFG = {
  '金': { bg: 'rgba(232,150,10,0.1)', color: '#B5730A', border: 'rgba(232,150,10,0.3)' },
  '银': { bg: 'rgba(100,110,130,0.08)', color: '#5A6380', border: 'rgba(100,110,130,0.2)' },
  '铜': { bg: 'rgba(160,104,58,0.08)', color: '#8C5A2A', border: 'rgba(160,104,58,0.2)' }
};

// 赛事亮点（预宣传用，固定文案）
const HIGHLIGHTS = [
  { e: '🏆', t: '国家级专业评审团现场打分' },
  { e: '🎖️', t: '颁发金、银、铜等级荣誉证书' },
  { e: '📊', t: '积分计入2025年度积分排行榜' },
  { e: '🎭', t: '多舞种、多年龄组均可参赛' }
];

Page({
  data: {
    event: null,
    // 报名表单（报名中）
    formName: '',
    formAge: '',
    formStyle: '',
    submitted: false,
    // 预宣传
    reminded: false,
    // 已结束
    rankStyles: [],
    rankAges: [],
    rankStyle: '',
    rankAge: '',
    rankings: []
  },

  onLoad(options) {
    const id = options.id;
    const raw = COMPETITIONS.find((c) => c.id === id);

    if (!raw) {
      wx.showToast({ title: '赛事不存在', icon: 'none' });
      return;
    }

    const statusMeta = STATUS[raw.status];
    const infoRows = [
      { icon: '📅', text: raw.date },
      { icon: '📍', text: raw.location },
      { icon: '👥', text: raw.organizer },
      { icon: '⏰', text: raw.deadline }
    ];

    // 报名进度（报名中）
    const progress = raw.maxParticipants ? Math.round((raw.registered / raw.maxParticipants) * 100) : 62;

    // 获奖名单筛选项（已结束）：只保留 古典舞 / 民族民间舞
    const rankStyles = (raw.styles || []).filter((s) => s === '古典舞' || s === '民族民间舞');
    const rankAges = raw.ageGroups || [];
    const rankStyle = rankStyles[0] || '';
    // 默认选中第一个“有成绩”的年龄组，避免进入即空白
    let rankAge = rankAges[0] || '';
    if (rankStyle) {
      const firstWithData = rankAges.find((a) =>
        (raw.rankings || []).some((r) => r.style === rankStyle && r.age === a)
      );
      if (firstWithData) rankAge = firstWithData;
    }

    const event = Object.assign({}, raw, {
      statusMeta,
      infoRows,
      progress,
      countdown: 47, // mock：距开放报名天数
      highlights: HIGHLIGHTS,
      schoolCount: Math.round(raw.registered / 7)
    });

    this._event = raw;
    this.setData({
      event,
      formAge: rankAges[0] || '',
      formStyle: raw.styles[0] || '',
      rankStyles,
      rankAges,
      rankStyle,
      rankAge
    });
    this.applyRankings();
  },

  // ===== 报名表单 =====
  onNameInput(e) {
    this.setData({ formName: e.detail.value });
  },
  onAgeSelect(e) {
    this.setData({ formAge: e.currentTarget.dataset.a });
  },
  onStyleSelect(e) {
    this.setData({ formStyle: e.currentTarget.dataset.s });
  },
  onSubmit() {
    if (!this.data.formName.trim()) return;
    this.setData({ submitted: true });
  },

  // ===== 预宣传提醒 =====
  onRemind() {
    this.setData({ reminded: true });
  },

  // ===== 已结束：获奖名单筛选 =====
  onRankStyle(e) {
    this.setData({ rankStyle: e.currentTarget.dataset.s }, () => this.applyRankings());
  },
  onRankAge(e) {
    this.setData({ rankAge: e.currentTarget.dataset.a }, () => this.applyRankings());
  },
  applyRankings() {
    const raw = this._event;
    if (!raw || raw.status !== 'closed') return;
    const { rankStyle, rankAge } = this.data;
    const rankings = (raw.rankings || [])
      .filter((r) => r.style === rankStyle && r.age === rankAge)
      .map((r) => Object.assign({}, r, { awardMeta: AWARD_CFG[r.award] }));
    this.setData({ rankings });
  }
});
