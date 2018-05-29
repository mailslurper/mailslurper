// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
SMTPWorkerState defines states that a worker may be in. Typically
a worker starts IDLE, the moves to WORKING, finally going to
either DONE or ERROR.
*/
type SMTPWorkerState int

const (
	SMTP_WORKER_IDLE    SMTPWorkerState = 0
	SMTP_WORKER_WORKING SMTPWorkerState = 1
	SMTP_WORKER_DONE    SMTPWorkerState = 100
	SMTP_WORKER_ERROR   SMTPWorkerState = 101

	RECEIVE_BUFFER_LEN         = 1024
	CONNECTION_TIMEOUT_MINUTES = 10
	COMMAND_TIMEOUT_SECONDS    = 5
)
