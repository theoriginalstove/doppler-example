package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/appengine"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		project, err := projects.GetProject(ctx, &projects.GetProjectArgs{
			Filter: "name:avocagrow-internal-tools",
		})
		if err != nil {
			return err
		}

		_, err = appengine.NewApplication(ctx, "doppler-example", &appengine.ApplicationArgs{
			Project:    pulumi.String(project.Id),
			LocationId: pulumi.String("us-central"),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
