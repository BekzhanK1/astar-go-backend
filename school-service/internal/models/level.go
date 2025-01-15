package models

type Level struct {
	Name   string `gorm:"not null;index:idx_name_flow,unique" json:"name"`
	FlowID uint   `gorm:"not null;index:idx_name_flow,unique" json:"flow_id"`
	Flow   Flow   `gorm:"foreignKey:FlowID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"flow"`
}
