package tracelog

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/bj-wangjia/priority_queue"
	"github.com/hashicorp/consul/api"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

var (
	NotSupportType      = errors.New("not support type")
	KvNotFound          = errors.New("kv not found")
	GetConsulKvFailed   = errors.New("get consul kv info failed")
	JsonUnmarshalFailed = errors.New("json unmarshal failed")

	scheduleUnit = priority_queue.New()
	onSched      = sync.Once{}
)

type Ops struct {
	Type     string
	Address  string
	Path     string
	Interval time.Duration
	TryTimes int
	OnChange func(value interface{}, err error) bool
}

type ScheduleUnit struct {
	ops      *Ops
	deadline int64
	t        reflect.Type
}

func (su *ScheduleUnit) Less(i interface{}) bool {
	return su.deadline < i.(*ScheduleUnit).deadline
}

func RunSchedule(ctx context.Context) {
	onSched.Do(func() {
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Run Schedule exit.")
					return
				default:
					if scheduleUnit.Len() <= 0 {
						fmt.Println("no task, sleep 1s")
						time.Sleep(time.Second)
						continue
					}
					unit := scheduleUnit.Top().(*ScheduleUnit)
					now := time.Now().UnixNano()
					if unit.deadline > now {
						time.Sleep(time.Duration(unit.deadline - now))
						continue
					}
					runTask(unit)
					unit.deadline += int64(unit.ops.Interval)
					scheduleUnit.Fix(unit, 0)
				}
			}
		}()
	})
}

func runTask(unit *ScheduleUnit) {
	o := reflect.New(unit.t.Elem())
	i := o.Interface()
	err := getConfig(unit.ops, i)
	if unit.ops.OnChange != nil {
		unit.ops.OnChange(i, err)
	}
}

func GetConfig(ops *Ops, value interface{}) error {
	err := getConfig(ops, value)
	if ops.OnChange != nil {
		ops.OnChange(value, err)
	}
	if err != nil {
		return err
	}
	if ops.Interval > 0 {
		scheduleUnit.Push(&ScheduleUnit{
			ops:      ops,
			deadline: time.Now().UnixNano(),
			t:        reflect.TypeOf(value),
		})
	}
	return nil
}

func getConfig(ops *Ops, value interface{}) error {
	switch ops.Type {
	case "json":
		return getJsonConfig(ops, value)
	case "toml":
		return getTomlConfig(ops, value)
	default:
		return NotSupportType
	}
}

func GetValue(ops *Ops) (*api.KVPair, error) {
	config := api.DefaultConfig()
	config.Address = ops.Address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	kv := client.KV()
	if kv == nil {
		return nil, GetConsulKvFailed
	}
	pair, _, err := kv.Get(ops.Path, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, KvNotFound
	}
	return pair, nil
}

func getTomlConfig(ops *Ops, value interface{}) error {
	pair, err := GetValue(ops)
	if err != nil {
		return err
	}
	if pair == nil {
		return KvNotFound
	}
	if _, err = toml.Decode(*(*string)(unsafe.Pointer(&pair.Value)), value); err != nil {
		return err
	}
	return nil
}

func getJsonConfig(ops *Ops, value interface{}) error {
	pair, err := GetValue(ops)
	if err != nil {
		return err
	}
	if pair == nil {
		return KvNotFound
	}
	err = jsoniter.Unmarshal(pair.Value, value)
	if err != nil {
		return JsonUnmarshalFailed
	}
	return nil
}
