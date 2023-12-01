package rusty

type Arc[V any] struct {
	v V
}

func NewArc[V *x, x any](v V) Arc[V] {
	return Arc[V]{
		v: v,
	}
}

func (a Arc[V]) Read() V {
	return a.v
}
