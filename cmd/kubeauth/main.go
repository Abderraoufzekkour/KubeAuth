package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Abderraoufzekkour/KubeAuth/internal/controller"
	"github.com/Abderraoufzekkour/KubeAuth/internal/keycloak"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubeauth",
	Short: "KubeAuth - Production-grade Keycloak OIDC operator for Kubernetes",
	Long: `KubeAuth automates the full Keycloak to Kubernetes authentication lifecycle.
It bootstraps OIDC, syncs Keycloak groups to Kubernetes RBAC, and verifies
tokens end-to-end — all in one tool.`,
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Test OIDC flow and print token claims",
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("keycloak-url")
		realm, _ := cmd.Flags().GetString("realm")
		clientID, _ := cmd.Flags().GetString("client-id")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		if url == "" || realm == "" {
			return fmt.Errorf("--keycloak-url and --realm are required")
		}

		fmt.Println("-------------------------------------------")
		fmt.Println("KubeAuth - OIDC Verification")
		fmt.Println("-------------------------------------------")

		cfg := keycloak.Config{
			URL:      url,
			Realm:    realm,
			ClientID: clientID,
		}

		client, err := keycloak.New(context.Background(), cfg)
		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}

		if err := client.VerifyOIDC(context.Background(), username, password); err != nil {
			return fmt.Errorf("verification failed: %w", err)
		}

		fmt.Println("-------------------------------------------")
		fmt.Println("Verification complete.")
		fmt.Printf("Issuer URL : %s/realms/%s\n", url, realm)
		fmt.Printf("Client ID  : %s\n", clientID)
		fmt.Println("-------------------------------------------")
		return nil
	},
}

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Run the KubeAuth Kubernetes operator",
	RunE: func(cmd *cobra.Command, args []string) error {
		group, _ := cmd.Flags().GetString("group")
		clusterRole, _ := cmd.Flags().GetString("cluster-role")
		namespace, _ := cmd.Flags().GetString("namespace")
		prefix, _ := cmd.Flags().GetString("prefix")

		fmt.Println("-------------------------------------------")
		fmt.Println("KubeAuth - Operator")
		fmt.Println("-------------------------------------------")

		r := controller.NewReconciler()
		if err := r.Reconcile(context.Background(), group, clusterRole, namespace, prefix); err != nil {
			return fmt.Errorf("reconcile failed: %w", err)
		}

		fmt.Println("-------------------------------------------")
		fmt.Println("Operator reconcile complete.")
		fmt.Println("-------------------------------------------")
		return nil
	},
}

func init() {
	verifyCmd.Flags().String("keycloak-url", "", "Keycloak base URL (e.g. https://auth.example.com)")
	verifyCmd.Flags().String("realm", "", "Keycloak realm name")
	verifyCmd.Flags().String("client-id", "kubernetes", "Keycloak OIDC client ID")
	verifyCmd.Flags().String("username", "", "Test username")
	verifyCmd.Flags().String("password", "", "Test user password")

	operatorCmd.Flags().String("group", "developers", "Keycloak group to sync")
	operatorCmd.Flags().String("cluster-role", "view", "Kubernetes ClusterRole to bind")
	operatorCmd.Flags().String("namespace", "", "Namespace for RoleBinding (empty = ClusterRoleBinding)")
	operatorCmd.Flags().String("prefix", "oidc:", "Group prefix for RBAC subject")

	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(operatorCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
