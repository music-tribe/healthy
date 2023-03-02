package healthy

import "github.com/music-tribe/errors"

type Service interface {
	Name() string
	Checkers() []Checker
	Version() string
}

type service struct {
	serviceName string `validate:"required"`
	version     string `validate:"required,semver"`
	checkers    []Checker
}

func New(serviceName, version string, checkers ...Checker) (Service, error) {
	c := make([]Checker, 0)
	s := new(service)
	s.checkers = checkers
	if checkers == nil {
		s.checkers = c
	}

	if serviceName == "" {
		return s, errors.NewCloudError(400, "serviceName param is empty")
	}

	if version == "" {
		return s, errors.NewCloudError(400, "version param is empty")
	}

	s.serviceName = serviceName
	s.version = version

	return s, nil
}

func (s *service) Checkers() []Checker {
	return s.checkers
}

func (s *service) Name() string {
	return s.serviceName
}

func (s *service) Version() string {
	return s.version
}
