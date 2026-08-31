package session

import (
	"context"
	"testing"
)

func TestResolverResolve(
	t *testing.T,
) {
	t.Parallel()

	cookieReader, store :=
		newResolverTestDependencies(t)

	resolver, err := newResolver(
		cookieReader,
		store,
		resolverTestNow,
	)
	if err != nil {
		t.Fatalf(
			"newResolver() error = %v",
			err,
		)
	}

	ctx := context.Background()
	request := newResolverTestRequest()

	got, err := resolver.Resolve(
		ctx,
		request,
	)
	if err != nil {
		t.Fatalf(
			"Resolve() error = %v",
			err,
		)
	}

	want := store.session

	if got != want {
		t.Errorf(
			"Resolve() session = %+v, want %+v",
			got,
			want,
		)
	}

	if cookieReader.readCalls != 1 {
		t.Errorf(
			"CookieReader.Read() calls = %d, want 1",
			cookieReader.readCalls,
		)
	}

	if cookieReader.request != request {
		t.Errorf(
			"CookieReader.Read() request = %p, want %p",
			cookieReader.request,
			request,
		)
	}

	if store.getCalls != 1 {
		t.Errorf(
			"Store.Get() calls = %d, want 1",
			store.getCalls,
		)
	}

	if store.getCtx != ctx {
		t.Errorf(
			"Store.Get() context = %v, want %v",
			store.getCtx,
			ctx,
		)
	}

	wantIDHash := resolverTestSessionIDHash(t)

	if store.idHash != wantIDHash {
		t.Errorf(
			"Store.Get() ID hash = %q, want %q",
			store.idHash,
			wantIDHash,
		)
	}

	if store.saveCalls != 0 {
		t.Errorf(
			"Store.Save() calls = %d, want 0",
			store.saveCalls,
		)
	}

	if store.deleteCalls != 0 {
		t.Errorf(
			"Store.Delete() calls = %d, want 0",
			store.deleteCalls,
		)
	}
}
