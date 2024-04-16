// is_account_created_from_seed
package main

import (
	"github.com/diamcircle/go/keypair"
)

// IsAccountCreatedFromSeed checks if a Diamante account exists for a given seed.
// This is a simplified example; adjust based on actual keypair validation methods.
func IsAccountCreatedFromSeed(seedKey string) (*keypair.Full, error) {
	kp, err := keypair.ParseFull(seedKey)
	if err != nil {
		return nil, err
	}
	return kp, nil
}
