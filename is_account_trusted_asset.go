// is_account_trusted_asset
package main

import (
	"github.com/diamcircle/go/clients/auroraclient"
)

// IsAccountTrustedAsset checks if an account trusts a specified asset.
func IsAccountTrustedAsset(client *auroraclient.Client, accountId, assetCode, assetIssuer string) (bool, error) {
	account, err := LoadAccount(client, accountId)
	if err != nil {
		return false, err
	}

	for _, balance := range account.Balances {
		if balance.Asset.Code == assetCode && balance.Asset.Issuer == assetIssuer {
			return true, nil
		}
	}

	return false, nil
}
