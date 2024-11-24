package promauto

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promsafe"
)

// NewCounterVec behaves as promauto.NewCounterVec but with type-safe labels
func NewCounterVec[T promsafe.LabelsProviderMarker](opts prometheus.CounterOpts) *promsafe.CounterVec[T] {
	//_ = promauto.NewCounterVec // keeping for reference

	c := promsafe.NewCounterVec[T](opts)
	if prometheus.DefaultRegisterer != nil {
		prometheus.DefaultRegisterer.MustRegister(c.Unsafe())
	}
	return c
}

// Factory is a promauto-like factory that allows type-safe labels.
type Factory[T promsafe.LabelsProviderMarker] struct {
	r prometheus.Registerer
}

// With behaves same as promauto.With but with type-safe labels
func With[T promsafe.LabelsProviderMarker](r prometheus.Registerer) Factory[T] {
	return Factory[T]{r: r}
}

// NewCounterVec behaves like promauto.NewCounterVec but with type-safe labels
func (f Factory[T]) NewCounterVec(opts prometheus.CounterOpts) *promsafe.CounterVec[T] {
	c := NewCounterVec[T](opts)
	if f.r != nil {
		f.r.MustRegister(c.Unsafe())
	}
	return c
}

// NewCounter wraps promauto.NewCounter.
// As it doesn't require any labels, it's already type-safe, and we keep it for consistency.
func (f Factory[T]) NewCounter(opts prometheus.CounterOpts) prometheus.Counter {
	return promauto.With(f.r).NewCounter(opts)
}

// NewCounterFunc wraps promauto.NewCounterFunc.
// As it doesn't require any labels, it's already type-safe, and we keep it for consistency.
func (f Factory[T]) NewCounterFunc(opts prometheus.CounterOpts, function func() float64) prometheus.CounterFunc {
	return promauto.With(f.r).NewCounterFunc(opts, function)
}
