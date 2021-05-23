package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/google/uuid"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
)

/*

Interface RocketMessage represents one rocket message. It encodes our custom message to be sent to
rocketmq broker, and decodes the message back to the custom message.

This interface is to provide one unified way for encoding and decoding. For performance concern, the
real implementation might use other approach instead of json marshal to encode/decode.

Because int64/int can't support json.Unmarshal interface{}, there're int64/int setter/getter funcs for int(s).

For string, it is also supported. More types might be supported later.

For map values, ideally we need lock for each map, but considering the real usage scenario, no need to consider
concurrent cases. This lock might be needed later.

Example Usage:

Producer part:

```
	msg, err := rocketmq.NewMessage().SetString("order_id", orderID).Encode(constants.RocketMQTagOrderPaid)
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
	err = consume(ctx, rocketmsg.GetString("order_id"))

	// Set the message consumed in db.
	err := m.rocketmq.SetMessageConsumed(rocketmsg.ID())
```
*/

type RocketMessage interface {
	// return the message uuid.
	ID() string
	// encode the message to be ready for sent to rocketmq server
	Encode(tag constants.RocketMQTag) (*primitive.Message, error)

	// set message values, such as order id, customer id, etc.
	Set(key string, val interface{}) RocketMessage
	SetInt(key string, val int) RocketMessage
	SetInt64(key string, val int64) RocketMessage
	SetString(key string, val string) RocketMessage

	// get the message values, such as order id, customer id, etc
	Get(key string) interface{}
	GetInt(key string) int
	GetInt64(key string) int64
	GetString(key string) string
}

type rocketMessage struct {
	// id is uuid
	UUID string `json:"uuid"`
	// Values is to store non-int64 values
	Values map[string]interface{} `json:"values"`
	// ValuesInt is to store int values.
	ValuesInt map[string]int `json:"values_int"`
	// ValuesInt64 is to store int64 values.
	ValuesInt64 map[string]int64 `json:"values_int64"`
	// ValuesString is to store string values.
	ValuesString map[string]string `json:"values_string"`
}

func (r *rocketMessage) Encode(tag constants.RocketMQTag) (*primitive.Message, error) {
	// init id
	r.UUID = uuid.New().String()
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	message := primitive.NewMessage(constants.RocketMQTopic, bytes)
	message.WithTag(string(tag))
	// @todo since keys has been provided, maybe message uuid is not needed.
	message.WithKeys([]string{r.UUID})
	return message, nil
}

func (r *rocketMessage) ID() string {
	return r.UUID
}

func (r *rocketMessage) Get(key string) interface{} {
	if r.Values == nil {
		return nil
	}
	val, ok := r.Values[key]
	if !ok {
		return nil
	}
	return val
}

func (r *rocketMessage) GetInt(key string) int {
	msg := fmt.Sprintf("rocket message int key: %s not existing", key)
	if r.ValuesInt == nil {
		panic(msg)
	}
	val, ok := r.ValuesInt[key]
	if !ok {
		panic(msg)
	}
	return val
}

func (r *rocketMessage) GetString(key string) string {
	msg := fmt.Sprintf("rocket message string key: %s not existing", key)
	if r.ValuesString == nil {
		panic(msg)
	}
	val, ok := r.ValuesString[key]
	if !ok {
		panic(msg)
	}
	return val
}

func (r *rocketMessage) GetInt64(key string) int64 {
	msg := fmt.Sprintf("rocket message int64 key: %s not existing", key)
	if r.ValuesInt64 == nil {
		panic(msg)
	}
	val, ok := r.ValuesInt64[key]
	if !ok {
		panic(msg)
	}
	return val
}

func (r *rocketMessage) Set(key string, val interface{}) RocketMessage {
	// Ideally there should be a lock. Condering the real
	// usage, lock might not be necessary.
	if r.Values == nil {
		r.Values = map[string]interface{}{
			key: val,
		}
	} else {
		r.Values[key] = val
	}
	return r
}

func (r *rocketMessage) SetString(key string, val string) RocketMessage {
	// Ideally there should be a lock. Condering the real
	// usage, lock might not be necessary.
	if r.Values == nil {
		r.ValuesString = map[string]string{
			key: val,
		}
	} else {
		r.ValuesString[key] = val
	}
	return r
}

func (r *rocketMessage) SetInt(key string, val int) RocketMessage {
	if r.ValuesInt == nil {
		r.ValuesInt = map[string]int{
			key: val,
		}
	} else {
		r.ValuesInt[key] = val
	}
	return r
}

func (r *rocketMessage) SetInt64(key string, val int64) RocketMessage {
	if r.ValuesInt64 == nil {
		r.ValuesInt64 = map[string]int64{
			key: val,
		}
	} else {
		r.ValuesInt64[key] = val
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
