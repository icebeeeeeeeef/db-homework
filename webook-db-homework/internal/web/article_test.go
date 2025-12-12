package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"webook/internal/domain"
	"webook/internal/service"
	svcmocks "webook/internal/service/mocks"
	ijwt "webook/internal/web/jwt"

	"webook/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantBody Result[int64]
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				articleSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "test",
					Content: "test123",
					Author: domain.Author{
						ID: 666,
					},
				}).Return((int64)(1), nil)
				return articleSvc
			},
			reqBody: `{
				"title": "test",
				"content": "test123"
			}`,
			wantCode: http.StatusOK,
			wantBody: Result[int64]{
				Code: 0,
				Msg:  "新建并发表成功",
				Data: 1,
			},
		},
		{
			name: "发布失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				articleSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "test",
					Content: "test123",
					Author: domain.Author{
						ID: 666,
					},
				}).Return((int64)(0), errors.New("发布失败"))
				return articleSvc
			},
			reqBody: `{
				"title": "test",
				"content": "test123"
			}`,
			wantCode: http.StatusOK,
			wantBody: Result[int64]{
				Code: 500,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()

			server.Use(func(c *gin.Context) {
				c.Set("claims", &ijwt.UserClaims{
					UserId: 666, //直接在这里模拟用户的登录
				})
			})

			h := NewArticleHandler(tc.mock(ctrl), logger.NewZapLogger(zap.NewExample()))

			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			var webResult Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webResult)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webResult)
		})
	}

}
