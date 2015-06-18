
# Terraform-ation for infrastructure
Hashicorp Terraform version >= 0.5.1 is required

### Pre-configuration
```
export AWS_ACCESS_KEY_ID='YOUR_AWS_API_KEY'
export AWS_SECRET_ACCESS_KEY='YOUR_AWS_API_SECRET_KEY'

# Shared state source configuration 
(in case multiple people are supposed to manage infrastructure)
terraform remote config -backend=S3 \
 -backend-config="bucket=elodina-terraform" \
 -backend-config="key=tf-state" \
 -backend-config="region=us-west-1"
```
If you're expected to be the only person managing an infrastructure, you don't have to perform a remote config, so you can proceed directly to deployment step.

### Deployment to AWS
Pretty much everything that you might want to tweak resides in deployment.tf file. It's self-explanatory.
Use it wisely, as Terraform is kinda gentle, and barely has some human-related errors prevention mechanism.
```
cd aws

terraform plan -input=false # Use default variables
terraform apply -input=false
```
Give cluster a few minutes to spin up.

### Tearing down the cluster
```
terraform destroy -input=false
```

## Notes
Terraform is designed to manage single deployment at a time, so please, don't try to deploy multiple clusters with it without completely wiping off previous deployment, as you just won't be able to destroy old resources without manual intervention or other type of sourcery.
