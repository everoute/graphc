package crcwatch

import (
	"errors"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
	"github.com/smartxworks/cloudtower-go-sdk/v2/watchor"

	"github.com/everoute/graphc/pkg/client"
)

func TestCrcwatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crcwatch Suite")
}

type fakeWatchClient struct {
	lock      sync.Mutex
	eventCh   chan *models.ResourceChangeEvent
	errCh     chan *watchor.ErrorEvent
	warnCh    chan *watchor.WarningEvent
	startErr  error
	startCall int
}

func newFakeWatchClient() *fakeWatchClient {
	return &fakeWatchClient{
		eventCh: make(chan *models.ResourceChangeEvent, 1),
		errCh:   make(chan *watchor.ErrorEvent, 1),
		warnCh:  make(chan *watchor.WarningEvent, 1),
	}
}

func (f *fakeWatchClient) Start(_ *watchor.ResourceChangeWatchStartParams) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.startCall++
	return f.startErr
}

func (f *fakeWatchClient) getStartCall() int {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.startCall
}
func (f *fakeWatchClient) Channel() <-chan *models.ResourceChangeEvent  { return f.eventCh }
func (f *fakeWatchClient) ErrorChannel() <-chan *watchor.ErrorEvent     { return f.errCh }
func (f *fakeWatchClient) WarningChannel() <-chan *watchor.WarningEvent { return f.warnCh }

var _ = Describe("Watch", func() {
	var (
		w       *Watch
		fakeCli *fakeWatchClient
		stopCh  chan struct{}
	)

	BeforeEach(func() {
		stopCh = make(chan struct{})
		fakeCli = newFakeWatchClient()

		w = &Watch{
			opts:           &Options{},
			resourceTypes:  []string{"vm"},
			crcWatchClient: fakeCli, // directly inject fake
		}
	})

	AfterEach(func() {
		close(stopCh)
	})

	It("should set options via OptionFunc", func() {
		userInfo := &client.UserInfo{Username: "u", Password: "p", Source: "s"}
		opts := &Options{}
		SetUserInfo(userInfo)(opts)
		Expect(opts.UserInfo).To(Equal(userInfo))

		SetAPIAuth("user", "pass")(opts)
		Expect(opts.APIUsername).To(Equal("user"))
		Expect(opts.APIPassword).To(Equal("pass"))

		SetHost("127.0.0.1")(opts)
		Expect(opts.Host).To(Equal("127.0.0.1"))
	})

	It("should call handler when event received", func() {
		called := make(chan bool, 1)
		w.RegistryHandler(func(e *models.ResourceChangeEvent) {
			called <- true
		})

		go w.Start(stopCh)

		resType := "vm"
		fakeCli.eventCh <- &models.ResourceChangeEvent{ResourceType: &resType}
		Eventually(called, 2*time.Second).Should(Receive(BeTrue()))
	})

	It("should handle error event for unsupport event", func() {
		w.RegistryHandler(func(e *models.ResourceChangeEvent) {})
		go w.Start(stopCh)

		fakeCli.errCh <- &watchor.ErrorEvent{Err: errors.New("test-error"), Type: watchor.ErrorEventTypeUnsupported}
		Eventually(func() int { return fakeCli.getStartCall() }, 2*time.Second).Should(Equal(1))
	})

	It("should handle warning event", func() {
		w.RegistryHandler(func(e *models.ResourceChangeEvent) {})
		go w.Start(stopCh)

		fakeCli.warnCh <- &watchor.WarningEvent{Err: errors.New("warn")}
		Eventually(func() int { return fakeCli.getStartCall() }, 2*time.Second).Should(BeNumerically("==", 1))
	})
})

var _ = Describe("NewWatch", func() {
	It("should return nil when NewWatchClient returns nil", func() {
		w, err := NewWatch([]string{"TestResource"}, SetUserInfo(&client.UserInfo{}))
		Expect(w).To(BeNil())
		Expect(err).ToNot(BeNil())
	})
})
