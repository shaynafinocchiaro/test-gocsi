package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

var nodeGetCapabilitiesCmd = &cobra.Command{
	Use:     "get-capabilities",
	Aliases: []string{"capabilities"},
	Short:   `invokes the rpc "NodeGetCapabilities"`,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx, cancel := context.WithTimeout(root.ctx, root.timeout)
		defer cancel()

		rep, err := node.client.NodeGetCapabilities(
			ctx,
			&csi.NodeGetCapabilitiesRequest{})
		if err != nil {
			return err
		}

		for _, cap := range rep.Capabilities {
			fmt.Println(cap.Type)
		}

		return nil
	},
}

func init() {
	nodeCmd.AddCommand(nodeGetCapabilitiesCmd)
}
