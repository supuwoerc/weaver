package models

import (
	"encoding/json"

	"github.com/supuwoerc/weaver/pkg/database"
)

type Department struct {
	Name      string        `json:"name" gorm:"not null;"`
	ParentID  *uint         `json:"parent_id"` // 父级部门ID,nil则为顶级部门
	Parent    *Department   `json:"parent" gorm:"foreignKey:ParentID;references:ID"`
	Children  []*Department `json:"children" gorm:"foreignKey:ParentID"` // 不建议使用 gorm 预加载
	Ancestors *string       `json:"ancestors"`                           // 祖先部门路径逗号拼接的字符串
	Leaders   []*User       `json:"leaders" gorm:"many2many:department_leader;"`
	Users     []*User       `json:"users" gorm:"many2many:user_department;"`
	CreatorID uint          `json:"-" gorm:"not null;"`
	Creator   User          `json:"creator" gorm:"foreignKey:CreatorID;references:ID"`
	UpdaterID uint          `json:"-" gorm:"not null;"`
	Updater   User          `json:"updater" gorm:"foreignKey:UpdaterID;references:ID"`
	database.BasicModel
}

// MarshalBinary 实现 encoding.BinaryMarshaler 接口
func (d *Department) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

// UnmarshalBinary 实现 encoding.BinaryUnmarshaler 接口
func (d *Department) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

type Departments []*Department

func (dl Departments) MarshalBinary() ([]byte, error) {
	return json.Marshal(dl)
}

func (dl Departments) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &dl)
}
