package models

func (Species) TableName() string {
	return "Species"
}

type Species struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"type:varchar(100);not null;unique"`
}
