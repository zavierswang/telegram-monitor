package models

type Address struct {
	ID        int    `gorm:"column:id;autoIncrement;primaryKey"`
	UserID    string `gorm:"column:user_id;not null"`
	Username  string `gorm:"column:username"`
	Address   string `gorm:"column:address;not null"`
	IsMonitor bool   `gorm:"column:is_monitor;default:false"`
	Group     Group  `gorm:"embedded;embeddedPrefix:group_"`
	Remark    string `gorm:"column:remark;default:null"`
	Avator    string `gorm:"column:avator;default:null"`
}

type Group struct {
	ChatID   string `gorm:"default:null"`
	Username string `gorm:"default:null;size:64"`
	Title    string `gorm:"default:null;size:64"`
}

func (a *Address) TableName() string {
	return "tb_address"
}
