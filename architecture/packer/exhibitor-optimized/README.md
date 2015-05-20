exhibitor-optimized base AMIs
=====

This Packer script optimizes the Zookeeper/Exhibitor Ansible configuration process 
by pre-installing bare essential packages. It is recommended that complementing
Ansible playbooks/tasks still be executed as a sanity check.


## Pre-Built AMIs

The following AMI(s) were last built on 5/20/15.

```
us-west-1: ami-d5608991
```


## Building with Packer

`packer build packer.json`

