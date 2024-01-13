package gost_test

import (
	"net/http"
	"testing"

	"github.com/egorgasay/gost"
)

func will() gost.ErrX {
	return gost.NewErrX(_notFound)
}

func wont() gost.ErrX {
	return gost.Nil
}

const (
	_unknown    = 0
	_notFound   = 404
	_order      = 1000
	_repository = 134
)

func x() gost.ErrX {
	err := will().Extend(_order).Extend(_order, "test")

	_ = err

	return gost.Nil
}

func repository_DeleteOrder(id string) gost.ErrX {

	// let's throw an error

	return gost.NewErrX(_notFound, "not found").
		Extend(_order, "order") // we can extend it!
}

func usecase_DeleteOrder(id string) gost.ErrX {
	if err := repository_DeleteOrder(id); err.IsErr() {
		if err.BaseCode() == _unknown {
			return err.Extend(_repository, "problem in repository") // another extend
		}
		return err
	}

	return gost.Nil
}

func http_FindOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if err := usecase_DeleteOrder(id); err.IsErr() {
		http_handleErrX(w, err)
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}

	return
}

func http_handleErrX(w http.ResponseWriter, err gost.ErrX) {
	switch err.BaseCode() {
	case _notFound:
		if err.CmpExt(_order) {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(404)
		}
	case _repository:
		w.WriteHeader(503)
	default:
		w.WriteHeader(500)
	}

	w.Write([]byte(err.Error()))
}

func TestErrXError(t *testing.T) {
	want := `404: not found, 1000: , 134: test`

	got := gost.NewErrX(_notFound, "not found").Extend(_order).Extend(134, "test")

	if !got.IsErr() {
		t.Fatalf("want error, got no %v", got)
	}

	if !got.CmpBase(_notFound) {
		t.Fatalf("wanted got.baseCode: %d, got: %d", _notFound, got.BaseCode())
	}

	if !got.CmpExt(_order) {
		t.Fatalf("wanted got.extCode: %d, got: %d", _order, got.ExtCodes()[0])
	}

	if !got.CmpExt(134) {
		t.Fatalf("wanted got.extCode: %d, got: %d", 134, got.ExtCodes()[1])
	}

	if got.CmpBase(14) {
		t.Fatalf("wanted to fail, got: %v", got)
	}

	if got.Error() != want {
		t.Fatalf("unexpected error: [%s] != [%s]", got, want)
	}

	t.Logf("got: %s", got)
}
