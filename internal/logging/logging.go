package logging

import (
	"context"
	"europm/internal/config"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var tracer oteltrace.Tracer

// var tp *sdktrace.TracerProvider
var logger *zap.Logger
var logLevel zapcore.Level
var shutdown func(context.Context) error

func Init() (err error) {
	// Init logger
	var cores []zapcore.Core

	if config.GetBool("logger.enable_console") {
		level := getZapLevel(config.GetString("logger.console_level"))
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getZapEncoder(config.GetBool("logger.file_json_format")), writer, level)
		cores = append(cores, core)
		logLevel = level
	}

	/*if config.EnableFile {
		level := getZapLevel(config.FileLevel)
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.FileLocation,
			MaxSize:    100,
			Compress:   true,
			MaxAge:     0,
			MaxBackups: 0,
		})
		core := zapcore.NewCore(getEncoder(config.FileJSONFormat), writer, level)
		cores = append(cores, core)
	}*/

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skips 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	logger = zap.New(
		combinedCore,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	)

	// Init tracer
	/*tp, _ = tracerProvider(config.GetString("tracer.url"))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))*/

	shutdown, err = initTracerProvider()
	if err != nil {
		return err
	}

	tracer = otel.Tracer("go.opentelemetry.io/otel")

	return err
}

func Destroy() {
	if err := shutdown(context.Background()); err != nil {
		Errorf(nil, "failed to shutdown TracerProvider: %w", err)
	}
	/*if err := tp.Shutdown(context.Background()); err != nil {
		Errorf(nil, "error shutting down tracer provider: ", err)
	}*/
}

func StartTrace(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	return tracer.Start(ctx, spanName, opts...)
}

func Debug(span oteltrace.Span, args ...interface{}) {
	Debugf(span, "", args...)
}

func Debugf(span oteltrace.Span, template string, args ...interface{}) {
	logf(zapcore.DebugLevel, span, template, args...)
}

func Info(span oteltrace.Span, args ...interface{}) {
	Infof(span, "", args...)
}

func Infof(span oteltrace.Span, template string, args ...interface{}) {
	logf(zapcore.InfoLevel, span, template, args...)
}

func Warn(span oteltrace.Span, args ...interface{}) {
	Warnf(span, "", args...)
}

func Warnf(span oteltrace.Span, template string, args ...interface{}) {
	logf(zapcore.WarnLevel, span, template, args...)
}

func Error(span oteltrace.Span, args ...interface{}) {
	Errorf(span, "", args...)
}

func Errorf(span oteltrace.Span, template string, args ...interface{}) {
	logf(zapcore.ErrorLevel, span, template, args...)
}

func Fatal(span oteltrace.Span, args ...interface{}) {
	Fatalf(span, "", args...)
}

func Fatalf(span oteltrace.Span, template string, args ...interface{}) {
	logf(zapcore.FatalLevel, span, template, args...)
}

func Panic(span oteltrace.Span, args ...interface{}) {
	Panicf(span, "", args...)
}

func Panicf(span oteltrace.Span, template string, args ...interface{}) {
	logf(zapcore.PanicLevel, span, template, args...)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func getZapEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time" // This will change the key from ts to time
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
/*func tracerProvider(url string) (*sdktrace.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	//hostName, _ := os.Hostname()
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in an Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("ipn-go"),
			//attribute.String("host_name", hostName),
			attribute.Int("pid", os.Getpid()),
			attribute.String("start_time", time.Now().UTC().Format("2006-01-02T15:04:05Z")),
		)),
	)

	return tp, nil
}*/

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initTracerProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("iportal-fdm"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	// conn, err := grpc.DialContext(ctx, "localhost:30080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	// }

	// Set up a trace exporter
	// traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	// }

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	// bsp := trace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func logf(level zapcore.Level, span oteltrace.Span, template string, args ...interface{}) {
	if level >= logLevel {
		msg := getMessage(template, args)
		switch level {
		case zapcore.DebugLevel:
			logger.Debug(msg)
		case zapcore.InfoLevel:
			logger.Info(msg)
		case zapcore.WarnLevel:
			logger.Warn(msg)
		case zapcore.ErrorLevel:
			logger.Error(msg)
		case zapcore.FatalLevel:
			logger.Fatal(msg)
		case zapcore.PanicLevel:
			logger.Panic(msg)
		default:
			logger.Info(msg)
		}
		if span != nil {
			span.AddEvent(msg)
		}
	}
}

// getMessage format with Sprint, Sprintf, or neither.
func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}
