// 赛事 / 排行榜 mock 数据（参照「优舞盟」设计稿）
// 后续接入后端时，把这里的常量替换为 api.getActivities() / api.getRankings() 即可

// 赛事状态元信息
const STATUS = {
  pre: { label: '预宣传', color: '#7C3AED', bg: 'rgba(139,92,246,0.1)', dot: '#8B5CF6' },
  open: { label: '报名中', color: '#059669', bg: 'rgba(16,185,129,0.1)', dot: '#10B981' },
  closed: { label: '已结束', color: '#6B7280', bg: 'rgba(107,114,128,0.1)', dot: '#9CA3AF' }
};

// 赛事列表
const COMPETITIONS = [
  {
    id: '1',
    name: '第七届全国青少年舞蹈锦标赛',
    status: 'open',
    date: '2024年12月15-17日',
    location: '北京国家体育馆',
    organizer: '中国舞蹈协会',
    styles: ['民族民间舞', '古典舞'],
    ageGroups: ['U11', 'U13', 'U15'],
    deadline: '12月1日截止报名',
    coverGradient: 'linear-gradient(135deg, #1a0833 0%, #3d1a5c 60%, #200840 100%)',
    registered: 312,
    maxParticipants: 500,
    description: '全国青少年舞蹈锦标赛由中国舞蹈协会主办，是国家级青少年专业舞蹈赛事，旨在发现培养舞蹈新秀。本届赛事设民族民间舞与古典舞两大舞种，分U11、U13、U15三个年龄组别角逐，总奖金池30万元，参赛成绩全部计入年度积分榜。'
  },
  {
    id: '2',
    name: '2025优舞盟冬季邀请赛',
    status: 'pre',
    date: '2025年1月18-20日',
    location: '上海东方艺术中心',
    organizer: '优舞盟赛事委员会',
    styles: ['民族民间舞', '古典舞', '现代舞'],
    ageGroups: ['U11', 'U13', 'U15'],
    deadline: '2025年1月5日开放报名',
    coverGradient: 'linear-gradient(135deg, #030d1a 0%, #0a1e40 60%, #041228 100%)',
    registered: 0,
    description: '2025优舞盟冬季邀请赛将于新年伊始在沪盛大开幕，汇聚全国各地舞蹈少年，同台竞技，共谱华章。本赛事采用专业评审团打分制，颁发年度积分，是积分榜重要积分来源之一。'
  },
  {
    id: '3',
    name: '第三届粤港澳青少年舞蹈大赛',
    status: 'closed',
    date: '2024年10月5-7日',
    location: '广州大剧院',
    organizer: '粤港澳舞蹈协会',
    styles: ['民族民间舞', '古典舞'],
    ageGroups: ['U11', 'U13', 'U15'],
    deadline: '已截止',
    coverGradient: 'linear-gradient(135deg, #1a0808 0%, #3d1010 60%, #200a0a 100%)',
    registered: 428,
    description: '第三届粤港澳青少年舞蹈大赛已圆满落幕，共吸引来自粤港澳三地428名选手参赛，现场精彩纷呈，感谢所有参赛选手与工作人员的辛勤付出。',
    rankings: [
      { rank: 1, name: '林晓彤', school: '北京舞蹈学院附中', score: '98.6', award: '金', style: '古典舞', age: 'U13' },
      { rank: 2, name: '王雅琪', school: '上海舞蹈学校', score: '97.8', award: '金', style: '古典舞', age: 'U13' },
      { rank: 3, name: '张若璃', school: '广州市少年宫', score: '96.4', award: '金', style: '民族民间舞', age: 'U11' },
      { rank: 4, name: '陈思颖', school: '成都艺术学校', score: '95.2', award: '银', style: '民族民间舞', age: 'U15' },
      { rank: 5, name: '李梦瑶', school: '杭州舞蹈团', score: '94.8', award: '银', style: '古典舞', age: 'U15' },
      { rank: 6, name: '赵欣妍', school: '天津艺术学院', score: '94.1', award: '银', style: '古典舞', age: 'U13' },
      { rank: 7, name: '孙雨桐', school: '南京艺术学校', score: '93.5', award: '铜', style: '民族民间舞', age: 'U11' },
      { rank: 8, name: '黄思嘉', school: '深圳舞蹈学校', score: '92.8', award: '铜', style: '古典舞', age: 'U15' }
    ]
  },
  {
    id: '4',
    name: '全国少儿舞蹈展演积分赛',
    status: 'open',
    date: '2024年12月22-24日',
    location: '成都锦城艺术宫',
    organizer: '文化和旅游部艺术司',
    styles: ['民族民间舞'],
    ageGroups: ['U11', 'U13'],
    deadline: '12月10日截止报名',
    coverGradient: 'linear-gradient(135deg, #0a1a0e 0%, #0d3020 60%, #061510 100%)',
    registered: 186,
    maxParticipants: 300,
    description: '全国少儿舞蹈展演由文化和旅游部艺术司主办，兼具展演与积分双重属性，本届仅设民族民间舞项目，参赛成绩按比例计入年度积分榜。'
  },
  {
    id: '5',
    name: '华北区少儿舞蹈联赛春季赛',
    status: 'pre',
    date: '2025年3月8-9日',
    location: '天津大剧院',
    organizer: '华北舞蹈协会',
    styles: ['古典舞', '民族民间舞'],
    ageGroups: ['U11', 'U13', 'U15'],
    deadline: '2025年2月20日开放报名',
    coverGradient: 'linear-gradient(135deg, #1a150a 0%, #3d2f0a 60%, #201a06 100%)',
    registered: 0,
    description: '华北区少儿舞蹈联赛是华北地区最具影响力的青少年舞蹈赛事，每年吸引来自京津冀晋蒙的优秀选手参赛。本赛季春季赛将在天津大剧院盛大举办。'
  },
  {
    id: '6',
    name: '第五届江南杯舞蹈邀请赛',
    status: 'closed',
    date: '2024年9月14-15日',
    location: '杭州大剧院',
    organizer: '浙江省舞蹈协会',
    styles: ['民族民间舞', '古典舞'],
    ageGroups: ['U11', 'U13', 'U15'],
    deadline: '已截止',
    coverGradient: 'linear-gradient(135deg, #0a0f1a 0%, #151f3d 60%, #060c24 100%)',
    registered: 356,
    description: '第五届江南杯舞蹈邀请赛已圆满落幕，共吸引356名选手参赛，是华东地区年度精彩赛事之一。',
    rankings: [
      { rank: 1, name: '刘梓嫣', school: '浙江艺术学校', score: '97.9', award: '金', style: '古典舞', age: 'U13' },
      { rank: 2, name: '张若璃', school: '广州市少年宫', score: '97.2', award: '金', style: '民族民间舞', age: 'U11' },
      { rank: 3, name: '周沐阳', school: '上海舞蹈学校', score: '96.8', award: '金', style: '古典舞', age: 'U15' },
      { rank: 4, name: '林晓彤', school: '北京舞蹈学院附中', score: '96.1', award: '银', style: '古典舞', age: 'U13' },
      { rank: 5, name: '吴思羽', school: '南京艺术学院', score: '95.4', award: '银', style: '民族民间舞', age: 'U15' },
      { rank: 6, name: '陈思颖', school: '成都艺术学校', score: '94.6', award: '银', style: '民族民间舞', age: 'U15' }
    ]
  }
];

