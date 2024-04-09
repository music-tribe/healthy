package healthy

import (
	"fmt"

	"github.com/music-tribe/errors"
)

type Service interface {
	Name() string
	Checkers() map[string]Checker
	Version() string
}

type service struct {
	serviceName string `validate:"required"`
	version     string `validate:"required,semver"`
	checkers    map[string]Checker
}

func New(serviceName, version string, checkers ...Checker) (Service, error) {
	s := new(service)
	s.checkers = make(map[string]Checker)
	if serviceName == "" {
		return s, errors.NewCloudError(400, "serviceName param is empty")
	}

	if version == "" {
		return s, errors.NewCloudError(400, "version param is empty")
	}

	s.serviceName = serviceName
	s.version = version

	for _, checker := range checkers {
		_, ok := s.checkers[checker.Name()]
		if ok {
			return s, errors.NewCloudError(400, fmt.Sprintf("Duplicate checker name: %s", checker.Name()))
		}
		s.checkers[checker.Name()] = checker
	}

	return s, nil
}

func (s *service) Checkers() map[string]Checker {
	return s.checkers
}

func (s *service) Name() string {
	return s.serviceName
}

func (s *service) Version() string {
	return s.version
}
