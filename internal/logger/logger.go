package logger

import "go.uber.org/zap"

func NewLogger() (*zap.SugaredLogger, error) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return zapLogger.Sugar(), nil
}
