package model

type DailyCommit struct {
	Id            string `json:"id"`            // id
	CommitTime    string `json:"commitTime"`    // 提交时间
	CommitContent string `json:"commitContent"` // 提交内容
}
