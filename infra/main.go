package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
        _, err := storage.NewBucket(ctx, "doppler-example", &storage.BucketArgs{
			Location: pulumi.String("US"),
		})
		if err != nil {
			return err
		}

		// Set up load balancer for serverless application (i.e. appengine)

		return nil
	})
}
