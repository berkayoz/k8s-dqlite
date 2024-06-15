package embeddedctl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/canonical/k8s-dqlite/pkg/embedded"
	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func newEmbeddedClient(storageDir string) (*clientv3.Client, error) {
	instance, err := embedded.New(storageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize instance: %w", err)
	}
	return instance.NewLocalClient()
}

func jsonOutput(i any) error {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON output: %w", err)
	}
	fmt.Println(string(b))
	return nil
}

func newCommand(f func(context.Context, *clientv3.Client) (any, error)) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := newEmbeddedClient(flagStorageDir)
		if err != nil {
			return fmt.Errorf("failed to initialize embedded client: %w", err)
		}
		resp, err := f(cmd.Context(), client)
		return jsonOutput(resp)
	}
}
