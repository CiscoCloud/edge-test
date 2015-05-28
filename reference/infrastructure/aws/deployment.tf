# Placement config

variable "region" {
  default = "us-east-1"
}

# This is not a counter, it's rather a switch
# Make sure that if you disable availability zone
# you don't place any instances in it

variable "availability_zones" {
  default = {
    us-east-1a = 1 
    us-east-1b = 1
    us-east-1c = 1
  }
}

# Cluster config

variable "cluster_name" {
  default = "dev_cluster"
}

variable "mesos_masters" {
  default = {
    instance_type = "t2.small"
    us-east-1a = 2
    us-east-1b = 2
    us-east-1c = 1
  }
}

variable "mesos_slaves" {
  default = {
    instance_type = "t2.small"
    us-east-1a = 1
    us-east-1b = 1
    us-east-1c = 1
  }
}

variable "exhibitors" {
  default = {
    instance_type = "t2.small"
    us-east-1a = 1
    us-east-1b = 1
    us-east-1c = 1
  }
}

# Misc
variable "vpc_name" {
  default = "elodina"
}

# Please, provide existing(!) key name
variable "key_name" {
  default = "reference"
}

variable "amis" {
  default = {
    # Ubuntu 14.04 HVM/EBS
    us-east-1 = "ami-7eb2a716" 
    us-west-1 = "ami-9b7f90df"
    us-west-2 = "ami-3389b803"
  }
}

variable "ingress_tcp_ports" {
  default = "22,2181,2888,3888,5050,5051,6066,7000,7001,7077,7199,8080,8081,9042,9090,9092,9160"
}
