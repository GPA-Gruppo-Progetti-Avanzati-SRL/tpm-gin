package main

import (
	_ "embed"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/httpsrv"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/httpsrv/resource/health"
	_ "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/httpsrv/resource/metrics"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AppConfig struct {
	Http       httpsrv.Config
	MwRegistry middleware.HandlerCatalogConfig `yaml:"mw-handler-registry" mapstructure:"mw-handler-registry"`
}

type S1 struct {
	Nome string
}

func (m *AppConfig) PostProcess() error {
	return nil
}

/*
func (m *AppConfig) GetDefaults() []configuration.VarDefinition {

	vd := make([]configuration.VarDefinition, 0, 20)
	vd = append(vd, httpsrv.GetConfigDefaults()...)
	vd = append(vd, middleware.GetConfigDefaults("config.mw-handler-registry")...)
	return vd
}
*/

//go:embed config.yml
var configFile []byte

func main() {

	appCfg := AppConfig{}

	/*
		_, err := configuration.NewConfiguration(
			configuration.WithType("yaml"),
			configuration.WithName("tpm-gin"),
			configuration.WithReader(bytes.NewBuffer([]byte(configFile))),
			configuration.WithData(&appCfg))
		if nil != err {
			log.Fatal().Err(err).Send()
		}
	*/

	log.Info().Msgf("read in config is: %+v\n", appCfg)

	if appCfg.MwRegistry != nil {
		if err := middleware.InitializeHandlerRegistry(appCfg.MwRegistry, appCfg.Http.MwUse); err != nil {
			log.Fatal().Err(err).Send()
		}
	}

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	jc, err := initGlobalTracer()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer jc.Close()

	s, err := httpsrv.NewServer(appCfg.Http, httpsrv.WithListenPort(9090), httpsrv.WithDocumentRoot("/www", "/tmp", false))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := s.Start(); err != nil {
		log.Fatal().Err(err).Send()
	}
	defer s.Stop()

	for !s.IsReady() {
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

	sig := <-shutdownChannel
	log.Debug().Interface("signal", sig).Msg("got termination signal")
}

func initGlobalTracer() (io.Closer, error) {
	var tracer opentracing.Tracer
	var closer io.Closer

	jcfg := jaegercfg.Configuration{
		ServiceName: "gintest",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: 1.0,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := jcfg.NewTracer(
		jaegercfg.Logger(&jlogger{}),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if nil != err {
		log.Error().Err(err).Msg("Error in NewTracer")
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}

type jlogger struct{}

func (l *jlogger) Error(msg string) {
	log.Error().Msg("(jaeger) " + msg)
}

func (l *jlogger) Infof(msg string, args ...interface{}) {
	log.Info().Msgf("(jaeger) "+msg, args...)
}
