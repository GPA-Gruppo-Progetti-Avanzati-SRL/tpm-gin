package sysmiddleware_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/sysmiddleware"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestPromHttpMetricsHandlerConfig(t *testing.T) {

	b, err := yaml.Marshal(sysmiddleware.DefaultMetricsConfig)
	require.NoError(t, err)

	t.Log(string(b))
}
