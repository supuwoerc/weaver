package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supuwoerc/weaver/pkg/database"
)

func TestPermission_GetRoleIds(t *testing.T) {
	p := &Permission{
		Roles: []*Role{
			{
				BasicModel: database.BasicModel{
					ID: 1,
				},
			},
			{
				BasicModel: database.BasicModel{
					ID: 2,
				},
			},
			{
				BasicModel: database.BasicModel{
					ID: 3,
				},
			},
		},
	}
	t.Run("Get Role Ids", func(t *testing.T) {
		ids := p.GetRoleIds()
		assert.Equal(t, ids[0], uint(1))
		assert.Equal(t, ids[1], uint(2))
		assert.Equal(t, ids[2], uint(3))
		assert.Equal(t, len(ids), len(p.Roles))
	})
}
