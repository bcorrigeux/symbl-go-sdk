// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package interfaces

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

/*
	custom headers and configuration options
*/
// signer
type Signer interface {
	SignRequest(*http.Request) error
}

type SignerContext struct{}

func WithSigner(ctx context.Context, s Signer) context.Context {
	return context.WithValue(ctx, SignerContext{}, s)
}

// header
type HeadersContext struct{}

func WithCustomHeaders(ctx context.Context, headers http.Header) context.Context {
	return context.WithValue(ctx, HeadersContext{}, headers)
}

// query
type ParametersContext struct{}

func WithCustomParameters(ctx context.Context, params map[string][]string) context.Context {
	return context.WithValue(ctx, ParametersContext{}, params)
}

/*
	RawResponse may be used with the Do method as the resBody argument in order
	to capture the raw response data.
*/
type RawResponse struct {
	bytes.Buffer
}

/*
	Error handling
*/
type StatusError struct {
	Resp *http.Response
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Resp.Request.Method, e.Resp.Request.URL, e.Resp.Status)
}