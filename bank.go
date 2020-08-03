package checkout

import (
	"fmt"
	"sync"
	"time"
)

// Bank implements a one or multiple Banking services to route payments to.
type Bank interface {
	Authorize(id, pan string, ccv string, exp time.Time, amount uint64, currency string) error
	Capture(id string, amount uint64) error
	Refund(id string, amount uint64) error
}

type method int

const (
	methodAuthorisation method = iota
	methodCapture
	methodRefund
)

type rule struct {
	method method
	pan    string
}

type ruleBank struct {
	rules map[rule]error // List of rules

	mu   sync.RWMutex
	pans map[string]string // Maps ID to PAN in memory, would be stored securely.
}

// MockBank implements the following rules, all else succeed:
// 4000 0000 0000 0119: authorisation failure
// 4000 0000 0000 0259: capture failure
// 4000 0000 0000 3238: refund failure
func NewMockBank() Bank {
	return &ruleBank{
		rules: map[rule]error{
			rule{methodAuthorisation, "4000000000000119"}: fmt.Errorf("authorisation failure"),
			rule{methodCapture, "4000000000000259"}:       fmt.Errorf("capture failure"),
			rule{methodRefund, "4000000000003238"}:        fmt.Errorf("refund failure"),
		},
	}
}

func (b *ruleBank) getPAN(id string) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	pan, ok := b.pans[id]
	if !ok {
		return "", fmt.Errorf("missing pan")
	}
	return pan, nil
}

func (b *ruleBank) setPAN(id, pan string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.pans == nil {
		b.pans = make(map[string]string)
	}
	b.pans[id] = pan
}

func (b *ruleBank) Authorize(
	id, pan, ccv string, exp time.Time, amount uint64, currency string,
) error {
	b.setPAN(id, pan)
	return b.rules[rule{methodAuthorisation, pan}]
}
func (b *ruleBank) Capture(id string, amount uint64) error {
	pan, err := b.getPAN(id)
	if err != nil {
		return err
	}
	return b.rules[rule{methodCapture, pan}]
}
func (b *ruleBank) Refund(id string, amount uint64) error {
	pan, err := b.getPAN(id)
	if err != nil {
		return err
	}
	return b.rules[rule{methodRefund, pan}]
}
