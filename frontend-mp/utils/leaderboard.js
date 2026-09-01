// 排行榜 mock 数据（表格式榜单，参照新设计稿）
// 舞种固定为 民族民间舞 / 古典舞；名字池为占位 mock（民族民间舞复用“民族舞”池、古典舞复用“现代舞”池）
// 后续接入后端时，用 api.getRankings() 替换 getDancers 即可

const AVATARS = [
  "https://images.unsplash.com/photo-1531746020798-e6953c6e8e04?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1508214751196-bcfd4ca60f91?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1517841905240-472988babdf9?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1502823403499-6ccfcf4fb453?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1529626455594-4ff0802cfb7e?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1520975916090-3105956dac38?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1504257432389-52343af06ae3?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1573497019940-1c28c88b4f3e?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1552058544-f2b08422138a?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1521119989659-a83eee488004?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1463453091185-61582044d556?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1558898479-33c0057a5d12?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1541823709867-1b206113eafd?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1520573687861-62f4c6f91b2b?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1546961342-ea5f62d5a27b?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1548142813-c348350df52b?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1560087637-bf797bc7796a?w=64&h=64&fit=crop&auto=format",
  "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?w=64&h=64&fit=crop&auto=format"
];

// rankChange: 正=上升、负=下降、0=不变、999=新上榜
const CHANGES = [0, 3, -2, 999, 1, -1, 5, 0, -3, 2, 999, 0, 4, -2, 1, -4, 0, 3, 999, -1, 2, 0, -2, 1, 3, -1, 0, 999, -3, 2];

const LEVELS = ['U11', 'U13', 'U15'];
const DANCE_STYLES = ['民族民间舞', '古典舞'];
const ADVANCE_QUOTA = 20;

function buildPool(names, directs, reasons, base) {
  return names.map((name, i) => ({
    id: i + 1,
    name,
    avatar: AVATARS[i % AVATARS.length],
    totalScore: Math.max(200, Math.round(base - i * 2.3 + (i % 5) * 1.1)),
    rankChange: CHANGES[i % CHANGES.length],
    isDirect: directs.includes(i),
    directReason: directs.includes(i) ? reasons[directs.indexOf(i)] : ''
  }));
}

const NAMES = {
  'U11-民族民间舞': {
    names: ['苗依依', '白雪儿', '花晨晨', '柏紫嫣', '茜子涵', '菡萏云', '荷语桐', '蕊晓彤', '芷梦琪', '兰云熙', '蓉诗晴', '菊子涵', '莲雅琳', '梅语桐', '兰晓桐', '芙蓉云', '桂思远', '棠紫涵', '杏雨欣', '槿梦云', '芸子墨', '葵晴天', '艾晓雨', '贝思远', '可梓萌', '丁芯瑜', '鄂欣彤', '樊紫涵', '葛佳琪'],
    directs: [1], reasons: ['全国少儿赛冠军'], base: 270
  },
  'U11-古典舞': {
    names: ['程紫涵', '薛语嫣', '潘晓宇', '魏芊羽', '尤雨欣', '卢梦云', '聂子墨', '戚晴天', '滕雅琳', '殷芯怡', '翟梓萌', '蒋雨桐', '柳晓彤', '武欣悦', '岳梦琪', '齐云熙', '康诗晴', '贺子涵', '文雅琳', '伍语桐', '关晓桐', '池思远', '仓芯瑜', '曹欣彤', '岑紫涵', '柴佳琪', '常诗雨', '车晨曦', '陈悦宁'],
    directs: [], reasons: [], base: 275
  },
  'U13-民族民间舞': {
    names: ['丁芯瑜', '鄂欣彤', '樊紫涵', '葛佳琪', '郝诗雨', '嵇晨曦', '纪悦宁', '寇梦琪', '冷子晴', '蒙雨桐', '那思佳', '欧嘉欣', '裴子墨', '祁晓桐', '戎梓涵', '沈云熙', '唐雨萱', '伍雪凝', '辛子妍', '严晓雨', '尹思远', '臧梓萌', '展芯瑜', '甄欣彤', '郑紫涵', '钟佳琪', '朱晨曦', '卓悦宁', '邹梦琪'],
    directs: [], reasons: [], base: 279
  },
  'U13-古典舞': {
    names: ['柏语嫣', '茜子涵', '菡萏云', '荷语桐', '蕊晓彤', '芷梦琪', '兰云熙', '蓉诗晴', '菊子涵', '莲雅琳', '梅语桐', '兰晓桐', '芙蓉云', '桂思远', '棠紫涵', '杏雨欣', '槿梦云', '芸子墨', '葵晴天', '艾晓雨', '贝思远', '可梓萌', '丁芯瑜', '鄂欣彤', '樊紫涵', '葛佳琪', '郝诗雨', '嵇晨曦', '纪悦宁'],
    directs: [0], reasons: ['全国青少年赛冠军'], base: 285
  },
  'U15-民族民间舞': {
    names: ['姜芯瑜', '焦欣彤', '解紫涵', '金佳琪', '靳诗雨', '经晨曦', '井悦宁', '景梦琪', '居子晴', '鞠雨桐', '阚思佳', '康嘉欣', '柯子墨', '孔晓桐', '匡梓涵', '旷云熙', '邝雨萱', '蒯雪凝', '况子妍', '赖晓雨', '蓝思远', '郎梓萌', '劳芯瑜', '乐欣彤', '雷紫涵', '冷佳琪', '黎诗雨', '李晨曦', '连悦宁'],
    directs: [0], reasons: ['全国民族舞冠军'], base: 283
  },
  'U15-古典舞': {
    names: ['管芯瑜', '郭欣彤', '韩紫涵', '杭佳琪', '郝诗雨', '何晨曦', '贺悦宁', '洪梦琪', '侯子晴', '胡雨桐', '花思佳', '滑嘉欣', '怀子墨', '黄晓桐', '惠梓涵', '霍云熙', '姬雨萱', '季雪凝', '纪子妍', '贾晓雨', '简思远', '江梓萌', '姜芯瑜', '焦欣彤', '解紫涵', '金佳琪', '靳诗雨', '经晨曦', '井悦宁'],
    directs: [], reasons: [], base: 287
  }
};

const DATA = {};
Object.keys(NAMES).forEach((key) => {
  const { names, directs, reasons, base } = NAMES[key];
  DATA[key] = buildPool(names, directs, reasons, base);
});

// 按积分降序返回某个 级别-舞种 的选手列表
function getDancers(level, danceStyle) {
  const list = DATA[`${level}-${danceStyle}`] || [];
  return list.slice().sort((a, b) => b.totalScore - a.totalScore);
}

module.exports = {
  LEVELS,
  DANCE_STYLES,
  ADVANCE_QUOTA,
  getDancers
};
