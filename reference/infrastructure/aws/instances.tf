resource "aws_instance" "exhibitor" {
  count = "${var.exhibitor_count}"

  ami = "${var.ami}"
  instance_type = "${var.exhibitor_instance_type}"
  key_name = "${var.key_name}"

  subnet_id = "${var.subnet_id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_exhibitor"
  }
}

resource "aws_instance" "kafka" {
  count = "${var.kafka_count}"

  ami = "${var.ami}"
  instance_type = "${var.kafka_instance_type}"
  key_name = "${var.key_name}"

  subnet_id = "${var.subnet_id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_kafka"
  }
}

resource "aws_instance" "mesos_master" {
  count = "${var.mesos_master_count}"

  ami = "${var.ami}"
  instance_type = "${var.mesos_master_instance_type}"
  key_name = "${var.key_name}"

  subnet_id = "${var.subnet_id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_mesos_master"
  }
}

resource "aws_instance" "mesos_slave" {
  count = "${var.mesos_slave_count}"

  ami = "${var.ami}"
  instance_type = "${var.mesos_slave_instance_type}"
  key_name = "${var.key_name}"

  subnet_id = "${var.subnet_id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_mesos_slave"
  }
}

