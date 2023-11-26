package xres

type Option[V any] struct {
	value     V
	isPresent bool
}

func (o *Option[V]) IsPresent() bool {
	return o.isPresent
}

func (o *Option[V]) Unwrap() V {
	if !o.isPresent {
		panic("option is empty")
	}

	return o.value
}

func (o *Option[V]) UnwrapOrElse(fn func() V) V {
	if !o.isPresent {
		return fn()
	}

	return o.value
}

func (o *Option[V]) UnwrapOrDefault() V {
	return o.value
}
