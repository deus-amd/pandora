// Copyright (c) 2017 Yandex LLC. All rights reserved.
// Use of this source code is governed by a MPL 2.0
// license that can be found in the LICENSE file.
// Author: Vladimir Skipor <skipor@yandex-team.ru>

package schedule

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/yandex/pandora/core"
	"github.com/yandex/pandora/core/coretest"
	"go.uber.org/atomic"
)

var _ = Describe("composite", func() {
	It("empty", func() {
		testee := NewComposite()
		coretest.ExpectScheduleNexts(testee, 0)
	})

	It("only", func() {
		testee := NewComposite(NewConst(1, time.Second))
		coretest.ExpectScheduleNexts(testee, time.Second, time.Second)
	})

	It("composite", func() {
		testee := NewComposite(
			NewConst(1, 2*time.Second),
			NewOnce(2),
			NewConst(0, 5*time.Second),
			NewOnce(0),
			NewOnce(1),
		)
		coretest.ExpectScheduleNexts(testee,
			time.Second,
			2*time.Second,

			2*time.Second,
			2*time.Second,

			7*time.Second,
			7*time.Second, // Finish.
		)
	})

	// Load concurrently, and let race detector do it's work.
	It("race", func() {
		var (
			nexts          []core.Schedule
			tokensGot      atomic.Int64
			tokensExpected int64
		)
		addOnce := func(v int64) {
			nexts = append(nexts, NewOnce(v))
			tokensExpected += v
		}
		addOnce(100000) // Delay to start concurrent readers.
		for i := 0; i < 100000; i++ {
			// Some work for catching races.
			addOnce(0)
			addOnce(1)
		}
		testee := NewCompositeConf(CompositeConf{nexts})
		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					_, ok := testee.Next()
					if !ok {
						return
					}
					tokensGot.Inc()
				}
			}()
		}
		wg.Wait()
		Expect(tokensGot.Load()).To(Equal(tokensExpected))
	})
})
