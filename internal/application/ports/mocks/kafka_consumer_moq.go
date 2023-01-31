// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/saltpay/go-kafka-driver"
	"github.com/saltpay/transaction-card-validator/internal/application/ports"
	"sync"
)

// Ensure, that ConsumerMock does implement ports.Consumer.
// If this is not the case, regenerate this file with moq.
var _ ports.Consumer = &ConsumerMock{}

// ConsumerMock is a mock implementation of ports.Consumer.
//
// 	func TestSomethingThatUsesConsumer(t *testing.T) {
//
// 		// make and configure a mocked ports.Consumer
// 		mockedConsumer := &ConsumerMock{
// 			ListenFunc: func(ctx context.Context, processor kafka.Processor, commitStrategy kafka.CommitStrategy, ps kafka.PauseStrategy)  {
// 				panic("mock out the Listen method")
// 			},
// 		}
//
// 		// use mockedConsumer in code that requires ports.Consumer
// 		// and then make assertions.
//
// 	}
type ConsumerMock struct {
	// ListenFunc mocks the Listen method.
	ListenFunc func(ctx context.Context, processor kafka.Processor, commitStrategy kafka.CommitStrategy, ps kafka.PauseStrategy)

	// calls tracks calls to the methods.
	calls struct {
		// Listen holds details about calls to the Listen method.
		Listen []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Processor is the processor argument value.
			Processor kafka.Processor
			// CommitStrategy is the commitStrategy argument value.
			CommitStrategy kafka.CommitStrategy
			// Ps is the ps argument value.
			Ps kafka.PauseStrategy
		}
	}
	lockListen sync.RWMutex
}

// Listen calls ListenFunc.
func (mock *ConsumerMock) Listen(ctx context.Context, processor kafka.Processor, commitStrategy kafka.CommitStrategy, ps kafka.PauseStrategy) {
	if mock.ListenFunc == nil {
		panic("ConsumerMock.ListenFunc: method is nil but Consumer.Listen was just called")
	}
	callInfo := struct {
		Ctx            context.Context
		Processor      kafka.Processor
		CommitStrategy kafka.CommitStrategy
		Ps             kafka.PauseStrategy
	}{
		Ctx:            ctx,
		Processor:      processor,
		CommitStrategy: commitStrategy,
		Ps:             ps,
	}
	mock.lockListen.Lock()
	mock.calls.Listen = append(mock.calls.Listen, callInfo)
	mock.lockListen.Unlock()
	mock.ListenFunc(ctx, processor, commitStrategy, ps)
}

// ListenCalls gets all the calls that were made to Listen.
// Check the length with:
//     len(mockedConsumer.ListenCalls())
func (mock *ConsumerMock) ListenCalls() []struct {
	Ctx            context.Context
	Processor      kafka.Processor
	CommitStrategy kafka.CommitStrategy
	Ps             kafka.PauseStrategy
} {
	var calls []struct {
		Ctx            context.Context
		Processor      kafka.Processor
		CommitStrategy kafka.CommitStrategy
		Ps             kafka.PauseStrategy
	}
	mock.lockListen.RLock()
	calls = mock.calls.Listen
	mock.lockListen.RUnlock()
	return calls
}