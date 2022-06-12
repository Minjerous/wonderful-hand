// cache_mq 用于在不同实例中同步共享的缓存，比如验证码等缓存

package helper

import (
	"context"
	"encoding/json"
	"fmt"
	common "github.com/igxnon/cachepool/pkg/cache"
	"github.com/streadway/amqp"
	"time"
)

const exchangeName = "exchange.__cache_sync__"

func runSyncFromMQ(ctx context.Context, cache common.ICache, ch *amqp.Channel, name string) error {
	err := ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		name,
		true,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		name,
		fmt.Sprintf("%s-key", name),
		exchangeName,
		true, nil,
	)

	if err != nil {
		return err
	}

	msg, err := ch.Consume(
		name,
		fmt.Sprintf("%s-consumer", name),
		true,
		true,
		false,
		true,
		nil,
	)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case m, ok := <-msg:
			if !ok {
				return nil
			}
			opt, key, value, exp, err := decode(m.Body)
			if err != nil {
				// log.Printf("Warning comsumer %s message %s err %v\n",
				// 	m.ConsumerTag, m.MessageId, err)
				continue
			}
			if opt {
				cache.Set(key, value, exp)
				continue
			}
			cache.Delete(key)
		}
	}
}

type data struct {
	Opt   bool          `json:"opt"` // true -> add, false -> delete
	Key   string        `json:"key"`
	Value any           `json:"value,omitempty"`
	Exp   time.Duration `json:"Exp,omitempty"`
}

func decode(b []byte) (opt bool, key string, value any, exp time.Duration, err error) {
	d := data{}
	err = json.Unmarshal(b, &d)
	return d.Opt, d.Key, d.Value, d.Exp, err
}

// Publish 将缓存同步到所有实例里
func Publish(ch *amqp.Channel, key string, value any, d time.Duration) error {
	data := data{
		Opt:   true,
		Key:   key,
		Value: value,
		Exp:   d,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ch.Publish(exchangeName, "", false, false, amqp.Publishing{
		Timestamp:    time.Now(),
		MessageId:    key,
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         b,
	})
}

func PublishDel(ch *amqp.Channel, key string) error {
	data := data{
		Opt: false,
		Key: key,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ch.Publish(exchangeName, "", false, false, amqp.Publishing{
		Timestamp:    time.Now(),
		MessageId:    key,
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         b,
	})
}
