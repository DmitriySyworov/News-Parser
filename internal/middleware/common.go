package middleware

import (
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/response"
)

type ManagerMiddleware struct {
	Signature string
	resp      response.Response
	*ContextValues
	Logger loggers.Logger
}
type ContextValues struct {
	SessionID string
	UserUUID  string
	DataLog   loggers.DataLog
}

const (
	KeyContextValues = "keyContextValues"
)

func NewManagerMiddleware(signature string, logger *loggers.Logger) *ManagerMiddleware {
	return &ManagerMiddleware{
		Signature: signature,
		Logger:    *logger,
		ContextValues: &ContextValues{
			DataLog: *logger.DataLog,
		},
	}
}
