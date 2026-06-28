package proxy

import "errors"

var ErrJobsServiceClosed = errors.New("jobs service is closed")
var ErrJobQueueFull = errors.New("job queue full")
