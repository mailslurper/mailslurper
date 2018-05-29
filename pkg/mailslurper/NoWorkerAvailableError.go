// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

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
