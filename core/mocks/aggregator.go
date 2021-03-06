// Code generated by mockery v1.0.0
package coremock

import "context"
import "github.com/yandex/pandora/core"
import "github.com/stretchr/testify/mock"

// Aggregator is an autogenerated mock type for the Aggregator type
type Aggregator struct {
	mock.Mock
}

// Report provides a mock function with given fields: _a0
func (_m *Aggregator) Report(_a0 core.Sample) {
	_m.Called(_a0)
}

// Run provides a mock function with given fields: _a0
func (_m *Aggregator) Run(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
