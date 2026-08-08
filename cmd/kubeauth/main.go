package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Abderraoufzekkour/KubeAuth/internal/controller"
	"github.com/Abderraoufzekkour/KubeAuth/internal/keycloak"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "verify":
		runVerify()
	case "operator":
		runOperator()
	case "--help", "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("KubeAuth - Production-grade Keycloak OIDC operator for Kubernetes")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  kubeauth verify     Test OIDC flow and print token claims")
	fmt.Println("  kubeauth operator   Run the Kubernetes operator")
	fmt.Println("  kubeauth help       Show this help message")
}

func runVerify() {
	cfg := keycloak.Config{
		URL:      "https://auth.example.com",
		Realm:    "myrealm",
		ClientID: "kubernetes",
	}

	client, err := keycloak.New(context.Background(), cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := client.VerifyOIDC(context.Background(), "testuser", "testpass"); err != nil {
		fmt.Printf("Verification failed: %v\n", err)
		os.Exit(1)
	}
}

func runOperator() {
	fmt.Println("Starting KubeAuth operator...")
	r := controller.NewReconciler()
	err := r.Reconcile(context.Background(), "developers", "view", "", "oidc:")
	if err != nil {
		fmt.Printf("Reconcile error: %v\n", err)
		os.Exit(1)
	}
}
