package flags

import "github.com/spf13/cobra"

type TagCreateOptions struct {
	Description string
}

func (o *TagCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "Description of this tag")
}
