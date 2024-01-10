package data

import v1metrics "go.opentelemetry.io/proto/otlp/metrics/v1"

type ResourceMetrics struct {
	raw *v1metrics.ResourceMetrics
}

func NewResourceMetrics(resource Resource, scopeMetrics ...ScopeMetrics) ResourceMetrics {
	rm := ResourceMetrics{
		raw: &v1metrics.ResourceMetrics{
			Resource:     resource.raw,
			ScopeMetrics: make([]*v1metrics.ScopeMetrics, 0, len(scopeMetrics)),
		},
	}

	for i := 0; i < len(scopeMetrics); i++ {
		rm.raw.ScopeMetrics = append(rm.raw.ScopeMetrics, scopeMetrics[i].raw)
	}

	return rm
}
