package mailslurper

/*
NoWorkerAvailableError is an error used when no worker is available to
service a SMTP connection request.
*/
type NoWorkerAvailableError struct{}

/*
NoWorkerAvailable returns a new instance of the No Worker Available error
*/
func NoWorkerAvailable() NoWorkerAvailableError {
	return NoWorkerAvailableError{}
}

func (err NoWorkerAvailableError) Error() string {
	return "No worker available. Timeout has been exceeded"
}
