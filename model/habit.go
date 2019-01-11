package model

type Profile struct {
	Username    string   `json:"username"`           // 用户名
	Password    string   `json:"password" xorm:"->"` // 密码
	Avatar      []int8   `json:"avatar"`             // 用户头像
	Description string   `json:"description"`        // 用户描述
	Habits      []string `json:"habits"`             // 习惯id列表
}

type Habit struct {
	Id                  string `json:"id"`          // id
	TimeQuantum         string `json:"timeQuantum"` // 时间段
	Name                string `json:"name"`        // 名称
	Icon                string `json:"icon"`        // 图标
	File                bool   `json:"file"`        // 归档/结束
	Color               string `json:"color"`
	ReminderTime        string `json:"reminderTime"`        // 提醒时间
	Encourage           string `json:"encourage"`           // 激励语句
	Important           bool   `json:"important"`           // 是否重要
	Notification        bool   `json:"notification"`        // 是否通知
	RecentPunchTime     string `json:"recentPunchTime"`     // 最近打卡时间
	LastRecentPunchTime string `json:"lastRecentPunchTime"` // 上次打卡时间

	TotalPunch int    `json:"totalPunch"` // 总打卡数
	CurrcPunch int    `json:"currcPunch"` // 当前连续打卡数
	OncecPunch int    `json:"oncecPunch"` // 曾经最大连续打卡数
	CreateAt   string `json:"createAt"`   // 创建时间
}

type Month struct {
	Id          string   `json:"id"`          // 月份id `2018-12`
	PlanPunch   int      `json:"planPunch"`   // 计划打卡数
	ActualPunch int      `json:"actualPunch"` // 实际打卡数
	MissPunch   int      `json:"missPunch"`   // 错过打卡数
	Days        []string `json:"days"`        // 每天打卡id列表
}

type Day struct {
	Id   string `json:"id"`   // 日期id
	Time string `json:"time"` // 打卡时间`16:30`
	Log  string `json:"log"`  // 打卡信息
}
