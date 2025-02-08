package response

import "gin-web/models"

// DepartmentTreeResponse 全量部门树结构
type DepartmentTreeResponse struct {
	*models.Department
	Users       any     `json:"users,omitempty"`
	Permissions any     `json:"permissions,omitempty"`
	Creator     Creator `json:"creator"`
	Updater     Updater `json:"updater"`
}
