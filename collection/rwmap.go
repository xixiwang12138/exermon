package collection

import "sync"

// RWMap 使用读写锁保护的并发安全map
type RWMap[K string | int | int64 | uintptr, V any] struct {
	data map[K]V
	lock *sync.RWMutex
}

func NewRWMap[K string | int | int64, V any]() *RWMap[K, V] {
	return &RWMap[K, V]{
		data: make(map[K]V, 0),
		lock: &sync.RWMutex{},
	}
}

func (m *RWMap[K, V]) Get(k K) (V, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	v, existed := m.data[k]
	return v, existed
}

func (m *RWMap[K, V]) GetOrSet(k K, f func() V) V {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, existed := m.data[k]
	if existed {
		return v
	}
	newV := f()
	m.data[k] = newV
	return newV
}

func (m *RWMap[K, V]) Set(k K, v V) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[k] = v
}

func (m *RWMap[K, V]) Delete(k K) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.data, k)
}

func (m *RWMap[K, V]) Len() int { // map的长度
	m.lock.RLock() // 锁保护
	defer m.lock.RUnlock()
	return len(m.data)
}

func (m *RWMap[K, V]) Each(f func(k K, v V) bool) { // 遍历map
	m.lock.RLock()
	defer m.lock.RUnlock()

	for k, v := range m.data {
		if !f(k, v) {
			return
		}
	}
}
