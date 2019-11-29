package json

import (
	"encoding/json"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/jeremywohl/flatten"
)

type serializer struct {
	TimestampUnits time.Duration
}

func NewSerializer(timestampUnits time.Duration) (*serializer, error) {
	s := &serializer{
		TimestampUnits: truncateDuration(timestampUnits),
	}
	return s, nil
}

func (s *serializer) Serialize(metric telegraf.Metric) ([]byte, error) {
	m := s.createObject(metric)
	serialized, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	serialized = append(serialized, '\n')

	return serialized, nil
}

func (s *serializer) SerializeBatch(metrics []telegraf.Metric) ([]byte, error) {
	objects := make([]interface{}, 0, len(metrics))
	for _, metric := range metrics {
		m := s.createObject(metric)
		objects = append(objects, m)
	}

	//obj := map[string]interface{}{
	//	"metrics": objects,
	//}

	serialized, err := json.Marshal(objects)
	if err != nil {
		return []byte{}, err
	}
	return serialized, nil
}

func (s *serializer) createObject(metric telegraf.Metric) map[string]interface{} {
	m := make(map[string]interface{}, 4)

	//m["tags"] = metric.Tags()

	for k, v := range metric.Tags() {
		if k == "resource" {
			m["resource"] = v
		} else {
	//		m["tag."+k] = v
		}
	}

	//m["fields"] = metric.Fields()

        for k, v := range metric.Fields() {
               if k == "value" {
		      m[metric.Name()] = v
	       }
	
	//m["name"] = metric.Name()
	
	m["timestamp"] = metric.Time().UnixNano() / int64(s.TimestampUnits)
	
	//flat, err := flatten.Flatten(m, "", flatten.DotStyle)
	//if err != nil {
	//	m["err"] = err
	//	return m
	//}
	
	return m
}
func truncateDuration(units time.Duration) time.Duration {
	// Default precision is 1s
	if units <= 0 {
		return time.Second
	}

	// Search for the power of ten less than the duration
	d := time.Nanosecond
	for {
		if d*10 > units {
			return d
		}
		d = d * 10
	}
}
