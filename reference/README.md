# Elodina

### Using hybrid inventory
Assuming you're in a repo directory
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
ansible -i inventory -m ping 'exhibitor:&test_cluster' 
```

### TODO's:
 * Vagrantfile - port forwarding, groups naming standadization
