package web

import (
	"bytes"
	"geek-hw-week5/internal/domain"
	"geek-hw-week5/internal/service"
	svcmocks "geek-hw-week5/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name string

		mockBuildServices func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		buildReq func(t *testing.T) *http.Request

		expectedCode int
		expectedBody string
	}{
		{
			name: "注册成功",
			mockBuildServices: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "5578@qq.com",
					Password: "hello@world123!",
				}).Return(nil)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			buildReq: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
						"email": "5578@qq.com",
						"password": "hello@world123!",
						"confirmPassword": "hello@world123!"
					}`)),
				)
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedCode: http.StatusOK,
			expectedBody: "注册成功",
		},
		{
			name: "非法邮箱格式",
			mockBuildServices: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			buildReq: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
						"email": "5578@qq",
						"password": "hello@world123!",
						"confirmPassword": "hello@world123!"
					}`)),
				)
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "非法邮箱格式",
		},
		{
			name: "两次输入密码不对",
			mockBuildServices: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			buildReq: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
						"email": "5578@qq.com",
						"password": "hello@world123!",
						"confirmPassword": "hello@world123"
					}`)),
				)
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "两次输入密码不对",
		},
		{
			name: "密码格式错误",
			mockBuildServices: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			buildReq: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
						"email": "5578@qq.com",
						"password": "hello@world",
						"confirmPassword": "hello@world"
					}`)),
				)
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "密码必须包含字母、数字、特殊字符，并且不少于八位",
		},
		{
			name: "邮箱冲突",
			mockBuildServices: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "5578@qq.com",
					Password: "hello@world123!",
				}).Return(service.ErrDuplicateEmail)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			buildReq: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
						"email": "5578@qq.com",
						"password": "hello@world123!",
						"confirmPassword": "hello@world123!"
					}`)),
				)
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "邮箱冲突，请换一个",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := testCase.mockBuildServices(ctrl)
			handler := NewUserHandler(userSvc, codeSvc)

			gin.SetMode(gin.ReleaseMode)
			server := gin.Default()
			handler.RegisterRoutes(server)

			req := testCase.buildReq(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedCode, recorder.Code)
			assert.Equal(t, testCase.expectedBody, recorder.Body.String())
		})
	}
}
