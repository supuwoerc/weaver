package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock localizer
type mockLocalizer struct {
	localizeFunc func(*i18n.LocalizeConfig) (string, error)
}

func (m *mockLocalizer) MustLocalize(config *i18n.LocalizeConfig) string {
	result, _ := m.localizeFunc(config)
	return result
}

func TestHttpResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		code         StatusCode
		data         interface{}
		config       *i18n.LocalizeConfig
		message      *string
		setupContext func(*gin.Context)
		expectedCode int
		expectedData interface{}
		expectedMsg  string
	}{
		{
			name:         "success with data and message",
			code:         Ok,
			data:         map[string]string{"key": "value"},
			config:       nil,
			message:      lo.ToPtr("Custom message"),
			setupContext: nil,
			expectedCode: 10000,
			expectedData: map[string]interface{}{"key": "value"},
			expectedMsg:  "Custom message",
		},
		{
			name:         "success without message",
			code:         Ok,
			data:         "test data",
			config:       nil,
			message:      nil,
			setupContext: nil,
			expectedCode: 10000,
			expectedData: "test data",
			expectedMsg:  "",
		},
		{
			name:    "with i18n translator",
			code:    Ok,
			data:    nil,
			config:  nil,
			message: nil,
			setupContext: func(c *gin.Context) {
				// 模拟 i18n localizer
				localizer := &mockLocalizer{
					localizeFunc: func(config *i18n.LocalizeConfig) (string, error) {
						return "localized message", nil
					},
				}
				c.Set(string(I18nTranslatorKey), localizer)
			},
			expectedCode: 10000,
			expectedData: nil,
			expectedMsg:  "localized message",
		},
		{
			name:    "with custom localize config",
			code:    Error,
			data:    nil,
			config:  &i18n.LocalizeConfig{MessageID: "customError"},
			message: nil,
			setupContext: func(c *gin.Context) {
				localizer := &mockLocalizer{
					localizeFunc: func(config *i18n.LocalizeConfig) (string, error) {
						if config.MessageID == "customError" {
							return "custom error message", nil
						}
						return "default error message", nil
					},
				}
				c.Set(string(I18nTranslatorKey), localizer)
			},
			expectedCode: 10001,
			expectedData: nil,
			expectedMsg:  "custom error message",
		},
		{
			name:         "error response",
			code:         InvalidParams,
			data:         nil,
			config:       nil,
			message:      lo.ToPtr("Invalid parameters"),
			setupContext: nil,
			expectedCode: 10002,
			expectedData: nil,
			expectedMsg:  "Invalid parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试响应记录器和上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			// 设置上下文
			if tt.setupContext != nil {
				tt.setupContext(c)
			}
			// 调用 HttpResponse
			HttpResponse(c, tt.code, tt.data, tt.config, tt.message)
			// 验证 HTTP 状态码
			assert.Equal(t, http.StatusOK, w.Code)
			var response BasicResponse[any]
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)
			// 验证响应内容
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, tt.expectedData, response.Data)
			assert.Equal(t, tt.expectedMsg, response.Message)
		})
	}
}

func TestSuccess(t *testing.T) {
	t.Run("success without data", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		Success(ctx)
		assert.Equal(t, http.StatusOK, w.Code)
		var response BasicResponse[any]
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, int(Ok), response.Code)
		assert.Nil(t, response.Data)
		assert.Equal(t, "", response.Message)
	})
}

func TestSuccessWithData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		data         interface{}
		expectedData interface{}
	}{
		{
			name:         "string data",
			data:         "test string",
			expectedData: "test string",
		},
		{
			name:         "map data",
			data:         map[string]interface{}{"key": "value", "count": 42},
			expectedData: map[string]interface{}{"key": "value", "count": float64(42)}, // JSON 数字会被解析为 float64
		},
		{
			name: "struct data",
			data: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{ID: 1, Name: "test"},
			expectedData: map[string]interface{}{"id": float64(1), "name": "test"},
		},
		{
			name:         "slice data",
			data:         []string{"a", "b", "c"},
			expectedData: []interface{}{"a", "b", "c"},
		},
		{
			name:         "nil data",
			data:         nil,
			expectedData: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			SuccessWithData(c, tt.data)
			// 验证 HTTP 状态码
			assert.Equal(t, http.StatusOK, w.Code)
			// 解析响应 JSON
			var response BasicResponse[any]
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)
			// 验证响应内容
			assert.Equal(t, int(Ok), response.Code)
			assert.Equal(t, tt.expectedData, response.Data)
			assert.Equal(t, "", response.Message)
		})
	}
}

func TestSuccessWithMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		message         string
		expectedMessage string
	}{
		{
			name:            "normal message",
			message:         "operation successful",
			expectedMessage: "operation successful",
		},
		{
			name:            "empty message",
			message:         "",
			expectedMessage: "",
		},
		{
			name:            "message with special characters",
			message:         "success! 用户创建成功 😊",
			expectedMessage: "success! 用户创建成功 😊",
		},
		{
			name:            "long message",
			message:         "this is a very long success message that contains multiple words and should be handled properly by the response system",
			expectedMessage: "this is a very long success message that contains multiple words and should be handled properly by the response system",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			SuccessWithMessage(c, tt.message)
			assert.Equal(t, http.StatusOK, w.Code)
			// 解析响应 JSON
			var response BasicResponse[any]
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)
			// 验证响应内容
			assert.Equal(t, int(Ok), response.Code)
			assert.Nil(t, response.Data)
			assert.Equal(t, tt.expectedMessage, response.Message)
		})
	}
}

func TestSuccessWithPageData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		total        int64
		list         interface{}
		expectedData DataList[interface{}]
	}{
		{
			name:  "normal page data",
			total: 100,
			list:  []string{"item1", "item2", "item3"},
			expectedData: DataList[interface{}]{
				Total: 100,
				List:  []interface{}{"item1", "item2", "item3"},
			},
		},
		{
			name:  "empty list",
			total: 0,
			list:  []string{},
			expectedData: DataList[interface{}]{
				Total: 0,
				List:  []interface{}{},
			},
		},
		{
			name:  "nil list",
			total: 0,
			list:  nil,
			expectedData: DataList[interface{}]{
				Total: 0,
				List:  nil,
			},
		},
		{
			name:  "large total",
			total: 999999,
			list:  []int{1, 2, 3, 4, 5},
			expectedData: DataList[interface{}]{
				Total: 999999,
				List:  []interface{}{float64(1), float64(2), float64(3), float64(4), float64(5)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			switch v := tt.list.(type) {
			case []string:
				SuccessWithPageData(c, tt.total, v)
			case []int:
				SuccessWithPageData(c, tt.total, v)
			case nil:
				SuccessWithPageData(c, tt.total, []string(nil))
			default:
				SuccessWithPageData(c, tt.total, []interface{}{})
			}
			assert.Equal(t, http.StatusOK, w.Code)

			var response BasicResponse[DataList[interface{}]]
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)
			// 验证响应内容
			assert.Equal(t, int(Ok), response.Code)
			assert.Equal(t, tt.expectedData.Total, response.Data.Total)
			assert.Equal(t, tt.expectedData.List, response.Data.List)
			assert.Equal(t, "", response.Message)
		})
	}
}

func TestSuccessWithPageDataGenericTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("struct list", func(t *testing.T) {
		type User struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		users := []User{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		}
		total := int64(2)
		SuccessWithPageData(c, total, users)
		assert.Equal(t, http.StatusOK, w.Code)
		var response BasicResponse[DataList[User]]
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		// 验证响应内容
		assert.Equal(t, int(Ok), response.Code)
		assert.Equal(t, total, response.Data.Total)
		assert.Equal(t, users, response.Data.List)
		assert.Equal(t, "", response.Message)
	})

	t.Run("pointer list", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		numbers := []*int{lo.ToPtr(1), lo.ToPtr(2), lo.ToPtr(3)}
		total := int64(3)
		SuccessWithPageData(c, total, numbers)
		assert.Equal(t, http.StatusOK, w.Code)
		// 解析响应 JSON
		var response BasicResponse[DataList[*int]]
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		// 验证响应内容
		assert.Equal(t, int(Ok), response.Code)
		assert.Equal(t, total, response.Data.Total)
		assert.Len(t, response.Data.List, 3)
		assert.Equal(t, "", response.Message)
	})
}

func TestFailWithMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		message         string
		expectedMessage string
	}{
		{
			name:            "normal error message",
			message:         "something went wrong",
			expectedMessage: "something went wrong",
		},
		{
			name:            "empty message",
			message:         "",
			expectedMessage: "",
		},
		{
			name:            "message with special characters",
			message:         "error! 系统错误 😭",
			expectedMessage: "error! 系统错误 😭",
		},
		{
			name:            "long error message",
			message:         "this is a very long error message that describes what went wrong in great detail and should be handled properly",
			expectedMessage: "this is a very long error message that describes what went wrong in great detail and should be handled properly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			FailWithMessage(c, tt.message)
			assert.Equal(t, http.StatusOK, w.Code)
			// 解析响应 JSON
			var response BasicResponse[interface{}]
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)
			// 验证响应内容
			assert.Equal(t, int(Error), response.Code)
			assert.Nil(t, response.Data)
			assert.Equal(t, tt.expectedMessage, response.Message)
		})
	}
}

func TestFailWithCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		code         StatusCode
		expectedCode int
	}{
		{
			name:         "invalid params",
			code:         InvalidParams,
			expectedCode: 10002,
		},
		{
			name:         "invalid token",
			code:         InvalidToken,
			expectedCode: 10003,
		},
		{
			name:         "user not exist",
			code:         UserNotExist,
			expectedCode: 20005,
		},
		{
			name:         "permission not exist",
			code:         PermissionNotExist,
			expectedCode: 50001,
		},
		{
			name:         "timeout error",
			code:         TimeoutErr,
			expectedCode: 10010,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			FailWithCode(c, tt.code)
			assert.Equal(t, http.StatusOK, w.Code)
			// 解析响应 JSON
			var response BasicResponse[interface{}]
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)
			// 验证响应内容
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Nil(t, response.Data)
			assert.Equal(t, "", response.Message)
		})
	}
}

func TestFailWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "status code error",
			err:          UserNotExist,
			expectedCode: 20005,
			expectedMsg:  "",
		},
		{
			name:         "context canceled",
			err:          context.Canceled,
			expectedCode: 10004,
			expectedMsg:  "",
		},
		{
			name:         "context deadline exceeded",
			err:          context.DeadlineExceeded,
			expectedCode: 10010,
			expectedMsg:  "",
		},
		{
			name:         "wrapped status code error",
			err:          fmt.Errorf("wrapped error: %w", InvalidToken),
			expectedCode: 10003,
			expectedMsg:  "",
		},
		{
			name:         "normal error",
			err:          errors.New("something went wrong"),
			expectedCode: 10001,
			expectedMsg:  "something went wrong",
		},
		{
			name:         "wrapped context canceled",
			err:          fmt.Errorf("request failed: %w", context.Canceled),
			expectedCode: 10004,
			expectedMsg:  "",
		},
		{
			name:         "wrapped context deadline exceeded",
			err:          fmt.Errorf("timeout: %w", context.DeadlineExceeded),
			expectedCode: 10010,
			expectedMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			FailWithError(c, tt.err)
			// 验证 HTTP 状态码
			assert.Equal(t, http.StatusOK, w.Code)
			// 解析响应 JSON
			var response BasicResponse[interface{}]
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			// 验证响应内容
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Nil(t, response.Data)
			assert.Equal(t, tt.expectedMsg, response.Message)
		})
	}
}

// Mock validator.FieldError
type mockFieldError struct {
	field       string
	translation string
}

func (m *mockFieldError) Tag() string             { return "required" }
func (m *mockFieldError) ActualTag() string       { return "required" }
func (m *mockFieldError) Namespace() string       { return "Test." + m.field }
func (m *mockFieldError) StructNamespace() string { return "Test." + m.field }
func (m *mockFieldError) Field() string           { return m.field }
func (m *mockFieldError) StructField() string     { return m.field }
func (m *mockFieldError) Value() interface{}      { return nil }
func (m *mockFieldError) Param() string           { return "" }
func (m *mockFieldError) Kind() reflect.Kind      { return reflect.String }
func (m *mockFieldError) Type() reflect.Type      { return reflect.TypeOf("") }
func (m *mockFieldError) Translate(ut ut.Translator) string {
	return m.translation
}
func (m *mockFieldError) Error() string { return m.translation }

type mockTranslator struct {
	ut.Translator
}

