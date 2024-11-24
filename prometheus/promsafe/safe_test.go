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

package promsafe_test

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promsafe"
)

// These are Examples that can be treated as basic smoke tests

func ExampleNewCounterVecT_multiple_labels_manual() {
	// Manually registering with multiple labels

	type MyCounterLabels struct {
		promsafe.StructLabelProvider
		EventType string
		Success   bool
		Position  uint8 // yes, it's a number, but be careful with high-cardinality labels

		ShouldNotBeUsed string `promsafe:"-"`
	}

	c := promsafe.NewCounterVec[MyCounterLabels](prometheus.CounterOpts{
		Name: "items_counted_detailed",
	})

	// Manually register the counter
	if err := prometheus.Register(c.Unsafe()); err != nil {
		log.Fatal("could not register: ", err.Error())
	}

	// and now, because of generics we can call Inc() with filled struct of labels:
	counter := c.With(MyCounterLabels{
		EventType: "reservation", Success: true, Position: 1,
	})
	counter.Inc()

	// Output:
}

// FastMyLabels is a struct that will have a custom method that converts to prometheus.Labels
type FastMyLabels struct {
	promsafe.StructLabelProvider
	EventType string
	Source    string
}

// ToPrometheusLabels does a superfast conversion to labels. So no reflection is involved.
func (f FastMyLabels) ToPrometheusLabels() prometheus.Labels {
	return prometheus.Labels{"event_type": f.EventType, "source": f.Source}
}

// ToLabelNames does a superfast label names list. So no reflection is involved.
func (f FastMyLabels) ToLabelNames() []string {
	return []string{"event_type", "source"}
}

func ExampleNewCounterVecT_fast_safe_labels_provider() {
	// Note: fast labels provider has a drawback: they can't be declared as inline structs
	//       as we need methods

	c := promsafe.NewCounterVec[FastMyLabels](prometheus.CounterOpts{
		Name: "items_counted_detailed_fast",
	})

	// Manually register the counter
	if err := prometheus.Register(c.Unsafe()); err != nil {
		log.Fatal("could not register: ", err.Error())
	}

	counter := c.With(FastMyLabels{
		EventType: "reservation", Source: "source1",
	})
	counter.Inc()

	// Output:
}
