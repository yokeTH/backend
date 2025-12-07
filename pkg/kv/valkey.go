package kv

import (
	"context"
	"time"

	"github.com/valkey-io/valkey-go"
)

type Valkey struct {
	client valkey.Client
}

func NewValkeyClient(connectionString string) (*Valkey, error) {
	opt, err := valkey.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	client, err := valkey.NewClient(opt)
	if err != nil {
		return nil, err
	}

	return &Valkey{client: client}, nil
}

func (v *Valkey) Close() {
	v.client.Close()
}

// Set executes: SET key val [EX ttl]
func (v *Valkey) Set(ctx context.Context, key, val string, ttl time.Duration) error {
	cmd := v.client.B().Set().Key(key).Value(val)
	if ttl != 0 {
		cmd.Ex(ttl)
	}
	return v.client.Do(ctx, cmd.Build()).Error()
}

// Get executes: GET key
// If the key is missing, it returns ("", nil).
func (v *Valkey) Get(ctx context.Context, key string) (string, error) {
	msg := v.client.Do(ctx, v.client.B().Get().Key(key).Build())
	result, err := msg.ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			// cache miss, but not an error at the API level
			return "", nil
		}
		return "", err
	}
	return result, nil
}

// Del executes: DEL key and returns the deleted count.
func (v *Valkey) Del(ctx context.Context, key string) (int, error) {
	n, err := v.client.Do(ctx, v.client.B().Del().Key(key).Build()).AsInt64()
	if err != nil {
		return 0, err
	}
	if n < 0 {
		n = 0
	}
	return int(n), nil
}

// SetNX executes: SET key val NX [EX ttl]
// Returns true if the key was set, false if it already existed.
func (v *Valkey) SetNX(ctx context.Context, key, val string, ttl time.Duration) (bool, error) {
	cmd := v.client.B().Set().Key(key).Value(val).Nx()
	if ttl != 0 {
		cmd.Ex(ttl)
	}

	err := v.client.Do(ctx, cmd.Build()).Error()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Incr executes: INCR key
func (v *Valkey) Incr(ctx context.Context, key string) (int64, error) {
	result, err := v.client.Do(ctx, v.client.B().Incr().Key(key).Build()).AsInt64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Expire executes: EXPIRE key seconds
// ttl == 0 is treated as immediate expiration (seconds 0).
func (v *Valkey) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	secs := max(int64(ttl/time.Second), 0)

	n, err := v.client.Do(
		ctx,
		v.client.B().Expire().Key(key).Seconds(secs).Build(),
	).AsInt64()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

// Exists executes: EXISTS key
func (v *Valkey) Exists(ctx context.Context, key string) (bool, error) {
	n, err := v.client.Do(
		ctx,
		v.client.B().Exists().Key(key).Build(),
	).AsInt64()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

var casDelScript = valkey.NewLuaScript(`
if redis.call('GET', KEYS[1]) == ARGV[1] then
  return redis.call('DEL', KEYS[1])
else
  return 0
end
`)

// CasDel deletes key only if its value matches expected, returning the number
// of keys removed (0 or 1).
func (v *Valkey) CasDel(ctx context.Context, key, expected string) (int, error) {
	n, err := casDelScript.
		Exec(ctx, v.client, []string{key}, []string{expected}).
		AsInt64()
	if err != nil {
		return 0, err
	}
	if n < 0 {
		n = 0
	}
	return int(n), nil
}