func TestParamsValidateFail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		err           error
		setupContext  func(*gin.Context)
		expectedCode  int
		expectedData  interface{}
		expectedMsg   string
		checkDataType func(interface{}) bool
	}{
		{
			name:          "normal error (not ValidationErrors)",
			err:           errors.New("invalid parameter"),
			setupContext:  nil,
			expectedCode:  10002,
			expectedData:  nil,
			expectedMsg:   "invalid parameter",
			checkDataType: func(data interface{}) bool { return data == nil },
		},
		{
			name: "validation errors without translator",
			err: validator.ValidationErrors{
				&mockFieldError{field: "Name", translation: "Name is required"},
				&mockFieldError{field: "Email", translation: "Email is invalid"},
			},
			setupContext:  nil,
			expectedCode:  10002,
			expectedData:  nil,
			expectedMsg:   "Name is required\nEmail is invalid",
			checkDataType: func(data interface{}) bool { return data == nil },
		},
		{
			name: "validation errors with invalid translator type",
			err: validator.ValidationErrors{
				&mockFieldError{field: "Name", translation: "Name is required"},
			},
			setupContext: func(c *gin.Context) {
				c.Set(string(ValidatorTranslatorKey), "not a translator")
			},
			expectedCode:  10002,
			expectedData:  nil,
			expectedMsg:   "Name is required",
			checkDataType: func(data interface{}) bool { return data == nil },
		},
		{
			name: "validation errors with valid translator",
			err: validator.ValidationErrors{
				&mockFieldError{field: "Name", translation: "Name field is required"},
				&mockFieldError{field: "Email", translation: "Email format is invalid"},
			},
			setupContext: func(c *gin.Context) {
				translator := &mockTranslator{}
				c.Set(string(ValidatorTranslatorKey), translator)
			},
			expectedCode: 10002,
			expectedData: map[string]string{
				"Name":  "Name field is required",
				"Email": "Email format is invalid",
			},
			expectedMsg: "",
			checkDataType: func(data interface{}) bool {
				errMap, ok := data.(map[string]interface{}) // 反序列化为map[string]interface{}
				return ok && len(errMap) == 2
			},
		},
		{
			name: "validation errors with namespace field names",
			err: validator.ValidationErrors{
				&mockFieldError{field: "User[0].Name", translation: "User name is required"},
				&mockFieldError{field: "Items[1].Price", translation: "Price must be positive"},
			},
			setupContext: func(c *gin.Context) {
				translator := &mockTranslator{}
				c.Set(string(ValidatorTranslatorKey), translator)
			},
			expectedCode: 10002,
			expectedData: map[string]string{
				"User.Name":   "User name is required",
				"Items.Price": "Price must be positive",
			},
			expectedMsg: "",
			checkDataType: func(data interface{}) bool {
				errMap, ok := data.(map[string]interface{}) // 反序列化为map[string]interface{}
				return ok && len(errMap) == 2
			},
		},
		{
			name: "validation errors with regex replacement failure",
			err: validator.ValidationErrors{
				&mockFieldError{field: "InvalidRegexField[", translation: "Invalid field"},
			},
			setupContext: func(c *gin.Context) {
				translator := &mockTranslator{}
				c.Set(string(ValidatorTranslatorKey), translator)
			},
			expectedCode: 10002,
			expectedData: map[string]string{
				"InvalidRegexField[": "Invalid field",
			},
			expectedMsg: "",
			checkDataType: func(data interface{}) bool {
				errMap, ok := data.(map[string]interface{}) // 反序列化为map[string]interface{}
				return ok && len(errMap) == 1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.setupContext != nil {
				tt.setupContext(c)
			}
			ParamsValidateFail(c, tt.err)
			assert.Equal(t, http.StatusOK, w.Code)
			// 解析响应 JSON
			var response BasicResponse[any]
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)
			// 验证响应内容
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, tt.expectedMsg, response.Message)
			// 验证数据内容
			if tt.checkDataType != nil {
				require.True(t, tt.checkDataType(response.Data), "Data type check failed")
			}
			// 对于有具体期望数据的测试，进行详细比较
			if tt.expectedData != nil {
				if expectedMap, ok := tt.expectedData.(map[string]string); ok {
					responseMap, ok := response.Data.(map[string]interface{})
					require.True(t, ok, "Response data should be a map")

					assert.Equal(t, len(expectedMap), len(responseMap), "Map sizes should match")

					for key, expectedValue := range expectedMap {
						actualValue, exists := responseMap[key]
						assert.True(t, exists, "Key %s should exist", key)
						assert.Equal(t, expectedValue, actualValue, "Value for key %s should match", key)
					}
				} else {
					assert.Equal(t, tt.expectedData, response.Data)
				}
			}
		})
	}
}
