/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client

import (
	"context"
	"testing"

	"github.com/liaogang/hertz/pkg/protocol"
)

var (
	biz       = "Biz"
	beforeMW0 = "BeforeMiddleware0"
	afterMW0  = "AfterMiddleware0"
	beforeMW1 = "BeforeMiddleware1"
	afterMW1  = "AfterMiddleware1"
)

func invoke(ctx context.Context, req *protocol.Request, resp *protocol.Response) (err error) {
	req.BodyBuffer().WriteString(biz)
	return nil
}

func mockMW0(next Endpoint) Endpoint {
	return func(ctx context.Context, req *protocol.Request, resp *protocol.Response) (err error) {
		req.BodyBuffer().WriteString(beforeMW0)
		err = next(ctx, req, resp)
		if err != nil {
			return err
		}
		req.BodyBuffer().WriteString(afterMW0)
		return nil
	}
}

func mockMW1(next Endpoint) Endpoint {
	return func(ctx context.Context, req *protocol.Request, resp *protocol.Response) (err error) {
		req.BodyBuffer().WriteString(beforeMW1)
		err = next(ctx, req, resp)
		if err != nil {
			return err
		}
		req.BodyBuffer().WriteString(afterMW1)
		return nil
	}
}

func TestChain(t *testing.T) {
	mws := chain(mockMW0, mockMW1)
	req := protocol.AcquireRequest()
	mws(invoke)(context.Background(), req, nil)
	final := beforeMW0 + beforeMW1 + biz + afterMW1 + afterMW0
	if req.BodyBuffer().String() != final {
		t.Errorf("unexpected %#v, expected %#v", req.BodyBuffer().String(), final)
	}
}
