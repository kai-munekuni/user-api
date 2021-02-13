package config

import "os"

// GCPSAKey gcp sa key json
func GCPSAKey() []byte {
	return []byte(os.Getenv("GCP_SA_KEY"))
}
