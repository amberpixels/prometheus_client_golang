package promauto

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promsafe"
)

func ExampleNewCounterVecT_promauto_migrated() {
	// Examples on how to migrate from promauto to promsafe
	// When promauto was using a custom factory with custom registry

	myReg := prometheus.NewRegistry()

	counterOpts := prometheus.CounterOpts{
		Name: "items_counted_detailed_auto",
	}

	// Old unsafe code
	// promauto.With(myReg).NewCounterVec(counterOpts, []string{"event_type", "source"})
	// becomes:

	type MyLabels struct {
		promsafe.StructLabelProvider
		EventType string
		Source    string
	}
	c := With[MyLabels](myReg).NewCounterVec(counterOpts)

	c.With(MyLabels{
		EventType: "reservation", Source: "source1",
	}).Inc()

	// Output:
}

func ExampleNewCounterVecT_pointer_to_labels_promauto() {
	// It's possible to use pointer to labels struct
	myReg := prometheus.NewRegistry()

	counterOpts := prometheus.CounterOpts{
		Name: "items_counted_detailed_ptr",
	}

	type MyLabels struct {
		promsafe.StructLabelProvider
		EventType string
		Source    string
	}
	c := With[*MyLabels](myReg).NewCounterVec(counterOpts)

	c.With(&MyLabels{
		EventType: "reservation", Source: "source1",
	}).Inc()

	// Output:
}
