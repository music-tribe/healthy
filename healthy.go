package healthy

type Service interface {
	Name() string
	Checkers() []Checker
	Version() string
}

type service struct {
	serviceName, version string
	checkers             []Checker
}

func New(serviceName, version string, checkers ...Checker) Service {
	if checkers == nil {
		checkers = make([]Checker, 0)
	}
	return &service{
		serviceName: serviceName,
		version:     version,
		checkers:    checkers,
	}
}

func (s *service) Checkers() []Checker {
	return s.checkers
}

func (s *service) Name() string {
	return s.serviceName
}

func (s *service) Version() string {
	return s.serviceName
}
