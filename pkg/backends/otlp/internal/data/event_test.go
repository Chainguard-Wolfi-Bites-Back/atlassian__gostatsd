package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1common "go.opentelemetry.io/proto/otlp/common/v1"
	v1log "go.opentelemetry.io/proto/otlp/logs/v1"

	"github.com/atlassian/gostatsd"
)

func TestTransformToLog(t *testing.T) {
	tests := []struct {
		name              string
		gostatsdEvent     *gostatsd.Event
		titleAttrKey      string
		categoryAttrKey   string
		propertiesAttrKey string
		want              *v1log.LogRecord
	}{
		{
			name: "should convert event log record with default attributes fields",
			gostatsdEvent: &gostatsd.Event{
				Title:        "title",
				Text:         "text",
				DateHappened: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				Tags:         gostatsd.Tags{"tag1:1", "tag2:2"},
				Source:       "127.0.0.1",
				Priority:     gostatsd.PriNormal,
				AlertType:    gostatsd.AlertError,
			},
			want: &v1log.LogRecord{
				TimeUnixNano: uint64(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()),
				Attributes: []*v1common.KeyValue{
					{
						Key:   "title",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "title"}},
					},
					{
						Key:   "tag1",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "1"}},
					},
					{
						Key:   "tag2",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "2"}},
					},
					{
						Key:   "host",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "127.0.0.1"}},
					},
					{
						Key:   "priority",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: gostatsd.PriNormal.String()}},
					},
					{
						Key:   "alert_type",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: gostatsd.AlertError.String()}},
					},
					{
						Key: "properties",
						Value: &v1common.AnyValue{
							Value: &v1common.AnyValue_KvlistValue{
								KvlistValue: &v1common.KeyValueList{
									Values: []*v1common.KeyValue{
										{
											Key:   "text",
											Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "text"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "should convert event log record with custom attributes fields",
			gostatsdEvent: &gostatsd.Event{
				Title:        "title",
				Text:         "text",
				DateHappened: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				Tags:         gostatsd.Tags{"tag1:1", "tag2:2"},
				Source:       "127.0.0.1",
				Priority:     gostatsd.PriNormal,
				AlertType:    gostatsd.AlertError,
			},
			titleAttrKey:      "com.atlassian.title",
			propertiesAttrKey: "com.atlassian.properties",
			want: &v1log.LogRecord{
				TimeUnixNano: uint64(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()),
				Attributes: []*v1common.KeyValue{
					{
						Key:   "com.atlassian.title",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "title"}},
					},
					{
						Key:   "tag1",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "1"}},
					},
					{
						Key:   "tag2",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "2"}},
					},
					{
						Key:   "host",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "127.0.0.1"}},
					},
					{
						Key:   "priority",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: gostatsd.PriNormal.String()}},
					},
					{
						Key:   "alert_type",
						Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: gostatsd.AlertError.String()}},
					},
					{
						Key: "com.atlassian.properties",
						Value: &v1common.AnyValue{
							Value: &v1common.AnyValue_KvlistValue{
								KvlistValue: &v1common.KeyValueList{
									Values: []*v1common.KeyValue{
										{
											Key:   "text",
											Value: &v1common.AnyValue{Value: &v1common.AnyValue_StringValue{StringValue: "text"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewOtlpEvent(tt.gostatsdEvent)
			assert.NoError(t, err)

			if tt.titleAttrKey != "" {
				s.titleAttrKey = tt.titleAttrKey
			}
			if tt.propertiesAttrKey != "" {
				s.propertiesAttrKey = tt.propertiesAttrKey
			}
			record := s.TransformToLog()
			assert.Equal(t, tt.want.TimeUnixNano, record.TimeUnixNano)
			assert.Equal(t, len(tt.want.Attributes), len(record.Attributes))
			for _, kv := range tt.want.Attributes {
				found := false
				for _, attr := range record.Attributes {
					if kv.Key == attr.Key {
						assert.Equal(t, kv, attr)
						found = true
						break
					}
				}
				if !found {
					t.Errorf("attribute %s not found", kv.Key)
				}
			}
		})
	}
}
