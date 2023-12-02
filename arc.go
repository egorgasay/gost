package gost

type Arc[V *X, X any] struct {
	v V
}

func NewArc[V *X, X any](v V) Arc[V, X] {
	return Arc[V, X]{
		v: v,
	}
}

func (a Arc[V, X]) Read() X {
	return *a.v
}

func (a Arc[V, X]) Ref() V {
	x := *a.v
	return &x
}
