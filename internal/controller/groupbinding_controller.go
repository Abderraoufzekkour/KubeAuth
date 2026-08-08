package controller

import (
	"context"
	"fmt"

	kube "github.com/Abderraoufzekkour/KubeAuth/internal/kubernetes"
)

type Reconciler struct {
	RBAC *kube.RBACManager
}

func NewReconciler() *Reconciler {
	return &Reconciler{
		RBAC: kube.NewRBACManager(),
	}
}

func (r *Reconciler) Reconcile(ctx context.Context, group, clusterRole, namespace, prefix string) error {
	fmt.Printf("Reconciling group: %s\n", group)

	opts := kube.GroupBindingOpts{
		GroupName:   group,
		GroupPrefix: prefix,
		ClusterRole: clusterRole,
		Namespace:   namespace,
	}

	if err := r.RBAC.SyncGroupBinding(ctx, opts); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	fmt.Printf("Reconcile complete: %s -> %s\n", group, clusterRole)
	return nil
}
