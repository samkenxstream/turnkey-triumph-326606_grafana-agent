// Package apache_exporter embeds https://github.com/Lusitaniae/apache_exporter
package apache_exporter //nolint:golint

import (
	"fmt"
	"net/http"
	"net/url"

	ae "github.com/Lusitaniae/apache_exporter/collector"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	integrations_v2 "github.com/grafana/agent/pkg/integrations/v2"
	"github.com/grafana/agent/pkg/integrations/v2/common"
	"github.com/grafana/agent/pkg/integrations/v2/metricsutils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DefaultConfig holds the default settings for the apache_exporter integration
var DefaultConfig = Config{
	ApacheAddr:         "http://localhost/server-status?auto",
	ApacheHostOverride: "",
	ApacheInsecure:     false,
}

// Config controls the apache_exporter integration.
type Config struct {
	ApacheAddr         string               `yaml:"scrape_uri,omitempty"`
	ApacheHostOverride string               `yaml:"host_override,omitempty"`
	ApacheInsecure     bool                 `yaml:"insecure,omitempty"`
	Common             common.MetricsConfig `yaml:",inline"`
}

// ApplyDefaults applies the integration's default configuration.
func (c *Config) ApplyDefaults(globals integrations_v2.Globals) error {
	c.Common.ApplyDefaults(globals.SubsystemOpts.Metrics.Autoscrape)
	return nil
}

// Identifier returns a string that identifies the integration.
func (c *Config) Identifier(globals integrations_v2.Globals) (string, error) {
	return *c.Common.InstanceKey, nil
}

// UnmarshalYAML implements yaml.Unmarshaler for Config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig

	type plain Config
	return unmarshal((*plain)(c))
}

// Name returns the name of the integration this config is for.
func (c *Config) Name() string {
	return "apache"
}

// InstanceKey returns the addr of the apache server.
func (c *Config) InstanceKey(agentKey string) (string, error) {

	u, err := url.Parse(c.ApacheAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", u.Hostname(), u.Port()), nil
}

type apacheHandler struct {
	cfg *Config
	log log.Logger
}

// NewIntegration instantiates a new integrations.MetricsIntegration
// which will handle requests to the apache exporter.
func (c *Config) NewIntegration(logger log.Logger, globals integrations_v2.Globals) (integrations_v2.Integration, error) {
	var metricsCfg common.MetricsConfig
	metricsCfg.ApplyDefaults(globals.SubsystemOpts.Metrics.Autoscrape)

	_, err := url.ParseRequestURI(c.ApacheAddr)
	if err != nil {
		level.Error(logger).Log("msg", "scrape_uri is invalid", "err", err)
		return nil, err
	}

	ah := &apacheHandler{cfg: c, log: logger}
	h, err := ah.createHandler()
	if err != nil {
		return nil, err
	}

	return metricsutils.NewMetricsHandlerIntegration(logger, c, metricsCfg, globals, h)
}

func (ah *apacheHandler) createHandler() (http.HandlerFunc, error) {

	aeExporter := ae.NewExporter(ah.log, &ae.Config{
		ScrapeURI:    ah.cfg.ApacheAddr,
		HostOverride: ah.cfg.ApacheHostOverride,
		Insecure:     ah.cfg.ApacheInsecure,
	})

	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		registry.MustRegister(aeExporter)
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}, nil
}

func init() {
	integrations_v2.Register(&Config{}, integrations_v2.TypeMultiplex)
}