// 年度积分排行榜：舞种 -> 年龄组 -> 榜单
const LEADERBOARD = {
  '民族民间舞': {
    'U11': [
      { rank: 1, name: '张若璃', school: '广州市少年宫', points: 2680, gold: 3, silver: 1 },
      { rank: 2, name: '孙雨桐', school: '南京艺术学校', points: 2540, gold: 2, silver: 2 },
      { rank: 3, name: '周小曼', school: '武汉艺术学院', points: 2380, gold: 2, silver: 1 },
      { rank: 4, name: '何依依', school: '西安舞蹈学校', points: 2210, gold: 1, silver: 3 },
      { rank: 5, name: '罗湘云', school: '长沙艺术学校', points: 2050, gold: 1, silver: 2 },
      { rank: 6, name: '谢婷婷', school: '福州少儿艺术团', points: 1920, gold: 1, silver: 1 },
      { rank: 7, name: '钱小悦', school: '苏州艺术学校', points: 1780, gold: 0, silver: 3 },
      { rank: 8, name: '方子怡', school: '合肥市少年宫', points: 1650, gold: 0, silver: 2 }
    ],
    'U13': [
      { rank: 1, name: '陈思颖', school: '成都艺术学校', points: 2720, gold: 3, silver: 2 },
      { rank: 2, name: '蒋晓薇', school: '重庆舞蹈学院', points: 2560, gold: 2, silver: 3 },
      { rank: 3, name: '汪雨柔', school: '昆明艺术学校', points: 2390, gold: 2, silver: 1 },
      { rank: 4, name: '彭思琦', school: '贵阳少儿艺术团', points: 2180, gold: 1, silver: 3 },
      { rank: 5, name: '冯玉兰', school: '太原艺术学校', points: 2010, gold: 1, silver: 2 },
      { rank: 6, name: '高璐璐', school: '石家庄少年宫', points: 1840, gold: 0, silver: 4 },
      { rank: 7, name: '叶晨曦', school: '郑州舞蹈学院', points: 1690, gold: 0, silver: 2 },
      { rank: 8, name: '马欣欣', school: '西宁少儿舞蹈团', points: 1540, gold: 0, silver: 1 }
    ],
    'U15': [
      { rank: 1, name: '吴思羽', school: '南京艺术学院', points: 2960, gold: 4, silver: 1 },
      { rank: 2, name: '李梦瑶', school: '杭州舞蹈团', points: 2800, gold: 3, silver: 2 },
      { rank: 3, name: '赵诗语', school: '厦门艺术学院', points: 2640, gold: 3, silver: 1 },
      { rank: 4, name: '林雨桐', school: '大连艺术学校', points: 2460, gold: 2, silver: 2 },
      { rank: 5, name: '朱晓晴', school: '沈阳舞蹈学院', points: 2250, gold: 1, silver: 4 },
      { rank: 6, name: '徐梦蝶', school: '哈尔滨艺术学校', points: 2080, gold: 1, silver: 2 },
      { rank: 7, name: '龚晓雨', school: '长春艺术学院', points: 1910, gold: 0, silver: 3 },
      { rank: 8, name: '唐小语', school: '南昌少儿艺术团', points: 1760, gold: 0, silver: 2 }
    ]
  },
  '古典舞': {
    'U11': [
      { rank: 1, name: '周沐阳', school: '上海舞蹈学校', points: 2720, gold: 3, silver: 2 },
      { rank: 2, name: '宋晴晴', school: '北京少年宫', points: 2580, gold: 3, silver: 0 },
      { rank: 3, name: '梁思思', school: '广州艺术学校', points: 2430, gold: 2, silver: 2 },
      { rank: 4, name: '程小芸', school: '深圳少儿艺术团', points: 2260, gold: 2, silver: 1 },
      { rank: 5, name: '倪兰心', school: '杭州艺术学院', points: 2090, gold: 1, silver: 3 },
      { rank: 6, name: '贺欣欣', school: '成都少年宫', points: 1930, gold: 1, silver: 2 },
      { rank: 7, name: '施雨薇', school: '重庆艺术学校', points: 1780, gold: 0, silver: 4 },
      { rank: 8, name: '卢雅颜', school: '武汉舞蹈学院', points: 1640, gold: 0, silver: 2 }
    ],
    'U13': [
      { rank: 1, name: '林晓彤', school: '北京舞蹈学院附中', points: 3120, gold: 5, silver: 0 },
      { rank: 2, name: '王雅琪', school: '上海舞蹈学校', points: 2940, gold: 4, silver: 1 },
      { rank: 3, name: '刘梓嫣', school: '浙江艺术学校', points: 2760, gold: 3, silver: 2 },
      { rank: 4, name: '赵欣妍', school: '天津艺术学院', points: 2580, gold: 2, silver: 3 },
      { rank: 5, name: '许晓艺', school: '山东艺术学院', points: 2380, gold: 2, silver: 1 },
      { rank: 6, name: '崔梦瑶', school: '辽宁舞蹈学校', points: 2200, gold: 1, silver: 3 },
      { rank: 7, name: '姜诗颖', school: '吉林艺术学院', points: 2030, gold: 1, silver: 2 },
      { rank: 8, name: '谭雅文', school: '湖南艺术学校', points: 1870, gold: 0, silver: 4 }
    ],
    'U15': [
      { rank: 1, name: '黄思嘉', school: '深圳舞蹈学校', points: 3050, gold: 4, silver: 2 },
      { rank: 2, name: '李梦瑶', school: '杭州舞蹈团', points: 2890, gold: 4, silver: 0 },
      { rank: 3, name: '周沐阳', school: '上海舞蹈学校', points: 2720, gold: 3, silver: 2 },
      { rank: 4, name: '潘欣悦', school: '福建艺术学院', points: 2540, gold: 2, silver: 3 },
      { rank: 5, name: '丁思佳', school: '江苏艺术学校', points: 2350, gold: 2, silver: 1 },
      { rank: 6, name: '庄晓雯', school: '安徽艺术学院', points: 2170, gold: 1, silver: 3 },
      { rank: 7, name: '苗雨微', school: '广西艺术学院', points: 2000, gold: 1, silver: 2 },
      { rank: 8, name: '付小菁', school: '云南艺术学校', points: 1840, gold: 0, silver: 3 }
    ]
  }
};

module.exports = {
  STATUS,
  COMPETITIONS,
  LEADERBOARD
};
