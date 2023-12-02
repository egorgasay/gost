package gost

type Option[V any] struct {
	value *V
}

func (Option[V]) Some(value V) Option[V] {
	return Option[V]{
		value: &value,
	}
}

func (Option[V]) None() Option[V] {
	return Option[V]{}
}

func Some[V any](value V) Option[V] {
	return Option[V]{
		value: &value,
	}
}

func None[V any]() Option[V] {
	return Option[V]{}
}

func (o Option[V]) IsSome() bool {
	return o.value != nil
}

func (o Option[V]) IsNone() bool {
	return o.value == nil
}

func (o Option[V]) Unwrap() V {
	if o.value == nil {
		panic("option is empty")
	}

	return *o.value
}

func (o Option[V]) UnwrapOrElse(fn func() V) V {
	if o.value == nil {
		return fn()
	}

	return *o.value
}

func (o Option[V]) UnwrapOrDefault() V {
	if o.value == nil {
		var v V
		return v
	}

	return *o.value
}
