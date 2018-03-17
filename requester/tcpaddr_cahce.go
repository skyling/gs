package requester

import (
	"net"
	"sync"
	"time"
)

var (
	// TCPAddrCache dns 缓存
	TCPAddrCache = tcpAddrCache{
		ta:       sync.Map{},
		lifeTime: 1 * time.Minute,
	}
)

// tcpAddrCache tcp 地址缓存,即dns解析后的ip地址
type tcpAddrCache struct {
	ta        sync.Map
	lifeTime  time.Duration // 生命周期
	gcStarted bool
}

// Set 设置
func (tac *tcpAddrCache) Set(address string, ta *net.TCPAddr) {
	tac.ta.Store(address, ta)
}

// Existed 检测是否存在
func (tac *tcpAddrCache) Existed(address string) bool {
	v, existed := tac.ta.Load(address)
	if existed && v == nil {
		return false
	}
	return existed
}

// Get 获取值
func (tac *tcpAddrCache) Get(address string) *net.TCPAddr {
	if tac.Existed(address) {
		value, _ := tac.ta.Load(address)
		return value.(*net.TCPAddr)
	}
	return nil
}

// SetLifeTime 设置生命周期
func (tac *tcpAddrCache) SetLifeTime(t time.Duration) {
	tac.lifeTime = t
}

// GC 缓存回收
func (tac *tcpAddrCache) GC() {
	if tac.gcStarted {
		return
	}
	go func() {
		for {
			time.Sleep(tac.lifeTime) // 这样可以动态修改lifetime
			tac.DelAll()
		}
	}()
}

// Del 删除
func (tac *tcpAddrCache) Del(address string) {
	tac.ta.Delete(address)
}

// DelAll 清空缓存
func (tac *tcpAddrCache) DelAll() {
	tac.ta.Range(func(address, _ interface{}) bool {
		tac.ta.Delete(address)
		return true
	})
}
