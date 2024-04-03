package ports

//go:generate mockgen -destination=./mocks/mock_logger.go -package=mocklogger -source=logger.go Logger
type Logger interface {
	Info(component, method string, args ...interface{})
	Debug(component, method string, args ...interface{})
	Warn(component, method string, args ...interface{})
	Fatal(component, method string, args ...interface{})
	Error(component, method string, args ...interface{})
	Panic(component, method string, args ...interface{})
}
