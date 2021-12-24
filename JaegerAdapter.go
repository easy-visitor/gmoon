package gmoon

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

type JaegerAdapter struct {
	Tracer opentracing.Tracer
}

func (this *JaegerAdapter) Name() string {

	return "JaegerAdapter"
}

func NewJaegerAdapter() *JaegerAdapter {
	jcfg := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		ServiceName: "serviceName",
	}
	report := jaegerConfig.ReporterConfig{
		LogSpans:           true,
		LocalAgentHostPort: globalConfig.Jaeger.LocalAgentHostPort,
	}
	reporter, _ := report.NewReporter(globalConfig.Jaeger.ServiceName, jaeger.NewNullMetrics(), jaeger.NullLogger)
	tracer, _, err := jcfg.NewTracer(
		jaegerConfig.Reporter(reporter),
	)
	if err != nil {

		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	//defer closer.Close()

	return &JaegerAdapter{
		Tracer: tracer,
	}
}
