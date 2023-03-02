package healthy

const defaultName = "healthcheck"

type Checker interface {
	Check() error
	Name() string
}

func NewChecker(name string, pinger Pinger) Checker {
	if name == "" {
		name = defaultName
	}
	if pinger == nil {
		pinger = new(stubPinger)
	}
	return &checker{
		name:   name,
		pinger: pinger,
	}
}

func NewMongoDbCheckerWithConnectionString(name, dbConnString string) Checker {
	return NewChecker(name, &MongoDbPinger{connString: dbConnString})
}

type checker struct {
	name   string
	pinger Pinger
}

func (c *checker) Check() error {
	return c.pinger.Ping()
}

func (c *checker) Name() string {
	return c.name
}
