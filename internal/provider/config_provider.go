package provider

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/provider/paramstore"
	"github.com/spf13/pflag"
	"sort"
	"strings"
)

type Interface interface {
	GetValue(key string) (string, error)
	GetValueTree(prefix string) (map[string]string, error)
	SetValue(key string, value string) error
}

type Options struct {
	ProviderName string
	Provider     Interface
	Region       string // a.k.a. Datacenter
	Namespace    string // a.k.a. Environment
	Debug        bool
}

// AddProviderOptions adds flags to a pflag.FlagSet
func (o *Options) AddProviderOptions(c *pflag.FlagSet) {
	c.StringVar(&o.Region, "region", "default", "region")
	c.StringVarP(&o.Namespace, "namespace", "n", "default", "namespace")
}

var SupportedProviders []string

func init() {
	SupportedProviders = []string{"aws"}
}

func New(o Options) (Interface, error) {
	// provider-by-key pattern
	//supportedProviders := map[string]func() app.Interface{
	//	"aws": func() app.Interface { return paramstore.NewSSMClient(o.Debug) },
	//}
	//
	//supportedProviderKeys := make([]string, 0)
	//for k, _ := range supportedProviders {
	//	supportedProviderKeys = append(supportedProviderKeys, k)
	//}
	//sort.Strings(supportedProviderKeys)
	//
	//if clientFunc, ok := supportedProviders[o.ProviderName]; ok {
	//	return clientFunc(), nil
	//} else {
	//	return nil, fmt.Errorf("unrecognized provider %#v: supported provider are: %#v", o.ProviderName, strings.Join(supportedProviderKeys, ", "))
	//}

	// provider-as-list pattern
	if strings.EqualFold(o.ProviderName, "aws") {
		return paramstore.NewSSMClient(o.Region, o.Debug), nil
	} else {
		sort.Strings(SupportedProviders)
		return nil, fmt.Errorf("unrecognized provider %#v: supported provider are: %#v", o.ProviderName, strings.Join(SupportedProviders, ", "))
	}
}
