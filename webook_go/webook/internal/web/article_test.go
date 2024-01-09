package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook_go/webook/internal/domain"
	"webook_go/webook/internal/service"
	svcmocks "webook_go/webook/internal/service/mocks"
	ijwt "webook_go/webook/internal/web/jwt"
	"webook_go/webook/pkg/logger"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) service.ArticleService
		reqBody string

		wantCode int
		wantRes  Result
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
{
	"title":"我的标题",
	"content":"我的内容"
}
`,
			wantCode: 200,
			wantRes: Result{
				Data: float64(1),
				Msg:  "OK",
			},
		},
		{
			name: "publish失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("publish error"))
				return svc
			},
			reqBody: `
{
	"title":"我的标题",
	"content":"我的内容"
}
`,
			wantCode: 200,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaims{
			Uid: 123,
		})
	})
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 用不上 codeSvc
			h := NewArticleHandler(tc.mock(ctrl), &logger.NopLogger{})
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/articles/publish",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			// 预期它肯定不会出问题
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
		})
	}
}
