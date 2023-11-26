package xres

type Option[V any] interface {
	IsPresent() bool
	Unwrap() V
	UnwrapOrElse(fn func() V) V
	UnwrapOrDefault() V
}

func Some[V any](value V) Option[V] {
	return &option[V]{
		value:     value,
		isPresent: true,
	}
}

func None[V any]() Option[V] {
	return &option[V]{
		isPresent: false,
	}
}

type option[V any] struct {
	value     V
	isPresent bool
}

func (o *option[V]) IsPresent() bool {
	return o.isPresent
}

func (o *option[V]) Unwrap() V {
	if !o.isPresent {
		panic("option is empty")
	}

	return o.value
}

func (o *option[V]) UnwrapOrElse(fn func() V) V {
	if !o.isPresent {
		return fn()
	}

	return o.value
}

func (o *option[V]) UnwrapOrDefault() V {
	return o.value
}
