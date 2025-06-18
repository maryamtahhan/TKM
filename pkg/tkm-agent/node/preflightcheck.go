package node

import (
	"log"

	"github.com/redhat-et/TKDK/tcv/pkg/accelerator"
	"github.com/redhat-et/TKDK/tcv/pkg/client"
)

func RunPreflightChecks(accs map[string]accelerator.Accelerator, imageName string) error {
	log.Printf("Performing preflight checks for image %s...", imageName)
	err := client.PreflightCheck(imageName)
	if err != nil {
		log.Fatalf("Incompatible system for image: %v", err)
	}

	return nil
}
