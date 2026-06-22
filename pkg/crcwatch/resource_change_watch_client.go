package crcwatch

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/smartxworks/cloudtower-go-sdk/v2/client"
	resource_change_client "github.com/smartxworks/cloudtower-go-sdk/v2/client/resource_change"
	userclient "github.com/smartxworks/cloudtower-go-sdk/v2/client/user"
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

func SetScheme(scheme string) OptionFunc {
	return func(o *Options) {
		o.Scheme = scheme
	}
}

func SetAllowInsecure(allow bool) OptionFunc {
	return func(o *Options) {
		o.AllowInsecure = allow
	}
}

func SetLimit(l int32) OptionFunc {
	return func(o *Options) {
		o.Limit = l
	}
}

func SetPollingInterval(i time.Duration) OptionFunc {
	return func(o *Options) {
		o.PollingInterval = i
	}
}

func SetCatchUpPollingInterval(i time.Duration) OptionFunc {
	return func(o *Options) {
		o.CatchUpPollingInterval = i
	}
}

type Options struct {
	UserInfo               *client.UserInfo
	Host                   string
	Scheme                 string
	AllowInsecure          bool
	APIUsername            string
	APIPassword            string
	PollingInterval        time.Duration
	CatchUpPollingInterval time.Duration
	Limit                  int32
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
	scheme := opts.Scheme
	if scheme == "" {
		scheme = "http"
	}

	towerclient, err := newTowerClient(apiclient.ClientConfig{
		Host:     opts.Host,
		BasePath: "v2/api",
		Schemes:  []string{scheme},
	}, apiclient.UserConfig{
		Name:     opts.UserInfo.Username,
		Password: opts.UserInfo.Password,
		Source:   models.UserSource(opts.UserInfo.Source),
	}, opts.AllowInsecure)

	if err != nil {
		klog.Errorf("Failed to init api client, err: %s", err)
		return nil, err
	}

	var options resource_change_client.ClientOption = func(op *runtime.ClientOperation) {
		op.AuthInfo = httptransport.Compose(
			op.AuthInfo,
			httptransport.BasicAuth(opts.APIUsername, opts.APIPassword),
		)
		op.Params = NewBypassWhiteListHeader(op.Params)
	}

	crcWatchClient, err := watchor.NewResourceChangeWatchClient(&watchor.NewResourceChangeWatchClientParams{
		Client:                 towerclient,
		ResourceID:             nil,
		PollingInterval:        opts.PollingInterval,
		CatchUpPollingInterval: opts.CatchUpPollingInterval,
		Limit:                  opts.Limit,
		ClientOptions:          options,
		ResourceTypes:          resourceTypes,
	})

	if err != nil {
		klog.Errorf("Failed to init crc client, err: %s", err)
		return nil, err
	}
	return crcWatchClient, nil
}

func newTowerClient(clientConfig apiclient.ClientConfig, userConfig apiclient.UserConfig, insecure bool) (*apiclient.Cloudtower, error) {
	tlsConfig, err := httptransport.TLSClientAuth(httptransport.TLSClientOptions{
		InsecureSkipVerify: insecure, // #nosec G402
	})
	if err != nil {
		return nil, err
	}
	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.TLSClientConfig = tlsConfig
	httpClient := &http.Client{Transport: httpTransport}

	transport := httptransport.NewWithClient(clientConfig.Host, clientConfig.BasePath, clientConfig.Schemes, httpClient)
	client := apiclient.New(transport, strfmt.Default)

	params := userclient.NewLoginParams()
	params.RequestBody = &models.LoginInput{
		Username: &userConfig.Name,
		Password: &userConfig.Password,
		Source:   userConfig.Source.Pointer(),
	}

	resp, err := client.User.Login(params)
	if err != nil {
		return nil, err
	}
	if resp.Payload == nil || resp.Payload.Data == nil || resp.Payload.Data.Token == nil {
		return nil, errors.New("login response missing token")
	}

	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", *resp.Payload.Data.Token)
	return client, nil
}
