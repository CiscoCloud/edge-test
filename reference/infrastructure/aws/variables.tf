variable "region" {
  default = "us-west-1"
}

variable "availability_zone" {
  default = "us-west-1a"
}

variable "vpc_id" {
  default = "vpc-a4c62cc1"
}

variable "subnet_id" {
  default = "subnet-a9d510cc"
}

variable "ami" {
  default = "ami-9b7f90df"
}

variable "key_name" {
  default = "reference"
}

# Cluster config
variable "cluster_name" {
  default = "dev_cluster"
}

variable "ingress_tcp_ports" {
  default = "22,2181,2888,3888,5050,5051,6066,7000,7001,7077,7199,8080,8081,9042,9090,9092,9160"
}

# Exhibitor
variable "exhibitor_instance_type" {
  default = "t2.small"
}

variable "exhibitor_count" {
  default = "3"
}

# Kafka
variable "kafka_instance_type" {
  default = "t2.small"
}

variable "kafka_count" {
  default = "3"
}

# Mesos Master
variable "mesos_master_instance_type" {
  default = "t2.small"
}

variable "mesos_master_count" {
  default = "3"
}

# Mesos Slave
variable "mesos_slave_instance_type" {
  default = "t2.small"
}

variable "mesos_slave_count" {
  default = "3"
}


