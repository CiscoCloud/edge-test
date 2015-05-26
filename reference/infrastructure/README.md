
# Terraform-ation for infrastructure
Hashicorp Terraform version >= 0.5.1 is required

Pre-configuration
```
export AWS_ACCESS_KEY_ID='YOUR_AWS_API_KEY'
export AWS_SECRET_ACCESS_KEY='YOUR_AWS_API_SECRET_KEY'

# Shared state source configuration
terraform remote config -backend=S3 \
 -backend-config="bucket=elodina-terraform" \
 -backend-config="key=tf-state" \
 -backend-config="region=us-west-1"
```

Deployment to AWS
(at the moment, dev_cluster only)
```
cd aws

terraform plan -input=false # Use default variables
terraform apply -input=false
```

When the work is done, tear down the cluster
```
terraform destroy -input=false
```