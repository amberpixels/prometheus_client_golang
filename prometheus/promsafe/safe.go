// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package promsafe provides safe labeling - strongly typed labels in prometheus metrics.
// Enjoy promsafe as you wish!
package promsafe

import (
	"github.com/prometheus/client_golang/prometheus"
)

// NewCounterVec creates a new CounterVec with type-safe labels.
func NewCounterVec[T LabelsProviderMarker](opts prometheus.CounterOpts) *CounterVec[T] {
	emptyLabels := NewEmptyLabels[T]()
	inner := prometheus.NewCounterVec(opts, extractLabelNames(emptyLabels))

	return &CounterVec[T]{inner: inner}
}

// CounterVec is a wrapper around prometheus.CounterVec that allows type-safe labels.
type CounterVec[T LabelsProviderMarker] struct {
	inner *prometheus.CounterVec
}

// GetMetricWithLabelValues covers prometheus.CounterVec.GetMetricWithLabelValues
// Deprecated: Use GetMetricWith() instead. We can't provide a []string safe implementation in promsafe
func (c *CounterVec[T]) GetMetricWithLabelValues(_ ...string) (prometheus.Counter, error) {
	panic("There can't be a SAFE GetMetricWithLabelValues(). Use GetMetricWith() instead")
}

// GetMetricWith behaves like prometheus.CounterVec.GetMetricWith but with type-safe labels.
func (c *CounterVec[T]) GetMetricWith(labels T) (prometheus.Counter, error) {
	return c.inner.GetMetricWith(extractLabelsWithValues(labels))
}

// WithLabelValues covers like prometheus.CounterVec.WithLabelValues.
// Deprecated: Use With() instead. We can't provide a []string safe implementation in promsafe
func (c *CounterVec[T]) WithLabelValues(_ ...string) prometheus.Counter {
	panic("There can't be a SAFE WithLabelValues(). Use With() instead")
}

// With behaves like prometheus.CounterVec.With but with type-safe labels.
func (c *CounterVec[T]) With(labels T) prometheus.Counter {
	return c.inner.With(extractLabelsWithValues(labels))
}

// CurryWith behaves like prometheus.CounterVec.CurryWith but with type-safe labels.
// It still returns a CounterVec, but it's inner prometheus.CounterVec is curried.
func (c *CounterVec[T]) CurryWith(labels T) (*CounterVec[T], error) {
	curriedInner, err := c.inner.CurryWith(extractLabelsWithValues(labels))
	if err != nil {
		return nil, err
	}
	c.inner = curriedInner
	return c, nil
}

// MustCurryWith behaves like prometheus.CounterVec.MustCurryWith but with type-safe labels.
// It still returns a CounterVec, but it's inner prometheus.CounterVec is curried.
func (c *CounterVec[T]) MustCurryWith(labels T) *CounterVec[T] {
	c.inner = c.inner.MustCurryWith(extractLabelsWithValues(labels))
	return c
}

// Unsafe returns the underlying prometheus.CounterVec
// it's used to call any other method of prometheus.CounterVec that doesn't require type-safe labels
func (c *CounterVec[T]) Unsafe() *prometheus.CounterVec {
	return c.inner
}

// NewCounter simply creates a new prometheus.Counter.
// As it doesn't have any labels, it's already type-safe.
// We keep this method just for consistency and interface fulfillment.
func NewCounter(opts prometheus.CounterOpts) prometheus.Counter {
	return prometheus.NewCounter(opts)
}

// NewCounterFunc wraps a new prometheus.CounterFunc.
// As it doesn't have any labels, it's already type-safe.
// We keep this method just for consistency and interface fulfillment.
func NewCounterFunc(opts prometheus.CounterOpts, function func() float64) prometheus.CounterFunc {
	return prometheus.NewCounterFunc(opts, function)
}
