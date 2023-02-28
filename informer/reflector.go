package informer

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/pager"
	"k8s.io/klog/v2"
	"k8s.io/utils/clock"
	"k8s.io/utils/trace"
	"math/rand"
)

type Reflector struct {
	// 类型
	object  interface{}
	// delta fifo
	store   Store
	// 特定类型的 list watch方法
	listerWatcher ListerWatcher
}

func NewReflector(object interface{}, store Store, listerWatcher ListerWatcher) *Reflector {
	return &Reflector{object: object, store: store, listerWatcher: listerWatcher}
}


func (r *Reflector) Run(stopCh <-chan struct{}) {

	go func() {
		if err := r.ListAndWatch(stopCh); err != nil {
			return
		}
	}()

}

func (r *Reflector) ListAndWatch(stopCh <-chan struct{}) error {


	// 1. 先使用list方法，拉取指定资源的list。
	err := r.list()
	if err != nil {
		return err
	}


	// 执行watch操作
	for {
		// give the stopCh a chance to stop the loop, even in case of continue statements further down on errors
		select {
		case <-stopCh:
			return nil
		default:
		}


		w, err := r.listerWatcher.Watch()
		if err != nil {
			return err
		}

		// 每当watch到时，需要区分不同的"事件"，再操作store(delta fifo)
		err = watchHandler(start, w, r.store, r.expectedType, r.expectedGVK, r.name, r.typeDescription, r.setLastSyncResourceVersion, r.clock, resyncerrc, stopCh)

		if err != nil {
			return err
		}
	}
}


func (r *Reflector) list() error {

	items, err := r.listerWatcher.List()
	if err != nil {
		return err
	}

	// 重要。 把list的[]runtime.Object结果放入store(delta fifo中)。使用接口的Replace方法
	if err := r.syncWith(items); err != nil {
		return fmt.Errorf("unable to sync list result: %v", err)
	}

	return nil
}

func (r *Reflector) syncWith(items []interface{}) error {
	found := make([]interface{}, 0, len(items))
	for _, item := range items {
		found = append(found, item)
	}
	return r.store.Replace(found)
}


func watchHandler(
	store Store,
	objectType interface{},
	errc chan error,
	stopCh <-chan struct{},
) error {


loop:
	for {
		select {
		case <-stopCh:
			return nil
		case err := <-errc:
			return err
		// 从 ResultChan取出对象
		case event, ok := <-w.ResultChan():
			if !ok {
				break loop
			}
			if event.Type == watch.Error {
				return apierrors.FromObject(event.Object)
			}
			if expectedType != nil {
				if e, a := expectedType, reflect.TypeOf(event.Object); e != a {
					utilruntime.HandleError(fmt.Errorf("%s: expected type %v, but watch event object had type %v", name, e, a))
					continue
				}
			}
			if expectedGVK != nil {
				if e, a := *expectedGVK, event.Object.GetObjectKind().GroupVersionKind(); e != a {
					utilruntime.HandleError(fmt.Errorf("%s: expected gvk %v, but watch event object had gvk %v", name, e, a))
					continue
				}
			}
			// 拿到object对象
			meta, err := meta.Accessor(event.Object)
			if err != nil {
				utilruntime.HandleError(fmt.Errorf("%s: unable to understand watch event %#v", name, event))
				continue
			}


			// 重要，区分不同事件
			switch event.Type {
			// add事件
			case watch.Added:
				err := store.Add(event.Object)
				if err != nil {
					return err
				}
			// update事件
			case watch.Modified:
				err := store.Update(event.Object)
				if err != nil {
					return err
				}
			// delete事件
			case watch.Deleted:
				// TODO: Will any consumers need access to the "last known
				// state", which is passed in event.Object? If so, may need
				// to change this.
				err := store.Delete(event.Object)
				if err != nil {
					return err
				}
			case watch.Bookmark:
				// A `Bookmark` means watch has synced here, just update the resourceVersion
			default:
			}

		}
	}

	return nil
}
