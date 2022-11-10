// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package async

import (
	"errors"
)

const (
	JobStatusInProgress string = "in_progress"
	JobStatusComplete   string = "completed"

	AudioTypeMP3  string = "mp3"
	AudioTypeMpeg string = "mpeg"
	AudioTypeWav  string = "wav"
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")

	// ErrJobStatusTimeout the job status check timed out
	ErrJobStatusTimeout = errors.New("the job status check timed out")

	// ErrInvalidWaitTime the time to wait agurment is invalid
	ErrInvalidWaitTime = errors.New("the time to wait agurment is invalid")

	// ErrInvalidURIExtension couldn't find a period to indicate a file extension
	ErrInvalidURIExtension = errors.New("couldn't find a period to indicate a file extension")
)
