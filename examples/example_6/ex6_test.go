package example_6_test

import (
	"bytes"
	_ "embed"
	"gitlab.alm.poste.it/go/configuration"
	"testing"
	"tpm-gin/httpsrv"
	"tpm-gin/middleware"
)

type AppConfig struct {
	Http       httpsrv.Config
	MwRegistry map[string]interface{} `yaml:"mw-handler-registry" mapstructure:"mw-handler-registry"`
}

func (m *AppConfig) PostProcess() error {
	return nil
}

func (m *AppConfig) GetDefaults() []configuration.VarDefinition {

	vd := make([]configuration.VarDefinition, 0, 20)
	vd = append(vd, httpsrv.GetConfigDefaults()...)
	vd = append(vd, middleware.GetConfigDefaults("config.mw-handler-registry")...)
	return vd
}

//go:embed config.yml
var configFile []byte

func TestConfigFile(t *testing.T) {

	appCfg := &AppConfig{}

	_, err := configuration.NewConfiguration(
		configuration.WithType("yaml"),
		configuration.WithName("tpm-gin"),
		configuration.WithReader(bytes.NewBuffer(configFile)),
		configuration.WithData(appCfg))
	if nil != err {
		t.Fatal(err)
	}

	t.Logf("%+v\n", appCfg)

	if appCfg.MwRegistry != nil {
		if err := middleware.InitializeHandlerRegistry(appCfg.MwRegistry); err != nil {
			t.Fatal(err)
		}
	}
}
