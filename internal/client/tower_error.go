package client

type towerError struct {
	err        error
	statusCode int
}

func (e towerError) Error() string {
	return e.err.Error()
}

func newTowerError(err error, statusCode int) towerError {
	return towerError{
		err:        err,
		statusCode: statusCode,
	}
}
