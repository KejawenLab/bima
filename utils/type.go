package utils

type Type struct {
	types map[string]string
}

func NewType() Type {
	return Type{
		types: map[string]string{
			"double":   "float64",
			"float":    "float32",
			"int32":    "int32",
			"int64":    "int64",
			"uint32":   "uint32",
			"uint64":   "uint64",
			"sint32":   "int32",
			"sint64":   "int64",
			"fixed32":  "uint32",
			"fixed64":  "uint64",
			"sfixed32": "int32",
			"sfixed64": "int64",
			"bool":     "bool",
			"string":   "string",
			"bytes":    "[]byte",
		},
	}
}

func (t Type) List() map[string]string {
	return t.types
}

func (t Type) Value(key string) string {
	if value, ok := t.types[key]; ok {
		return value
	}

	return ""
}
