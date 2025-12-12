package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Interactive struct {
	ID         int64  `gorm:"primaryKey,autoIncrement;column:id"`
	BizId      int64  `gorm:"uniqueIndex:biz_id_type,column:biz_id"`
	Biz        string `gorm:"uniqueIndex:biz_id_type;type:varchar(128);column:biz"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
	Readcnt    int64  `gorm:"column:readcnt"`
	Likecnt    int64  `gorm:"column:likecnt"`
	Collectcnt int64  `gorm:"column:collectcnt"`
}

type UserLikeSomething struct {
	ID        int64  `gorm:"primaryKey,autoIncrement;column:id"`
	BizId     int64  `gorm:"uniqueIndex:biz_id_type,column:biz_id"`
	Biz       string `gorm:"uniqueIndex:biz_id_type;type:varchar(128);column:biz"`
	UID       int64  `gorm:"uniqueIndex:biz_id_type;column:uid"`
	Status    bool   `gorm:"column:status"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

// 收藏夹业务
type Collection struct {
	ID        int64  `gorm:"primaryKey,autoIncrement;column:id"`
	Name      string `gorm:"column:name"`
	UID       int64  `gorm:"column:uid"`
	Status    bool   `gorm:"column:status"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

type UserCollectSomething struct {
	ID        int64  `gorm:"primaryKey,autoIncrement;column:id"`
	BizId     int64  `gorm:"uniqueIndex:biz_id_type,column:biz_id"`
	Biz       string `gorm:"uniqueIndex:biz_id_type;type:varchar(128);column:biz"`
	UID       int64  `gorm:"uniqueIndex:biz_id_type;column:uid"`
	CollectId int64  `gorm:"index;column:collect_id"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

/*
// 这里考虑收藏夹的业务怎么实现，需要用到join查询 也就是通过收藏夹的name或者id查询到收藏夹，然后通过id作为外键到存储usercollectsomething表中查询具体收藏实体的信息
func (dao *GORMInteractiveDAO) GetCollectionContentsByName(ctx context.Context, collectName string, uid int64) ([]UserCollectSomething, error) {
    var contents []UserCollectSomething

    err := dao.db.WithContext(ctx).
        Table("user_collect_somethings ucs").
        Joins("JOIN collections c ON ucs.collect_id = c.id").
        Where("c.name = ? AND c.uid = ? AND c.status = ?", collectName, uid, true).
        Select("ucs.*").
        Order("ucs.created_at DESC").
        Find(&contents).Error

    return contents, err
}

// 根据收藏夹ID查询收藏内容（也可以JOIN，获取更多信息）
func (dao *GORMInteractiveDAO) GetCollectionContentsWithInfo(ctx context.Context, collectId int64, uid int64) ([]CollectionContentVO, error) {
    var contents []CollectionContentVO

    err := dao.db.WithContext(ctx).
        Table("user_collect_somethings ucs").
        Joins("JOIN collections c ON ucs.collect_id = c.id").
        Where("ucs.collect_id = ? AND ucs.uid = ? AND c.status = ?", collectId, uid, true).
        Select("ucs.*, c.name as collection_name, c.created_at as collection_created_at").
        Order("ucs.created_at DESC").
        Find(&contents).Error

    return contents, err
}
*/

type InteractiveDAO interface {
	IncLike(ctx context.Context, biz string, id int64, uid int64) error
	DecLike(ctx context.Context, biz string, id int64, uid int64) error
	IncRead(ctx context.Context, biz string, id int64) error
	IncCollect(ctx context.Context, biz string, id int64, uid int64) error
	//DecCollect(ctx context.Context, biz string, id int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (Interactive, error)                                //获取互动信息
	GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeSomething, error)       //获取点赞信息
	GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectSomething, error) //获取收藏信息
	BatchIncRead(ctx context.Context, bizs []string, ids []int64) error
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]Interactive, error)
}

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func NewInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{db: db}
}

