/*
Copyright 2021 The Everoute Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package informer

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/everoute/graphc/pkg/client"
	"github.com/everoute/graphc/pkg/schema"
	"github.com/everoute/graphc/third_party/forked/client-go/informer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

// NewSharedInformerFactory constructs a new instance of sharedInformerFactory for all resources
func NewSharedInformerFactory(client *client.Client, defaultResync time.Duration) *SharedInformerFactory {
	factory := &SharedInformerFactory{
		client:           client,
		defaultResync:    defaultResync,
		informers:        make(map[reflect.Type]cache.SharedIndexInformer),
		startedInformers: make(map[reflect.Type]bool),
		customResync:     make(map[reflect.Type]time.Duration),
	}
	return factory
}

type SharedInformerFactory struct {
	client        *client.Client
	lock          sync.Mutex
	defaultResync time.Duration
	customResync  map[reflect.Type]time.Duration

	informers map[reflect.Type]cache.SharedIndexInformer
	// startedInformers is used for tracking which informers have been started.
	// This allows Start() to be called multiple times safely.
	startedInformers map[reflect.Type]bool
}

// Start implements SharedInformerFactory.Start
func (f *SharedInformerFactory) Start(stopCh <-chan struct{}) {
	f.lock.Lock()
	defer f.lock.Unlock()

	for informerType, sharedInformer := range f.informers {
		if !f.startedInformers[informerType] {
			go sharedInformer.Run(stopCh)
			f.startedInformers[informerType] = true
		}
	}
}

// WaitForCacheSync implements SharedInformerFactory.WaitForCacheSync
func (f *SharedInformerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool {
	informers := func() map[reflect.Type]cache.SharedIndexInformer {
		f.lock.Lock()
		defer f.lock.Unlock()

		informers := map[reflect.Type]cache.SharedIndexInformer{}
		for informerType, sharedInformer := range f.informers {
			if f.startedInformers[informerType] {
				informers[informerType] = sharedInformer
			}
		}
		return informers
	}()

	res := map[reflect.Type]bool{}
	for informType, sharedInformer := range informers {
		res[informType] = cache.WaitForCacheSync(stopCh, sharedInformer.HasSynced)
	}
	return res
}

// InformerFor implements SharedInformerFactory.InformerFor
func (f *SharedInformerFactory) InformerFor(obj schema.Object) cache.SharedIndexInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	informerType := reflect.TypeOf(obj)
	sharedInformer, exists := f.informers[informerType]
	if exists {
		return sharedInformer
	}

	resyncPeriod, exists := f.customResync[informerType]
	if !exists {
		resyncPeriod = f.defaultResync
	}

	sharedInformer = defaultNewInformerFunc(f.client, obj, resyncPeriod)
	f.informers[informerType] = sharedInformer

	return sharedInformer
}

func defaultNewInformerFunc(c *client.Client, obj schema.Object, resyncPeriod time.Duration) cache.SharedIndexInformer {
	var newReflectorFunc = NewReflectorBuilder(c)
	return informer.NewSharedIndexInformer(newReflectorFunc, obj, DefaultKeyFunc, resyncPeriod, cache.Indexers{})
}

func ReconcileWorker(name string, queue workqueue.RateLimitingInterface, processFunc func(string) error) func() {
	return func() {
		for {
			key, quit := queue.Get()
			if quit {
				return
			}

			err := processFunc(key.(string))
			if err != nil {
				queue.Done(key)
				queue.AddRateLimited(key)
				klog.Errorf("%s got error while sync %s: %s", name, key.(string), err)
				continue
			}

			// stop the rate limiter from tracking the key
			queue.Done(key)
			queue.Forget(key)
		}
	}
}

func DefaultKeyFunc(obj interface{}) (string, error) {
	if d, ok := obj.(cache.DeletedFinalStateUnknown); ok {
		return d.Key, nil
	}
	resource, ok := obj.(schema.Object)
	if ok && !reflect.ValueOf(resource).IsNil() {
		return resource.GetID(), nil
	}
	return "", fmt.Errorf("unsupport resource type %s, object: %v", obj, obj)
}
