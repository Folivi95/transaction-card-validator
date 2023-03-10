// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"github.com/saltpay/transaction-card-validator/internal/adapters/schemaregistry"
	"sync"
	"time"
)

// Ensure, that RefreshSchedulerMock does implement schemaregistry.RefreshScheduler.
// If this is not the case, regenerate this file with moq.
var _ schemaregistry.RefreshScheduler = &RefreshSchedulerMock{}

// RefreshSchedulerMock is a mock implementation of schemaregistry.RefreshScheduler.
//
// 	func TestSomethingThatUsesRefreshScheduler(t *testing.T) {
//
// 		// make and configure a mocked schemaregistry.RefreshScheduler
// 		mockedRefreshScheduler := &RefreshSchedulerMock{
// 			AfterFuncFunc: func(t time.Duration, f func()) *time.Timer {
// 				panic("mock out the AfterFunc method")
// 			},
// 		}
//
// 		// use mockedRefreshScheduler in code that requires schemaregistry.RefreshScheduler
// 		// and then make assertions.
//
// 	}
type RefreshSchedulerMock struct {
	// AfterFuncFunc mocks the AfterFunc method.
	AfterFuncFunc func(t time.Duration, f func()) *time.Timer

	// calls tracks calls to the methods.
	calls struct {
		// AfterFunc holds details about calls to the AfterFunc method.
		AfterFunc []struct {
			// T is the t argument value.
			T time.Duration
			// F is the f argument value.
			F func()
		}
	}
	lockAfterFunc sync.RWMutex
}

// AfterFunc calls AfterFuncFunc.
func (mock *RefreshSchedulerMock) AfterFunc(t time.Duration, f func()) *time.Timer {
	if mock.AfterFuncFunc == nil {
		panic("RefreshSchedulerMock.AfterFuncFunc: method is nil but RefreshScheduler.AfterFunc was just called")
	}
	callInfo := struct {
		T time.Duration
		F func()
	}{
		T: t,
		F: f,
	}
	mock.lockAfterFunc.Lock()
	mock.calls.AfterFunc = append(mock.calls.AfterFunc, callInfo)
	mock.lockAfterFunc.Unlock()
	return mock.AfterFuncFunc(t, f)
}

// AfterFuncCalls gets all the calls that were made to AfterFunc.
// Check the length with:
//     len(mockedRefreshScheduler.AfterFuncCalls())
func (mock *RefreshSchedulerMock) AfterFuncCalls() []struct {
	T time.Duration
	F func()
} {
	var calls []struct {
		T time.Duration
		F func()
	}
	mock.lockAfterFunc.RLock()
	calls = mock.calls.AfterFunc
	mock.lockAfterFunc.RUnlock()
	return calls
}
