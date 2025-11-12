package lock

import (
	"fmt"
	"sync"
	"time"

	"github.com/aid297/aid/dict"
)

type (
	MapLocker interface {
		implMapLocker()
		New() MapLocker
		Once() MapLocker
		Set(key string, val any) error
		SetMany(items map[string]any) error
		Destroy(key string)
		DestroyAll()
		Lock(key string, timeout time.Duration) (*itemLock, error)
		Try(key string) error
	}

	// MapLock 字典锁：一个锁的集合
	MapLock struct {
		lock  sync.RWMutex
		locks *dict.AnyDict[string, *itemLock]
	}

	// 锁项：一个集合锁中的每一项，包含：锁状态、锁值、超时时间、定时器
	itemLock struct {
		inUse   bool
		val     any
		timeout time.Duration
		timer   *time.Timer
	}
)

var (
	onceMapLock sync.Once
	mapLockIns  *MapLock
)

func (*MapLock) implMapLocker() {}

func (*MapLock) New() MapLocker { return &MapLock{locks: dict.Make[string, *itemLock]()} }

func (*MapLock) Once() MapLocker {
	onceMapLock.Do(func() { mapLockIns = &MapLock{locks: dict.Make[string, *itemLock]()} })

	return mapLockIns
}

func (*MapLock) set(key string, val any) (err error) {
	_, exists := mapLockIns.locks.Get(key)
	if exists {
		return fmt.Errorf("锁[%s]已存在", key)
	} else {
		mapLockIns.locks.Set(key, &itemLock{val: val})
	}

	return
}

// Set 创建锁
func (*MapLock) Set(key string, val any) error {
	mapLockIns.lock.Lock()
	defer mapLockIns.lock.Unlock()

	return mapLockIns.set(key, val)
}

// SetMany 批量创建锁
func (*MapLock) SetMany(items map[string]any) (err error) {
	mapLockIns.lock.Lock()
	defer mapLockIns.lock.Unlock()

	for idx, item := range items {
		if err = mapLockIns.set(idx, item); err != nil {
			mapLockIns.DestroyAll()
			return
		}
	}

	return
}

// Release 显式锁释放方法
func (r *itemLock) Release() {
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.inUse = false
}

func (*MapLock) destroy(key string) {
	if il, ok := mapLockIns.locks.Get(key); ok {
		il.Release()
		mapLockIns.locks.RemoveByKey(key) // 删除键值对，以便垃圾回收
	}
}

// Destroy 删除锁
func (*MapLock) Destroy(key string) {
	mapLockIns.lock.Lock()
	defer mapLockIns.lock.Unlock()

	mapLockIns.destroy(key)
}

// DestroyAll 删除所有锁
func (*MapLock) DestroyAll() {
	mapLockIns.lock.Lock()
	defer mapLockIns.lock.Unlock()

	mapLockIns.locks.Each(func(key string, value *itemLock) {
		mapLockIns.destroy(key)
	})
}

// Lock 加锁
func (*MapLock) Lock(key string, timeout time.Duration) (*itemLock, error) {
	mapLockIns.lock.RLock()
	defer mapLockIns.lock.RUnlock()

	if item, exists := mapLockIns.locks.Get(key); !exists {
		return nil, fmt.Errorf("锁[%s]不存在", key)
	} else {
		if item.inUse {
			return nil, fmt.Errorf("锁[%s]被占用", key)
		}

		// 设置锁占用
		item.inUse = true

		// 设置超时时间
		if timeout > 0 {
			item.timeout = timeout
			item.timer = time.AfterFunc(timeout, func() {
				if il, ok := mapLockIns.locks.Get(key); ok {
					if il.timer != nil {
						il.Release()
					}
				}
			})
		}

		return item, nil
	}
}

// Try 尝试获取锁
func (*MapLock) Try(key string) error {
	mapLockIns.lock.RLock()
	defer mapLockIns.lock.RUnlock()

	if item, exist := mapLockIns.locks.Get(key); !exist {
		return fmt.Errorf("锁[%s]不存在", key)
	} else {
		if item.inUse {
			return fmt.Errorf("锁[%s]被占用", key)
		}
		return nil
	}
}
