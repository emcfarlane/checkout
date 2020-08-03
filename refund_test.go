// +build integration

package checkout

import (
	"context"
	"testing"

	pb "github.com/emcfarlane/checkout/checkoutpb"
	"go.uber.org/zap/zaptest"
)

func TestRefund(t *testing.T) {
	s, err := NewService(zaptest.NewLogger(t), db, NewMockBank())
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	rsp, err := s.Authorize(ctx, testAuthReq("4532111111111112"))
	if err != nil {
		t.Fatal(err)
	}

	amount := rsp.Amount
	_, err = s.Capture(ctx, &pb.CaptureRequest{Id: rsp.Id, Amount: amount})
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Refund(ctx, &pb.RefundRequest{Id: rsp.Id, Amount: 1})
	if err != nil {
		t.Fatal(err)
	}
	// Captures should now fail
	_, err = s.Capture(ctx, &pb.CaptureRequest{Id: rsp.Id, Amount: amount})
	if err == nil {
		t.Fatal("expected error")
	}
	_, err = s.Refund(ctx, &pb.RefundRequest{Id: rsp.Id, Amount: amount})
	if err == nil {
		t.Fatal("expected error")
	}
	rsp, err = s.Refund(ctx, &pb.RefundRequest{Id: rsp.Id, Amount: amount - 1})
	if err != nil {
		t.Fatal(err)
	}
	if rsp.AmountCaptured != 0 {
		t.Fatalf("expected full request: %d", rsp.AmountCaptured)
	}
}
