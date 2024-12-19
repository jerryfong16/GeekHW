package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	integrationteststartup "geek-hw-week5/internal/integration_test/startup"
	"geek-hw-week5/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestSendSMSLoginCode(t *testing.T) {
	server := integrationteststartup.InitWebServer()
	redis := integrationteststartup.InitRedis()
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		phone string

		expectedCode   int
		expectedResult web.Result
	}{
		{
			name: "发送成功",

			before: func(t *testing.T) {},

			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13166668888"
				code, err := redis.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := redis.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				err = redis.Del(ctx, key).Err()
				assert.NoError(t, err)
			},

			phone: "13166668888",

			expectedCode:   http.StatusOK,
			expectedResult: web.Result{Msg: "发送成功"},
		},
		{
			name: "发送太频繁",

			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13166668888"
				err := redis.Set(ctx, key, "123456", time.Minute*9+time.Second*5).Err()
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13166668888"
				code, err := redis.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, code, "123456")
			},

			phone: "13166668888",

			expectedCode: http.StatusOK,
			expectedResult: web.Result{
				Code: 4,
				Msg:  "短信发送太频繁，请稍后再试",
			},
		},
		{
			name: "系统错误",

			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13166668888"
				err := redis.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13166668888"
				code, err := redis.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, code, "123456")
			},

			phone: "13166668888",

			expectedCode: http.StatusOK,
			expectedResult: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.before(t)
			defer testCase.after(t)

			req, err := http.NewRequest(
				http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone":"%s"}`, testCase.phone))),
			)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)
			assert.Equal(t, testCase.expectedCode, recorder.Code)
			if testCase.expectedCode != http.StatusOK {
				return
			}

			var rs web.Result
			err = json.NewDecoder(recorder.Body).Decode(&rs)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResult, rs)
		})
	}
}
