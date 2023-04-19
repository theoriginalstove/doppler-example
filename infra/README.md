# Deploying with Pulumi

To deploy this application and its associated infrastructure with Pulumi, you should:

- Have the Pulumi CLI installed
- Have Go installed (either of the latest 2 versions should work)
- Have GCP CLI installed and configured for your account

Then follow these steps:

 1. Clone this repository to your local system
 1. Switch to the infra directory
 1. Run `go mod tidy`
 1. Run `pulumi stack init <name>` to create a new stack.
 1. Set your desired GCP region with `pulumi config set gcp:region <region-name>`.
 1. Run `pulumi up`
