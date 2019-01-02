package model

type TokenItem struct {
	Username string `json:"username"`
	DH_TOKEN string `json:"dh_token" xorm:"dh_token"`
}
