// +build integration

package checkout

import (
	"context"
	"testing"

	pb "github.com/emcfarlane/checkout/checkoutpb"
	"go.uber.org/zap/zaptest"
)

func TestVoid(t *testing.T) {
	s, err := NewService(zaptest.NewLogger(t), db, NewMockBank())
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	rsp, err := s.Authorize(ctx, testAuthReq("4532111111111112"))
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Void(ctx, &pb.VoidRequest{Id: rsp.Id})
	if err != nil {
		t.Fatal(err)
	}
	// Duplicate void should fail.
	_, err = s.Void(ctx, &pb.VoidRequest{Id: rsp.Id})
	if err == nil {
		t.Fatal("expected error")
	}
}
