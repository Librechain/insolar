///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package main

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"math"
	"sync"
	"time"
)

func main() {
	ctx := initLogger(insolar.DebugLevel)

	wg := sync.WaitGroup{}

	parallelism := 10
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go generateLogs(ctx, i, wg)
	}
	wg.Wait()
}

func generateLogs(ctx context.Context, id int, wg sync.WaitGroup) {
	defer wg.Done()
	for i := uint32(0); i < math.MaxUint32; i++ {
		inslogger.FromContext(ctx).Debugf("some text: id=%d iter=%d time=%v", id, i, time.Now())
		time.Sleep(time.Millisecond)
	}
}

func initLogger(level insolar.LogLevel) context.Context {
	ctx := context.Background()
	logger := inslogger.FromContext(ctx).WithCaller(false)
	logger, _ = logger.WithLevelNumber(level)
	logger, _ = logger.WithFormat(insolar.TextFormat)
	ctx = inslogger.SetLogger(ctx, logger)
	return ctx
}
