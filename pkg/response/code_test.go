package response

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusCode_Error(t *testing.T) {
	tests := []struct {
		name string
		s    StatusCode
		want string
	}{
		{
			name: "response code ok",
			s:    Ok,
			want: "10000",
		},
		{
			name: "response code error",
			s:    Error,
			want: "10001",
		},
		{
			name: "response code invalidParams",
			s:    InvalidParams,
			want: "10002",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("customer statue code", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			invalidStatusCode := StatusCode(i)
			err := invalidStatusCode.Error()
			assert.Contains(t, err, fmt.Sprintf("%d", i))
		}
	})
}
