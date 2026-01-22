package crcwatch

import (
	"time"

	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
	watchor "github.com/smartxworks/cloudtower-go-sdk/v2/watchor"
	"k8s.io/klog"
)

type CrcEventHandlerFunc func(*models.ResourceChangeEvent)

type Watch struct {
	opts           *Options
	resourceTypes  []string
	crcWatchClient ResourceChangeWatcher
	eventHandler   CrcEventHandlerFunc
}

func NewWatch(resourcesTypes []string, opts ...OptionFunc) (*Watch, error) {
	w := &Watch{
		resourceTypes: resourcesTypes,
		opts:          &Options{PollingInterval: 10 * time.Second},
	}
	for _, f := range opts {
		f(w.opts)
	}
	var err error
	w.crcWatchClient, err = NewWatchClient(w.resourceTypes, w.opts)
	if w.crcWatchClient == nil || err != nil {
		klog.Errorf("Failed to init crc watch client: %s", err)
		return nil, err
	}
	return w, nil
}

func (w *Watch) Start(stopCh <-chan struct{}) {
	if w.eventHandler == nil {
		klog.Fatal("must registry crc event handler")
	}
	crcwLoop := func() {
		err := w.crcWatchClient.Start(&watchor.ResourceChangeWatchStartParams{
			StartRevision: nil,
		})

		if err != nil {
			klog.Fatalln(err)
		}

		klog.Infoln("crc watch start")

		for {
			select {
			case err := <-w.crcWatchClient.ErrorChannel():
				if err == nil {
					klog.Fatal("crc watch get nil error")
				}
				if err.CompactRevision != nil {
					klog.Fatalf("crc event missed, compacted error : %v, compact revision %s", err, *err.CompactRevision)
				} else if err.Err != nil {
					klog.Errorf("crc error event: %v", err)
					if err.Type == watchor.ErrorEventTypeUnsupported {
						// after unsupported error, crc will stop event loop
						return
					}
				}
			case warning := <-w.crcWatchClient.WarningChannel():
				if warning.Err != nil {
					klog.Warningf("crc warning event %s", warning.Err.Error())
				}
			case event := <-w.crcWatchClient.Channel():
				w.eventHandler(event)
			case <-stopCh:
				return
			}
		}
	}

	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				crcwLoop()
				// crcwLoop will return when crc not supported
				// restart after 10 minutes for tower upgrade
				time.Sleep(10 * time.Minute)
			}
		}
	}()
}

func (w *Watch) RegistryHandler(f CrcEventHandlerFunc) {
	w.eventHandler = f
}
