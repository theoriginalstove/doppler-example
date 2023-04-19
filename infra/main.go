package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/appengine"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
        _, err := appengine.NewApplication(ctx, "doppler-example", &appengine.ApplicationArgs{
			Project:    pulumi.String("avocagrow-internal-tools"),
			LocationId: pulumi.String("us-central"),
		})
		if err != nil {
			return err
		}

        // Set up load balancer for serverless application (i.e. appengine)

		return nil
	})
}
