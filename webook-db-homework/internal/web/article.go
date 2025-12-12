package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"webook/internal/domain"
	"webook/internal/service"
	"webook/pkg/logger"
)

type ArticleHandler struct {
	svc         service.ArticleService
	interactive service.InteractiveService
	l           logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService, inter service.InteractiveService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc:         svc,
		interactive: inter,
		l:           l,
	}
}

func (h *ArticleHandler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/articles")
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)
	g.POST("/list", h.List)
	g.GET("/detail/:id", h.Detail)
	g.GET("/pub/:id", h.PubDetail)

	pub := r.Group("/articles/pub")
	pub.POST("/list", h.PubList)
	pub.POST("/like", h.PubLike)
	pub.POST("/collect", h.PubCollect)
	pub.POST("/reward", h.PubReward)

	r.POST("/reward/detail", h.RewardDetail)
}

func (h *ArticleHandler) Edit(c *gin.Context) {
	var req ArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Result[int64]{Code: 400, Msg: "参数错误"})
		return
	}
	uid := c.GetInt64("userId")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, Result[int64]{Code: 401, Msg: "未登录"})
		return
	}
	id, err := h.svc.Save(c, req.ToDomain(uid))
	if err != nil {
		h.l.Error("保存文章失败", logger.Error(err))
		c.JSON(http.StatusInternalServerError, Result[int64]{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result[int64]{Code: 0, Msg: "编辑成功", Data: id})
}

func (h *ArticleHandler) Publish(c *gin.Context) {
	var req ArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Result[int64]{Code: 400, Msg: "参数错误"})
		return
	}
	uid := c.GetInt64("userId")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, Result[int64]{Code: 401, Msg: "未登录"})
		return
	}
	id, err := h.svc.Publish(c, req.ToDomain(uid))
	if err != nil {
		h.l.Error("发表文章失败", logger.Error(err))
		c.JSON(http.StatusInternalServerError, Result[int64]{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result[int64]{Code: 0, Msg: "发表成功", Data: id})
}

func (h *ArticleHandler) Withdraw(c *gin.Context) {
	var req ArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Result[int64]{Code: 400, Msg: "参数错误"})
		return
	}
	uid := c.GetInt64("userId")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, Result[int64]{Code: 401, Msg: "未登录"})
		return
	}
	err := h.svc.Withdraw(c, req.ToDomain(uid))
	if err != nil {
		h.l.Error("撤回文章失败", logger.Error(err))
		c.JSON(http.StatusInternalServerError, Result[int64]{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result[int64]{Code: 0, Msg: "撤回成功", Data: req.ID})
}

func (h *ArticleHandler) List(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Result[[]ArticleVO]{Code: 400, Msg: "参数错误"})
		return
	}
	uid := c.GetInt64("userId")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, Result[[]ArticleVO]{Code: 401, Msg: "未登录"})
		return
	}
	list, err := h.svc.List(c, uid, req.Offset, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result[[]ArticleVO]{Code: 500, Msg: "系统错误"})
		return
	}
	result := make([]ArticleVO, 0, len(list))
	for _, a := range list {
		result = append(result, ArticleVO{
			ID:        a.ID,
			Title:     a.Title,
			Abstract:  a.GenAbstract(),
			Author:    a.Author.Name,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			Status:    a.Status,
		})
	}
	c.JSON(http.StatusOK, Result[[]ArticleVO]{Code: 0, Msg: "获取文章列表成功", Data: result})
}

// PubList 返回所有已发布文章的列表，免登录
func (h *ArticleHandler) PubList(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Result[[]ArticleVO]{Code: 400, Msg: "参数错误"})
		return
	}
	if req.Limit == 0 {
		req.Limit = 20
	}
	list, err := h.svc.ListPub(c, req.Offset, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result[[]ArticleVO]{Code: 500, Msg: "系统错误"})
		return
	}
	result := make([]ArticleVO, 0, len(list))
	for _, a := range list {
		result = append(result, toVO(a))
	}
	c.JSON(http.StatusOK, Result[[]ArticleVO]{Code: 0, Msg: "获取公开文章列表成功", Data: result})
}

