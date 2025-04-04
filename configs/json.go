package configs

import "github.com/invopop/jsonschema"

func GenerateSchema[T any]() interface{} {
	ref := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var obj T
	schema := ref.Reflect(obj)
	return schema
}
