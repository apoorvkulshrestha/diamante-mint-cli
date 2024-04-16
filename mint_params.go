// mint_params.go
package main

import (
	"fmt"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/txnbuild"
)

// MintParams holds parameters for the minting operation.
type MintParams struct {
	IssuerSeedKey                    string
	IssuerAddress                    string
	AssetCode                        string
	Amount                           string
	DistributorAddress               string
	DistributorSeedKey               string
	IsDistributorAccountNeedToCreate bool
	IsDistributorAccountNeedToTrust  bool
	IsIssuerAccountLock              bool
}

func mintAsset(client *auroraclient.Client, params MintParams) error {
	// Parse the issuer keypair
	issuerKP, err := keypair.ParseFull(params.IssuerSeedKey)
	if err != nil {
		return fmt.Errorf("failed to parse issuer seed key: %v", err)
	}

	// Define the asset to be minted
	asset := txnbuild.CreditAsset{Code: params.AssetCode, Issuer: issuerKP.Address()}

	// Check if the distributor account needs to be created
	if params.IsDistributorAccountNeedToCreate {
		_, err := CreateAccount(client, params.DistributorSeedKey, "100") // Assuming minimum balance to be 100
		if err != nil {
			return fmt.Errorf("failed to create distributor account: %v", err)
		}
	}

	// Check if the distributor account needs to trust the asset
	if params.IsDistributorAccountNeedToTrust {
		err = CreateTrustline(client, params.DistributorSeedKey, asset)
		if err != nil {
			return fmt.Errorf("failed to create trustline: %v", err)
		}
	}

	// Load the issuer account
	issuerAccount, err := LoadAccount(client, issuerKP.Address())
	if err != nil {
		return fmt.Errorf("failed to load issuer account: %v", err)
	}

	// Issuing the asset from the issuer to the distributor
	op := txnbuild.Payment{
		Destination: params.DistributorAddress,
		Amount:      params.Amount,
		Asset:       asset,
	}

	// Build the transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        issuerAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Memo:                 txnbuild.MemoText("Minting new asset"),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to build transaction: %v", err)
	}

	// Sign the transaction with the issuer's keypair
	tx, err = tx.Sign(network.TestNetworkPassphrase, issuerKP)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Submit the transaction to the network
	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	fmt.Println("Transaction successfully submitted. Hash:", resp.Hash)

	// If the issuer account should be locked after minting
	if params.IsIssuerAccountLock {
		err := LockAccount(client, issuerKP)
		if err != nil {
			return fmt.Errorf("failed to lock issuer account: %v", err)
		}
	}

	return nil
}

// CreateAccount creates a new account on the network with the given seed and initial balance.
func CreateAccount(client *auroraclient.Client, sourceSeed string, initialBalance string) (*aurora.Account, error) {
	sourceKP, err := keypair.ParseFull(sourceSeed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source seed: %v", err)
	}

	// Generate a new keypair for the new account
	newAccountKP, err := keypair.Random()
	if err != nil {
		return nil, fmt.Errorf("failed to create keypair for new account: %v", err)
	}

	// Create an account creation operation
	createAccountOp := txnbuild.CreateAccount{
		Destination: newAccountKP.Address(),
		Amount:      initialBalance,
	}

	// Fetch the sequence number for the source account
	sourceAccount, err := client.AccountDetail(auroraclient.AccountRequest{AccountID: sourceKP.Address()})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source account details: %v", err)
	}

	// Build the transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&createAccountOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: %v", err)
	}

	// Sign the transaction with the source account's seed
	tx, err = tx.Sign(network.TestNetworkPassphrase, sourceKP)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Submit the transaction to the network
	_, err = client.SubmitTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction for account creation: %v", err)
	}

	// Return the details of the new account
	return &aurora.Account{
		AccountID: newAccountKP.Address(),
		Sequence:  "0", // New accounts start with sequence number 0
	}, nil
}

// CreateTrustline establishes a trustline from an account to an asset.
func CreateTrustline(client *auroraclient.Client, sourceSeed string, asset txnbuild.Asset) error {
	sourceKP, err := keypair.ParseFull(sourceSeed)
	if err != nil {
		return fmt.Errorf("failed to parse source seed: %v", err)
	}

	sourceAccount, err := client.AccountDetail(auroraclient.AccountRequest{AccountID: sourceKP.Address()})
	if err != nil {
		return fmt.Errorf("failed to fetch source account details: %v", err)
	}

	// Convert the provided asset to a ChangeTrustAsset
	changeTrustAsset, err := asset.ToChangeTrustAsset()
	if err != nil {
		return fmt.Errorf("failed to convert asset for ChangeTrust operation: %v", err)
	}

	// Create a trustline operation
	changeTrustOp := txnbuild.ChangeTrust{
		Line: changeTrustAsset,
	}

	// Build the transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&changeTrustOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to build transaction: %v", err)
	}

	// Sign the transaction with the source account's seed
	tx, err = tx.Sign(network.TestNetworkPassphrase, sourceKP)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Submit the transaction to the network
	_, err = client.SubmitTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to submit transaction for trustline creation: %v", err)
	}

	return nil
}

// LockAccount locks the issuer account by setting its master weight and thresholds to 0.
func LockAccount(client *auroraclient.Client, issuerKP *keypair.Full) error {
	sourceAccount, err := client.AccountDetail(auroraclient.AccountRequest{AccountID: issuerKP.Address()})
	if err != nil {
		return fmt.Errorf("failed to fetch source account details: %v", err)
	}

	setOptionsOp := txnbuild.SetOptions{
		MasterWeight:    txnbuild.NewThreshold(0),
		LowThreshold:    txnbuild.NewThreshold(0),
		MediumThreshold: txnbuild.NewThreshold(0),
		HighThreshold:   txnbuild.NewThreshold(0),
		Signer:          &txnbuild.Signer{Address: issuerKP.Address(), Weight: 0},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&setOptionsOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to build transaction: %v", err)
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, issuerKP)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	_, err = client.SubmitTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to submit transaction for account locking: %v", err)
	}

	return nil
}
