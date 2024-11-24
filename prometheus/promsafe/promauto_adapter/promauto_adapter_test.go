package promauto_adapter_test

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promsafe"
	promauto "github.com/prometheus/client_golang/prometheus/promsafe/promauto_adapter"
)

func ExampleNewCounterVec_promauto_adapted() {
	// Examples on how to migrate from promauto to promsafe gently
	//
	// Before:
	//   import "github.com/prometheus/client_golang/prometheus/promauto"
	// func main() {
	//   myReg := prometheus.NewRegistry()
	//   counterOpts := prometheus.CounterOpts{Name:"..."}
	//   promauto.With(myReg).NewCounterVec(counterOpts, []string{"event_type", "source"})
	// }
	//
	// After:
	//
	//   import (
	//     promauto "github.com/prometheus/client_golang/prometheus/promsafe/promauto_adapter"
	//   )
	//   ...

	myReg := prometheus.NewRegistry()
	counterOpts := prometheus.CounterOpts{
		Name: "items_counted_detailed_auto",
	}

	type MyLabels struct {
		promsafe.StructLabelProvider
		EventType string
		Source    string
	}
	c := promauto.With[MyLabels](myReg).NewCounterVec(counterOpts)

	c.With(MyLabels{
		EventType: "reservation", Source: "source1",
	}).Inc()

	// Output:
}
