package suspender

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	suite.Run(t, &testSuspender{})
}

type testSuspender struct {
	suite.Suite

	suspender *Suspender[uint64]
}

func (t *testSuspender) SetupTest() {
	t.suspender = New[uint64](Config{})
}

func (t *testSuspender) TestIncTwice() {
	testValue := uint64(1)

	t.Require().NoError(t.suspender.Inc(testValue))

	t.Require().ErrorIs(t.suspender.Inc(testValue), ErrCountOverflow)
}

func (t *testSuspender) TestIncDifferent() {
	testValue := uint64(1)
	testValue2 := uint64(2)

	t.Require().NoError(t.suspender.Inc(testValue))
	t.Require().NoError(t.suspender.Inc(testValue2))
}

func (t *testSuspender) TestIncDrop() {
	testValue := uint64(1)

	t.Require().NoError(t.suspender.Inc(testValue))
	t.Require().NoError(t.suspender.Dec(testValue))
	t.Require().NoError(t.suspender.Inc(testValue))
}

func (t *testSuspender) TestParallelInc() {
	testValue := uint64(1)

	t.Require().NoError(t.suspender.Inc(testValue))

	wg := sync.WaitGroup{}

	for index := uint64(0); index < 100000; index++ {
		wg.Add(1)

		go func(i uint64) {
			defer wg.Done()

			t.Require().ErrorIs(t.suspender.Inc(testValue), ErrCountOverflow)
		}(index)
	}

	wg.Wait()
}

func (t *testSuspender) TestIncWithCtx() {
	ctx, cancel := context.WithCancel(context.Background())

	testValue := uint64(1)

	t.Require().NoError(t.suspender.IncWithCtx(ctx, testValue))
	t.Require().ErrorIs(t.suspender.IncWithCtx(ctx, testValue), ErrCountOverflow)
	t.Require().ErrorIs(t.suspender.Inc(testValue), ErrCountOverflow)

	cancel()

	t.Require().ErrorIs(t.suspender.IncWithCtx(ctx, testValue), context.Canceled)
	// t.suspender.Inc(testValue) -- error before gorotine will see that context done
}
