package kubernetes

import (
	"context"
	"fmt"
)

type GroupBindingOpts struct {
	GroupName   string
	GroupPrefix string
	ClusterRole string
	Namespace   string
}

type RBACManager struct{}

func NewRBACManager() *RBACManager {
	return &RBACManager{}
}

func (r *RBACManager) SyncGroupBinding(ctx context.Context, opts GroupBindingOpts) error {
	subject := opts.GroupPrefix + opts.GroupName
	if opts.Namespace == "" {
		fmt.Printf("ClusterRoleBinding: %s → %s\n", subject, opts.ClusterRole)
	} else {
		fmt.Printf("RoleBinding in [%s]: %s → %s\n", opts.Namespace, subject, opts.ClusterRole)
	}
	return nil
}

func (r *RBACManager) DeleteGroupBinding(ctx context.Context, opts GroupBindingOpts) error {
	fmt.Printf("Deleted binding for group: %s\n", opts.GroupName)
	return nil
}
