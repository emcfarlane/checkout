// +build integration

package checkout

import (
	"context"
	"testing"

	pb "github.com/emcfarlane/checkout/checkoutpb"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func testAuthReq(pan string) *pb.AuthorizeRequest {
	return &pb.AuthorizeRequest{
		Pan:      pan,
		ExpMonth: 1,
		ExpYear:  2022,
		Cvv:      "000",
		Amount:   500,
		Currency: "gbp",
	}
}

func TestAuthorize(t *testing.T) {
	s, err := NewService(zaptest.NewLogger(t), db, NewMockBank())
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	rsp, err := s.Authorize(ctx, testAuthReq("4532111111111112"))
	if err != nil {
		t.Fatal(err)
	}
	if rsp.Id == "" {
		t.Fatalf("missing id")
	}

	// Check auth fails on MOCK pan
	_, err = s.Authorize(ctx, testAuthReq("4000 0000 0000 0119"))
	se, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected failure for PAN")
	}
	if se.Code() != codes.FailedPrecondition {
		t.Fatalf("invalid code %s", se.Code().String())
	}
}
