package service

import (
	"context"
	"testing"
	"webook/internal/domain"
	events "webook/internal/events/article"
	repository "webook/internal/repository/article"
	"webook/pkg/logger"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestArticleService_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (repository.ArticleRepository, events.Producer)
		reqBody domain.Article
		wantID  int64
		wantErr error
	}{ /*
			{
				name: "新建并发布成功",
				mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
					repo := repomocks.NewMockArticleAuthorRepository(ctrl)
					repo.EXPECT().Create(gomock.Any(), domain.Article{
						Title:   "test",
						Content: "test123",
						Author: domain.Author{
							ID: 666,
						},
					}).Return((int64)(1), nil)
					repo1 := repomocks.NewMockArticleReaderRepository(ctrl)
					repo1.EXPECT().Save(gomock.Any(), domain.Article{
						ID:      1,
						Title:   "test",
						Content: "test123",
						Author: domain.Author{
							ID: 666,
						},
					}).Return(nil)
					return repo, repo1
				},
				reqBody: domain.Article{
					Title:   "test",
					Content: "test123",
					Author: domain.Author{
						ID: 666,
					},
				},
				wantID:  1,
				wantErr: nil,
			},
			{
				name: "修改并发布成功",
				mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
					repo := repomocks.NewMockArticleAuthorRepository(ctrl)
					repo.EXPECT().Update(gomock.Any(), domain.Article{
						ID:      2,
						Title:   "test",
						Content: "test123",
						Author: domain.Author{
							ID: 666,
						},
					}).Return(nil)
					repo1 := repomocks.NewMockArticleReaderRepository(ctrl)
					repo1.EXPECT().Save(gomock.Any(), domain.Article{
						ID:      2,
						Title:   "test",
						Content: "test123",
						Author: domain.Author{
							ID: 666,
						},
					}).Return(nil)
					return repo, repo1
				},
				reqBody: domain.Article{
					ID:      2,
					Title:   "test",
					Content: "test123",
					Author: domain.Author{
						ID: 666,
					},
				},
				wantID:  2,
				wantErr: nil,
			},
			{
				name: "制作库失败（新建）",
				mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
					repo := repomocks.NewMockArticleAuthorRepository(ctrl)
					repo.EXPECT().Create(gomock.Any(), domain.Article{
						Title:   "test",
						Content: "test123",
						Author: domain.Author{
							ID: 666,
						},
					}).Return((int64)(0), errors.New("制作库失败"))
					repo1 := repomocks.NewMockArticleReaderRepository(ctrl)
					return repo, repo1
				},
				reqBody: domain.Article{
					Title:   "test",
					Content: "test123",
					Author: domain.Author{
						ID: 666,
					},
				},
				wantID:  0,
				wantErr: errors.New("制作库失败"),
			},

			{
				name: "制作库成功（新建），部分重试后成功",
				mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
					repo := repomocks.NewMockArticleAuthorRepository(ctrl)
					repo.EXPECT().Create(gomock.Any(), domain.Article{
						Title:   "test",
						Content: "test123",
						Author: domain.Author{
							ID: 666,
						},
					}).Return((int64)(1), nil)

					repo1 := repomocks.NewMockArticleReaderRepository(ctrl)
					gomock.InOrder(
						repo1.EXPECT().Save(gomock.Any(), domain.Article{
							ID:      1,
							Title:   "test",
							Content: "test123",
							Author: domain.Author{
								ID: 666,
							},
						}).Return(errors.New("发布到线上库失败")),
						repo1.EXPECT().Save(gomock.Any(), domain.Article{
							ID:      1,
							Title:   "test",
							Content: "test123",
							Author: domain.Author{
								ID: 666,
							},
						}).Return(nil))
					return repo, repo1
				},
				reqBody: domain.Article{
					Title:   "test",
					Content: "test123",
					Author: domain.Author{
						ID: 666,
					},
				},
				wantID:  1,
				wantErr: nil,
			},

			{
				name: "发布到线上库的时候重试结束都没成功（新建）",
				mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
					repo := repomocks.NewMockArticleAuthorRepository(ctrl)
					repo.EXPECT().Create(gomock.Any(), domain.Article{
						Title:   "test",
						Content: "test123",
						Author: domain.Author{
							ID: 666,
						},
					}).Return((int64)(1), nil)
					repo1 := repomocks.NewMockArticleReaderRepository(ctrl)
					repo1.EXPECT().Save(gomock.Any(), domain.Article{
						ID:      1,
						Title:   "test",
						Content: "test123",
						Author: domain.Author{
							ID: 666,
						},
					}).Return(errors.New("发布到线上库失败")).Times(3)
					return repo, repo1
				},
				reqBody: domain.Article{
					Title:   "test",
					Content: "test123",
					Author: domain.Author{
						ID: 666,
					},
				},
				wantID:  0,
				wantErr: errors.New("发布到线上库失败"),
			},
		*/ //这里这个测试用例是引入两个repo的版本，但是那个版本的针对与grpc的客户端改造没修改，这里注释掉了
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo, producer := tc.mock(ctrl)
			svc := NewArticleService(repo, logger.NewZapLogger(zap.NewExample()), producer)
			id, err := svc.Publish(context.Background(), tc.reqBody)
			assert.Equal(t, tc.wantID, id)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}
