package models

type Category struct {
	ID       string    `gorm:"column:id;type:uuid"`
	ParentId *string   `gorm:"column:parent_id;type:uuid"`
	Parent   *Category `gorm:"references:ID;foreignKey:ParentId"`
	Name     string    `gorm:"column:name"`
	Meta
}

func (Category) TableName() string {
	return "categories.categories"
}

