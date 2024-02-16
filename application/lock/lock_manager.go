package lock

type LockManager interface {
	Acquire(key string) error
	Release(key string) error
}
