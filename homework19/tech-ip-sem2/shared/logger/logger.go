package logger

import "go.uber.org/zap"

func New() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}

	l, err := cfg.Build()
	if err != nil {
		panic("logger: failed to build: " + err.Error())
	}
	return l
}

func Port(port string) zap.Field {
	return zap.String("port", port)
}

func Err(err error) zap.Field {
	return zap.Error(err)
}