func (dao *GORMInteractiveDAO) GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]Interactive, error) {
	var interactives []Interactive
	for _, id := range ids {
		var interactive Interactive
		err := dao.db.WithContext(ctx).Model(&Interactive{}).Where("biz = ? AND biz_id = ?", biz, id).First(&interactive).Error
		if err != nil {
			return nil, err
		}
		interactives = append(interactives, Interactive{
			BizId:      id,
			Biz:        biz,
			Readcnt:    interactive.Readcnt,
			Likecnt:    interactive.Likecnt,
			Collectcnt: interactive.Collectcnt,
			CreatedAt:  interactive.CreatedAt,
			UpdatedAt:  interactive.UpdatedAt,
		})
	}

	res := make(map[int64]Interactive)
	for _, interactive := range interactives {
		res[interactive.BizId] = interactive
	}
	return res, nil
}
func (dao *GORMInteractiveDAO) BatchIncRead(ctx context.Context, bizs []string, ids []int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := &GORMInteractiveDAO{db: tx}
		for i, biz := range bizs {
			err := txDAO.IncRead(ctx, biz, ids[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (dao *GORMInteractiveDAO) IncLike(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ul UserLikeSomething
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).
			First(&ul).Error
		switch err {
		case nil:
			if ul.Status {
				// 已点赞，无需重复计数
				return nil
			}
			if err = tx.Model(&UserLikeSomething{}).Where("id = ?", ul.ID).
				Updates(map[string]interface{}{
					"status":     true,
					"updated_at": now,
				}).Error; err != nil {
				return err
			}
		case gorm.ErrRecordNotFound:
			if err = tx.Create(&UserLikeSomething{
				BizId:     id,
				Biz:       biz,
				UID:       uid,
				Status:    true,
				CreatedAt: now,
				UpdatedAt: now,
			}).Error; err != nil {
				return err
			}
		default:
			return err
		}

		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "biz_id"}, {Name: "biz"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"updated_at": now, "likecnt": gorm.Expr("likecnt + 1")}),
		}).Create(&Interactive{BizId: id, Biz: biz, CreatedAt: now, UpdatedAt: now, Likecnt: 1}).Error
	})
}

func (dao *GORMInteractiveDAO) DecLike(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ul UserLikeSomething
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).
			First(&ul).Error
		switch err {
		case nil:
			if !ul.Status {
				// 未点赞，无需扣减
				return nil
			}
			if err = tx.Model(&UserLikeSomething{}).Where("id = ?", ul.ID).
				Updates(map[string]interface{}{
					"status":     false,
					"updated_at": now,
				}).Error; err != nil {
				return err
			}
		case gorm.ErrRecordNotFound:
			// 没有记录，不扣减
			return nil
		default:
			return err
		}
		return tx.Model(&Interactive{}).Where("biz = ? AND biz_id = ?", biz, id).Updates(map[string]interface{}{
			"updated_at": now,
			"likecnt":    gorm.Expr("likecnt - 1"),
		}).Error
	})
}

func (dao *GORMInteractiveDAO) IncRead(ctx context.Context, biz string, id int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		//Columns: []clause.Column{{Name: "biz_id"}},  这一行mysql不用写
		DoUpdates: clause.Assignments(map[string]interface{}{
			"updated_at": now,
			"readcnt":    gorm.Expr("readcnt + 1"),
		}),
	}).Create(&Interactive{BizId: id, Biz: biz, CreatedAt: now, UpdatedAt: now, Readcnt: 1}).Error
}

func (dao *GORMInteractiveDAO) IncCollect(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&UserCollectSomething{BizId: id, Biz: biz, UID: uid, CreatedAt: now, UpdatedAt: now}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": now,
				"collectcnt": gorm.Expr("collectcnt + 1"),
			}),
		}).Create(&Interactive{BizId: id, Biz: biz, CreatedAt: now, UpdatedAt: now, Collectcnt: 1}).Error
	})
}

func (dao *GORMInteractiveDAO) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	var Interactive Interactive
	err := dao.db.WithContext(ctx).Model(&Interactive).Where("biz = ? AND biz_id = ?", biz, id).First(&Interactive).Error
	return Interactive, err
}

func (dao *GORMInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeSomething, error) {
	var data UserLikeSomething
	err := dao.db.WithContext(ctx).Model(&UserLikeSomething{}).Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).First(&data).Error
	return data, err
}

func (dao *GORMInteractiveDAO) GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectSomething, error) {
	var userCollectSomething UserCollectSomething
	err := dao.db.WithContext(ctx).Model(&UserCollectSomething{}).Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).First(&userCollectSomething).Error
	return userCollectSomething, err
}
