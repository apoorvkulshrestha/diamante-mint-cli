// load_account.go
package main

import (
	"errors"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/protocols/aurora"
)

// LoadAccount loads an account's details from the Diamante blockchain.
func LoadAccount(client *auroraclient.Client, accountId string) (*aurora.Account, error) {
	if accountId == "" {
		return nil, errors.New("account ID is required")
	}

	// Create the account request
	accountRequest := auroraclient.AccountRequest{AccountID: accountId}
	// Fetch the account details
	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		return nil, err
	}

	// The account is returned directly from the AccountDetail call.
	// Make sure the horizon.Account type is correctly imported from the auroraclient package.
	return &account, nil
}
