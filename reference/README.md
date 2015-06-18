# Elodina
Please, make sure that cluster is running before actually trying to provision it.
Cluster deployment can be done using Terraform states in infrastructure folder.

### Using hybrid inventory
Assuming you're in a repo directory and your cluster is called 'dev-cluster'
```
export EC2_INI_PATH=$(pwd)/inventory/ec2.ini
```
AWS credentials should also be in env
```
export AWS_ACCESS_KEY_ID='YOUR_AWS_API_KEY'
export AWS_SECRET_ACCESS_KEY='YOUR_AWS_API_SECRET_KEY'
```
Finally
```
ansible-playbook site.yml --extra-vars="cluster=dev-cluster"
```

### TODO's:
 * Vagrantfile - port forwarding, groups naming standadization
