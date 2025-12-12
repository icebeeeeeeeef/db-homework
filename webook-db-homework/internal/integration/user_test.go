package integration

/*
func TestUserHandler_SendLoginSMSCode(t *testing.T) {

	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				//不需要，因为redis里什么数据都不需要
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:13800138000").Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody:  `{"phone":"13800138000"}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 0,
				Msg:  "发送验证码成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:13800138000", "123456", time.Minute*10).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:13800138000").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, val, "123456")
			},
			reqBody:  `{"phone":"13800138000"}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 429,
				Msg:  "发送验证码太频繁",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			//准备好要发的请求
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)                            //代表着如果出错就会报错（panic）
			req.Header.Set("Content-Type", "application/json") //设置请求头是json格式
			resp := httptest.NewRecorder()

			//这就是HTTP请求进入gin的入口
			//当你这样调用的时候gin就会处理这个请求
			//随后把返回值写回到resp中
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			var webResult web.Result
			err = json.NewDecoder(resp.Body).Decode(&webResult)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webResult)
			tc.after(t)
		})
	}
}
*/
