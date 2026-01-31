package workerpool

import (
	"context"
	"errors"
	"testing"
)

func TestPool_Run_PreservesOrder(t *testing.T) {
	p := New[string](3)
	in := []string{"aa", "b", "ccc"}

	out, err := p.Run(context.Background(), in, func(ctx context.Context, s string) (string, error) {
		return s + "!", nil
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(out) != 3 || out[0] != "aa!" || out[1] != "b!" || out[2] != "ccc!" {
		t.Fatalf("unexpected out: %#v", out)
	}
}

func TestPool_Run_EarlyCancelOnError(t *testing.T) {
	sentinel := errors.New("boom")
	p := New[int](2)
	_, err := p.Run(context.Background(), []int{1, 2, 3}, func(ctx context.Context, v int) (int, error) {
		if v == 2 {
			return 0, sentinel
		}
		return v, nil
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel err, got %v", err)
	}
}
