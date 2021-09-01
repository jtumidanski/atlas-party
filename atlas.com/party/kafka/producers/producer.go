package producers

import (
	"atlas-party/retry"
	"atlas-party/topic"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func CreateKey(key int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, uint32(key))
	return b
}

func ProduceEvent(l logrus.FieldLogger, topicToken string) func(key []byte, event interface{}) {
	name := topic.GetRegistry().Get(l, topicToken)
	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
		Topic:        name,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}

	return func(key []byte, event interface{}) {
		r, err := json.Marshal(event)
		l.WithField("message", string(r)).Debugf("Writing message to topic %s.", name)
		if err != nil {
			l.WithError(err).Fatalf("Unable to marshall event for topic %s.", name)
		}

		writeMessage := func(attempt int) (bool, error) {
			err = w.WriteMessages(context.Background(), kafka.Message{
				Key:   key,
				Value: r,
			})
			if err != nil {
				l.Warnf("Unable to emit event on topic %s, will retry.", name)
				return true, err
			}
			return false, err
		}

		err = retry.Try(writeMessage, 10)
		if err != nil {
			l.WithError(err).Fatalf("Unable to emit event on topic %s.", name)
		}
	}
}
