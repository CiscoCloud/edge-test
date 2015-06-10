Pre-configuration
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

cd aws

./terraform plan -input=false # Use default variables
./terraform apply -input=false

cd ../../
ansible-playbook site.yml --extra-vars="cluster=dev_cluster"

```

When the work is done, tear down the cluster
```
terraform destroy -input=false
```

## Notes
Terraform is designed to manage single deployment at a time, so please, don't try to deploy multiple clusters with it without completely wiping off previous deployment, as you just won't be able to destroy old resources without manual intervention or other type of sourcery.