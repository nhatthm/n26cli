package cli

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/imdario/mergo"
	"github.com/nhatthm/n26cli/internal/service"
)

var emptyUUID uuid.UUID

type globalConfig struct {
	Verbose    bool
	Debug      bool
	ConfigFile string
}

type transformersCallback struct {
	call func(typeOf reflect.Type) func(dst, src reflect.Value) error
}

func (c *transformersCallback) Transformer(typeOf reflect.Type) func(dst, src reflect.Value) error {
	return c.call(typeOf)
}

func newTransformersCallback(call func(typeOf reflect.Type) func(dst, src reflect.Value) error) *transformersCallback {
	return &transformersCallback{call: call}
}

func uuidTransformer() *transformersCallback {
	return newTransformersCallback(func(typeOf reflect.Type) func(dst reflect.Value, src reflect.Value) error {
		if typeOf != reflect.TypeOf(uuid.UUID{}) {
			return nil
		}

		return func(dst reflect.Value, src reflect.Value) error {
			if dst.CanSet() {
				if val := dst.Interface().(uuid.UUID); val == emptyUUID { // nolint: errcheck
					dst.Set(src)
				}
			}

			return nil
		}
	})
}

func mergeConfig(dest *service.Config, srcs ...service.Config) error {
	for _, src := range srcs {
		if err := mergo.Merge(dest, src, mergo.WithTransformers(uuidTransformer())); err != nil {
			return err
		}
	}

	return nil
}
