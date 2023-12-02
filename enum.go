package gost

type Enum[V any] Variant[V]

func (e Enum[V]) Enum() V {
	return e.Value
}

func (e Enum[V]) New(Value V) Enum[V] {
	e.Value = Value
	return e
}

type Variant[V any] struct {
	Value V
}

func (v Variant[V]) Enum() V {
	return v.Value
}

func (v Variant[V]) New(Value V) Enum[V] {
	v.Value = Value
	return Enum[V](v)
}

func NewVariant[V any](Value V) Variant[V] {
	return Variant[V]{
		Value: Value,
	}
}
