package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	"sv-printer/internal/license"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "keygen":
		keygen()
	case "sign":
		sign()
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `sv-license — issue sv-printer licenses

Usage:
  sv-license keygen
      Generate an Ed25519 keypair. Prints the private key (hex) to stderr
      and the public key (hex) to stdout. Keep the private key secret.

  sv-license sign -key <private-key-hex> -customer <name> -tier <trial|beta|full> [flags]
      Sign a license and print the license.key JSON to stdout.

Sign flags:
  -id         license id (default: auto-generated)
  -customer   customer name (required)
  -tier       trial | beta | full (required)
  -expiry     RFC3339 expiry date, e.g. 2027-01-01T00:00:00Z (optional)
  -features   comma-separated: raw_print,websocket,serial (optional)
  -fingerprint device fingerprint (optional)
  -key        private key hex, or env SV_LICENSE_KEY (required)`)
}

func keygen() {
	pub, priv, err := license.GenerateKey()
	if err != nil {
		fmt.Fprintln(os.Stderr, "keygen:", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "PRIVATE KEY (keep secret, do not commit):")
	fmt.Fprintln(os.Stderr, hex.EncodeToString(priv))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "PUBLIC KEY (embed into internal/license/public.key):")

	fmt.Println(hex.EncodeToString(pub))
}

func sign() {
	fs := flag.NewFlagSet("sign", flag.ContinueOnError)

	id := fs.String("id", "", "license id")
	customer := fs.String("customer", "", "customer name")
	tier := fs.String("tier", "", "trial | beta | full")
	expiry := fs.String("expiry", "", "RFC3339 expiry date")
	features := fs.String("features", "", "comma-separated features")
	fingerprint := fs.String("fingerprint", "", "device fingerprint")
	key := fs.String("key", os.Getenv("SV_LICENSE_KEY"), "private key hex")

	if err := fs.Parse(os.Args[2:]); err != nil {
		os.Exit(2)
	}

	if *customer == "" || *tier == "" {
		fmt.Fprintln(os.Stderr, "sign: -customer and -tier are required")
		os.Exit(2)
	}

	if *key == "" {
		fmt.Fprintln(os.Stderr, "sign: private key required (-key or SV_LICENSE_KEY)")
		os.Exit(2)
	}

	privBytes, err := hex.DecodeString(strings.TrimSpace(*key))
	if err != nil || len(privBytes) != ed25519.PrivateKeySize {
		fmt.Fprintln(os.Stderr, "sign: invalid private key")
		os.Exit(2)
	}
	priv := ed25519.PrivateKey(privBytes)

	lic := license.License{
		LicenseID: *id,
		Product:   license.ProductName,
		Customer:  *customer,
		Tier:      license.Tier(*tier),
	}

	if *expiry != "" {
		lic.Expiry = expiry
	}
	if *fingerprint != "" {
		lic.Fingerprint = fingerprint
	}
	if *features != "" {
		for _, f := range strings.Split(*features, ",") {
			if f = strings.TrimSpace(f); f != "" {
				lic.Features = append(lic.Features, f)
			}
		}
	}

	out, err := license.Sign(priv, lic)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sign:", err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}
