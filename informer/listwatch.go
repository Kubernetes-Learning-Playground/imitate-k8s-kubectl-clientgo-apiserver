package informer

// Lister is any object that knows how to perform an initial list.
type Lister interface {

	List() ([]interface{}, error)
}

// Watcher is any object that knows how to start a watch on a resource.
type Watcher interface {

	Watch() (interface{}, error)
}

// ListerWatcher is any object that knows how to perform an initial list and start a watch on a resource.
type ListerWatcher interface {
	Lister
	Watcher
}

// ListFunc knows how to list resources
type ListFunc func() (interface{}, error)

// WatchFunc knows how to watch resources
type WatchFunc func() (interface{}, error)

// ListWatch knows how to list and watch a set of apiserver resources.  It satisfies the ListerWatcher interface.
// It is a convenience function for users of NewReflector, etc.
// ListFunc and WatchFunc must not be nil
type ListWatch struct {
	ListFunc  ListFunc
	WatchFunc WatchFunc

}
