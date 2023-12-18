package main

import "math/big"

// EthToWei converts an Ethereum value in ETH (as a string) to wei (as *big.Int).
func EthToWei(eth string) (*big.Int, bool) {
	// Create a big.Float from the ETH string
	value, ok := new(big.Float).SetString(eth)
	if !ok {
		return nil, false // Could not parse ETH value
	}

	// Create a big.Float for the conversion factor (1 ETH = 10^18 wei)
	multiplier := new(big.Float).SetInt(big.NewInt(1e18))

	// Multiply the ETH value by the conversion factor to get wei
	value.Mul(value, multiplier)

	// Convert the big.Float result to a big.Int
	wei := new(big.Int)
	value.Int(wei) // Extracts the integer part of the big.Float

	return wei, true
}
