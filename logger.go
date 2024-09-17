package healthy

//go:generate mockgen -destination=./mock_logger.go -package=healthy -source=logger.go Logger
type Logger interface {
	Info(args ...interface{})
}
