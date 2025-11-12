package lock

import (
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	k8sLinks := map[string]any{
		"k8s-a": &struct{}{},
		"k8s-b": &struct{}{},
		"k8s-c": &struct{}{},
	}

	// 获取字典锁对象
	ml := APP.MapLock.Once()

	// 批量创建锁
	storeErr := ml.SetMany(k8sLinks)
	if storeErr != nil {
		// 处理err
		t.Fatal(storeErr.Error())
	}

	// 检测锁
	tryErr := ml.Try("k8s-a")
	if tryErr != nil {
		// 处理err
		t.Fatal(tryErr.Error())
	}

	// 获取锁
	lock, lockErr := ml.Lock("k8s-a", time.Second*10) // 10秒业务处理不完也会过期 设置为：0则为永不过期
	if lockErr != nil {
		t.Fatal(lockErr.Error())
	}
	defer lock.Release()

	// 处理业务...
}
