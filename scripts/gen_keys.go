package main

import (
	"fmt"
	"os"

	"bank-api/utils"
)

func main() {
	pub, priv, err := utils.GenerateTestPGPKeys()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== PGP PUBLIC KEY ===")
	fmt.Println(pub)
	fmt.Println("\n=== PGP PRIVATE KEY ===")
	fmt.Println(priv)
}
