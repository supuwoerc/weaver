package response

import (
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/supuwoerc/weaver/models"
)

// DepartmentTreeResponse 全量部门树结构
type DepartmentTreeResponse struct {
	*models.Department
	FullName  string                    `json:"full_name"`
	FullIds   []uint                    `json:"full_ids"`
	Ancestors []uint                    `json:"ancestors,omitempty"`
	Parent    any                       `json:"parent,omitempty"`
	Children  []*DepartmentTreeResponse `json:"children,omitempty"`
	Leaders   []*SimpleUser             `json:"leaders,omitempty"`
	Users     []*SimpleUser             `json:"users,omitempty"`
	Creator   *Creator                  `json:"creator,omitempty"`
	Updater   *Updater                  `json:"updater,omitempty"`
}

// ToDepartmentTreeResponse 将 dept 转为响应
func ToDepartmentTreeResponse(dept *models.Department, deptMap map[uint]*models.Department) (*DepartmentTreeResponse, error) {
	var ancestors string
	if dept.Ancestors != nil {
		ancestors = *dept.Ancestors
	}
	var splitAncestors = make([]uint, 0)
	if ancestors != "" {
		ids := strings.Split(ancestors, ",")
		for _, idString := range ids {
			atoi, err := strconv.Atoi(idString)
			if err != nil {
				return nil, err
			}
			splitAncestors = append(splitAncestors, uint(atoi))
		}
	}
	ancestorNames := lo.Map(splitAncestors, func(item uint, _ int) string {
		return deptMap[item].Name
	})
	ancestorNames = append(ancestorNames, dept.Name)
	fullName := strings.Join(ancestorNames, "/")
	fullIds := make([]uint, len(splitAncestors))
	copy(fullIds, splitAncestors)
	fullIds = append(fullIds, dept.ID)
	res := &DepartmentTreeResponse{
		Department: dept,
		Ancestors:  splitAncestors,
		FullName:   fullName,
		FullIds:    fullIds,
		Leaders: lo.Map(dept.Leaders, func(item *models.User, _ int) *SimpleUser {
			return &SimpleUser{
				User: item,
			}
		}),
		Users: lo.Map(dept.Users, func(item *models.User, _ int) *SimpleUser {
			return &SimpleUser{
				User: item,
			}
		}),
	}
	if dept.Creator.ID != 0 {
		res.Creator = &Creator{
			User: &dept.Creator,
		}
	}
	if dept.Updater.ID != 0 {
		res.Updater = &Updater{
			User: &dept.Updater,
		}
	}
	return res, nil
}
