// main.go
package main

import (
	"fmt"
	"log"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/manifoldco/promptui"
)

func main() {
	fmt.Println("Welcome to the Diamante Minting Application")

	issuerSeedKey := promptForInput("Enter Issuer Seed Key")
	assetCode := promptForInput("Enter Asset Code")
	amount := promptForInput("Enter Amount to Mint")
	distributorAddress := promptForInput("Enter Distributor Account ID")
	distributorSeedKey := promptForInput("Enter Distributor Seed Key")

	// Initialize the Diamante SDK client
	diamnetClient := auroraclient.DefaultTestNetClient

	// Prepare parameters for the minting operation
	mintParams := MintParams{
		IssuerSeedKey:                    issuerSeedKey,
		AssetCode:                        assetCode,
		Amount:                           amount,
		DistributorAddress:               distributorAddress,
		DistributorSeedKey:               distributorSeedKey,
		IsDistributorAccountNeedToCreate: false, // This should be set based on your application logic
		IsDistributorAccountNeedToTrust:  true,  // Typically set to true if we're minting a new asset
		IsIssuerAccountLock:              false, // Set this based on your requirements
	}

	fmt.Println("Preparing to mint asset...")
	err := mintAsset(diamnetClient, mintParams)
	if err != nil {
		log.Fatalf("Minting failed: %v", err)
	}

	fmt.Println("Asset minted successfully")
}

func promptForInput(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}
	return result
}
