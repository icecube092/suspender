package suspender_test

import (
	"sync"
	"testing"

	"suspender"

	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	suite.Run(t, &testSuspender{})
}

type testSuspender struct {
	suite.Suite

	suspender *suspender.Suspender[uint64]
}

func (t *testSuspender) SetupTest() {
	t.suspender = suspender.New[uint64](suspender.Config{})
}

func (t *testSuspender) TestAddTwice() {
	testValue := uint64(1)

	t.Require().NoError(t.suspender.Add(testValue))

	t.Require().ErrorIs(t.suspender.Add(testValue), suspender.ErrCountOverflow)
}

func (t *testSuspender) TestAddDrop() {
	testValue := uint64(1)

	t.Require().NoError(t.suspender.Add(testValue))
	t.Require().NoError(t.suspender.Drop(testValue))
	t.Require().NoError(t.suspender.Add(testValue))
}

func (t *testSuspender) TestParallelAdd() {
	testValue := uint64(1)

	t.Require().NoError(t.suspender.Add(testValue))

	wg := sync.WaitGroup{}

	for index := uint64(0); index < 100000; index++ {
		wg.Add(1)

		go func(i uint64) {
			defer wg.Done()

			t.Require().ErrorIs(t.suspender.Add(testValue), suspender.ErrCountOverflow)
		}(index)
	}

	wg.Wait()
}
