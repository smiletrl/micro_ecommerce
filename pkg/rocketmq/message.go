package rocketmq

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/google/uuid"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
)

/*
Example Usage:

Producer part:

```
	msg, err := rocketmq.NewMessage().Set("order_id", orderID).Encode(constants.RocketMQTagOrderPaid)
	if err != nil {
		return err
	}
	_, err = m.producer.SendSync(ctx, msg)
```

Consumer part:

```
	rocketmsg, err := rocketmq.DecodeMessage(msgs[0].Body)
	if err != nil {
		return consumer.Commit, err
	}

	// See if this message has been consumed already.
	has, err := m.rocketmq.HasMessageConsumed(rocketmsg.ID())
	if err != nil {
		m.logger.Errorw("rocketmq order message consumed", string(rocketmsg.ID()))

		return consumer.Commit, err
	}

	// If it has been consumed already, skip this message.
	if has {
		return consumer.ConsumeSuccess, nil
	}

	// Real consume happens here.
	// consume is a custom consume function for this message
	err = consume(ctx, rocketmsg.Get("order_id").(string))

	// Set the message consumed in db.
	err := m.rocketmq.SetMessageConsumed(rocketmsg.ID())
```
*/

type RocketMessage interface {
	// set message values, such as order id, customer id, etc.
	Set(key string, val interface{}) RocketMessage
	// encode the message to be ready for sent to rocketmq server
	Encode(tag constants.RocketMQTag) (*primitive.Message, error)

	// return the message uuid.
	ID() string
	// get the message values, such as order id, customer id, etc
	Get(key string) interface{}
}

type rocketMessage struct {
	// id is uuid
	id     string
	values map[string]interface{}
}

func (r *rocketMessage) Encode(tag constants.RocketMQTag) (*primitive.Message, error) {
	// init id
	r.id = uuid.New().String()
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	message := primitive.NewMessage(constants.RocketMQTopic, bytes)
	message.WithTag(string(tag))
	return message, nil
}

func (r *rocketMessage) ID() string {
	return r.id
}

func (r *rocketMessage) Get(key string) interface{} {
	if r.values == nil {
		return nil
	}
	val, ok := r.values[key]
	if !ok {
		return nil
	}
	return val
}

func (r *rocketMessage) Set(key string, val interface{}) RocketMessage {
	// Ideally there should be a lock. Condering the real
	// usage, lock might not be necessary.
	if r.values == nil {
		r.values = map[string]interface{}{
			key: val,
		}
	} else {
		r.values[key] = val
	}
	return r
}

// ---- util ---- //

type MessageOpt func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)

func DecodeMessage(bytes []byte) (RocketMessage, error) {
	m := rocketMessage{}
	err := json.Unmarshal(bytes, &m)
	return &m, err
}

func NewMessage() RocketMessage {
	return &rocketMessage{}
}
