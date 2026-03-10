package crcwatch

import (
	"errors"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	apiclient "github.com/smartxworks/cloudtower-go-sdk/v2/client"
	resource_change_client "github.com/smartxworks/cloudtower-go-sdk/v2/client/resource_change"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
	watchor "github.com/smartxworks/cloudtower-go-sdk/v2/watchor"
	"k8s.io/klog"

	"github.com/everoute/graphc/pkg/client"
)

var _ ResourceChangeWatcher = &watchor.ResourceChangeWatchClient{}

type ResourceChangeWatcher interface {
	Start(params *watchor.ResourceChangeWatchStartParams) error
	Channel() <-chan *models.ResourceChangeEvent
	ErrorChannel() <-chan *watchor.ErrorEvent
	WarningChannel() <-chan *watchor.WarningEvent
}

type OptionFunc func(*Options)

func SetUserInfo(u *client.UserInfo) OptionFunc {
	return func(o *Options) {
		o.UserInfo = u
	}
}

func SetAPIAuth(username string, password string) OptionFunc {
	return func(o *Options) {
		o.APIPassword = password
		o.APIUsername = username
	}
}

func SetHost(host string) OptionFunc {
	return func(o *Options) {
		o.Host = host
	}
}

func SetPollingInterval(i time.Duration) OptionFunc {
	return func(o *Options) {
		o.PollingInterval = i
	}
}

type Options struct {
	UserInfo        *client.UserInfo
	Host            string
	APIUsername     string
	APIPassword     string
	PollingInterval time.Duration
}

func NewWatchClient(resourceTypes []string, opts *Options) (ResourceChangeWatcher, error) {
	c, err := NewWatchOriClient(resourceTypes, opts)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("nil crc watch client")
	}
	return c, nil
}

func NewWatchOriClient(resourceTypes []string, opts *Options) (*watchor.ResourceChangeWatchClient, error) {
	towerclient, err := apiclient.NewWithUserConfig(apiclient.ClientConfig{
		Host:     opts.Host,
		BasePath: "v2/api",
		Schemes:  []string{"http"},
	}, apiclient.UserConfig{
		Name:     opts.UserInfo.Username,
		Password: opts.UserInfo.Password,
		Source:   models.UserSource(opts.UserInfo.Source),
	})

	if err != nil {
		klog.Errorf("Failed to init api client, err: %s", err)
		return nil, err
	}

	var options resource_change_client.ClientOption = func(op *runtime.ClientOperation) {
		op.AuthInfo = httptransport.BasicAuth(opts.APIUsername, opts.APIPassword)
		op.Params = NewBypassWhiteListHeader()
	}

	crcWatchClient, err := watchor.NewResourceChangeWatchClient(&watchor.NewResourceChangeWatchClientParams{
		Client:          towerclient,
		ResourceID:      nil,
		PollingInterval: opts.PollingInterval,
		ClientOptions:   options,
		ResourceTypes:   resourceTypes,
	})

	if err != nil {
		klog.Errorf("Failed to init crc client, err: %s", err)
		return nil, err
	}
	return crcWatchClient, nil
}
