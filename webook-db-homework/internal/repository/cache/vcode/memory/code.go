package memory

import (
	"context"
	"fmt"
	"sync"
	"time"
	"webook/internal/repository/cache/vcode"
)

type CodeInfo struct {
	Code  string
	Cnt   int64
	Utime int64
}
type MemoryCodeCache struct {
	mp    map[string]CodeInfo
	cache vcode.CodeCache
	mu    sync.RWMutex
}

func NewMemoryCodeCache(mp map[string]CodeInfo, cache vcode.CodeCache) vcode.CodeCache {
	return &MemoryCodeCache{
		mp:    mp,
		cache: cache,
	}
}

func (c *MemoryCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	key := c.key(biz, phone)
	var res int
	c.mu.RLock()
	info, ok := c.mp[key]
	if ok {
		_, ttl := c.parse(info)
		if ttl < 60 {
			res = -1
		} else if ttl < 0 {
			res = -2
		}
	}
	c.mu.RUnlock()
	if res == -1 {
		return vcode.ErrSetCodeBusy
	} else if res == -2 {
		return vcode.ErrSetCodeSystemError
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	newinfo, okNew := c.mp[key]
	if newinfo != info || (!ok && okNew) {
		return vcode.ErrSetCodeBusy
	}
	if ok && !okNew {
		//可能是校验后删掉了
		return vcode.ErrSetCodeSystemError
	}
	if okNew {
		_, ttl := c.parse(newinfo)
		if ttl > 60 {
			c.mp[key] = CodeInfo{
				Code:  code,
				Cnt:   3,
				Utime: time.Now().UnixMilli(),
			}
		} else if ttl < 0 {
			return vcode.ErrSetCodeSystemError
		} else {
			return vcode.ErrSetCodeBusy
		}
	} else {
		c.mp[key] = CodeInfo{
			Code:  code,
			Cnt:   3,
			Utime: time.Now().UnixMilli(),
		}
	}
	return nil
}

func (c *MemoryCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *MemoryCodeCache) parse(info CodeInfo) (int64, int64) {
	return info.Cnt, int64(time.Since(time.UnixMilli(info.Utime)).Seconds())
}

func (c *MemoryCodeCache) Verify(ctx context.Context, biz string, phone string, code string) error {
	key := c.key(biz, phone)

	c.mu.Lock()
	defer c.mu.Unlock()

	info, ok := c.mp[key]
	if !ok {
		return vcode.ErrVarifyCodeInvalid
	}

	cnt, ttl := c.parse(info)
	if ttl < 0 || ttl > 600 {
		return vcode.ErrVarifyCodeInvalid
	}

	if cnt <= 0 {
		return vcode.ErrVarifyCodeTooMany
	}

	if info.Code != code {
		// 验证码错误，减少剩余次数
		c.mp[key] = CodeInfo{
			Code:  info.Code,
			Cnt:   cnt - 1,
			Utime: info.Utime,
		}
		return vcode.ErrVarifyCodeInvalid
	}

	// 验证成功，删除验证码
	delete(c.mp, key)
	return nil
}