func (h *ArticleHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result[ArticleVO]{Code: 400, Msg: "参数错误"})
		return
	}
	uid := c.GetInt64("userId")
	article, err := h.svc.Detail(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result[ArticleVO]{Code: 500, Msg: "系统错误"})
		return
	}
	if uid != 0 && article.Author.ID != uid {
		c.JSON(http.StatusForbidden, Result[ArticleVO]{Code: 403, Msg: "无权限查看"})
		return
	}
	c.JSON(http.StatusOK, Result[ArticleVO]{Code: 0, Msg: "获取文章详情成功", Data: toVO(article)})
}

func (h *ArticleHandler) PubDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result[ArticleVO]{Code: 400, Msg: "参数错误"})
		return
	}
	uid := c.GetInt64("userId")
	article, err := h.svc.PubDetail(c, id, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, Result[ArticleVO]{Code: 404, Msg: "文章不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, Result[ArticleVO]{Code: 500, Msg: "系统错误"})
		return
	}
	info, _ := h.interactive.Get(c, "article", id, uid)
	_ = h.interactive.IncrRead(c, "article", id)
	vo := toVO(article)
	vo.ReadCnt = info.ReadCnt + 1
	vo.LikeCnt = info.LikeCnt
	vo.CollectCnt = info.CollectCnt
	vo.Liked = info.Liked
	vo.Collected = info.Collected
	c.JSON(http.StatusOK, Result[ArticleVO]{Code: 0, Msg: "获取文章详情成功", Data: vo})
}

// PubLike 简化版：仅回传成功，不做真实计数
func (h *ArticleHandler) PubLike(c *gin.Context) {
	var req LikeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "参数错误"})
		return
	}
	uid := c.GetInt64("userId")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, Result[string]{Code: 401, Msg: "未登录"})
		return
	}
	info, err := h.interactive.Like(c, "article", req.ID, uid, req.Like)
	if err != nil {
		h.l.Error("点赞失败", logger.Error(err))
		c.JSON(http.StatusInternalServerError, Result[string]{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result[domain.Interactive]{Code: 0, Msg: "点赞成功", Data: info})
}

type CollectReq struct {
	ID  int64 `json:"id"`
	CID int64 `json:"cid"`
}

func (h *ArticleHandler) PubCollect(c *gin.Context) {
	var req CollectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "参数错误"})
		return
	}
	uid := c.GetInt64("userId")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, Result[string]{Code: 401, Msg: "未登录"})
		return
	}
	info, err := h.interactive.Collect(c, "article", req.ID, uid)
	if err != nil {
		h.l.Error("收藏失败", logger.Error(err))
		c.JSON(http.StatusInternalServerError, Result[string]{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result[domain.Interactive]{Code: 0, Msg: "收藏成功", Data: info})
}

type RewardReq struct {
	ID  int64 `json:"id"`
	Amt int64 `json:"amt"`
}
type RewardResp struct {
	CodeURL string `json:"codeURL"`
	Rid     int64  `json:"rid"`
}

// PubReward 简化版：返回假二维码链接和 rid
func (h *ArticleHandler) PubReward(c *gin.Context) {
	var req RewardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result[string]{Code: 400, Msg: "参数错误"})
		return
	}
	c.JSON(http.StatusOK, Result[RewardResp]{Code: 0, Msg: "创建支付码成功", Data: RewardResp{
		CodeURL: "https://example.com/reward/" + strconv.FormatInt(req.ID, 10),
		Rid:     req.ID,
	}})
}

func (h *ArticleHandler) RewardDetail(c *gin.Context) {
	c.JSON(http.StatusOK, Result[string]{Code: 0, Msg: "查询成功", Data: "RewardStatusPayed"})
}

func toVO(a domain.Article) ArticleVO {
	return ArticleVO{
		ID:        a.ID,
		Title:     a.Title,
		Abstract:  a.GenAbstract(),
		Content:   a.Content,
		Author:    a.Author.Name,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Status:    a.Status,
	}
}
