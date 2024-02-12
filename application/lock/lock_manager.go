package lock

type LockManager interface {
	Acquire(key int)
}
