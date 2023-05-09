package models

type Like struct {
	Id    uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name  string `gorm:"size:36;not null;unique" json:"name"`
	Count int32  `gorm:"not null" json:"count"`
	Image string `gorm:"" json:"image"`
}
