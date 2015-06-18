# Instances

# Exhibitors

resource "aws_instance" "exhibitor_zone_a" {
  count = "${lookup(var.exhibitors,"${var.region}a")}"
  availability_zone = "${var.region}a"

  ami = "${lookup(var.amis,"${var.region}")}"
  instance_type = "${lookup(var.exhibitors,"instance_type")}"
  key_name = "${var.key_name}"

  subnet_id = "${aws_subnet.zone_a_public.id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  root_block_device {
    volume_size = "${lookup(var.exhibitors,"volume_size")}"
  }

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_exhibitor"
  }

  depends_on = ["aws_internet_gateway.igw"]
}

resource "aws_instance" "exhibitor_zone_b" {
  count = "${lookup(var.exhibitors,"${var.region}b")}"
  availability_zone = "${var.region}b"

  ami = "${lookup(var.amis,"${var.region}")}"
  instance_type = "${lookup(var.exhibitors,"instance_type")}"
  key_name = "${var.key_name}"

  subnet_id = "${aws_subnet.zone_b_public.id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  root_block_device {
    volume_size = "${lookup(var.exhibitors,"volume_size")}"
  }

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_exhibitor"
  }

  depends_on = ["aws_internet_gateway.igw"]
}

resource "aws_instance" "exhibitor_zone_c" {
  count = "${lookup(var.exhibitors,"${var.region}c")}"
  availability_zone = "${var.region}c"

  ami = "${lookup(var.amis,"${var.region}")}"
  instance_type = "${lookup(var.exhibitors,"instance_type")}"
  key_name = "${var.key_name}"

  subnet_id = "${aws_subnet.zone_c_public.id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  root_block_device {
    volume_size = "${lookup(var.exhibitors,"volume_size")}"
  }

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_exhibitor"
  }

  depends_on = ["aws_internet_gateway.igw"]
}

# Mesos masters

resource "aws_instance" "mesos_master_zone_a" {
  count = "${lookup(var.mesos_masters,"${var.region}a")}"
  availability_zone = "${var.region}a"

  ami = "${lookup(var.amis,"${var.region}")}"
  instance_type = "${lookup(var.mesos_masters,"instance_type")}"
  key_name = "${var.key_name}"

  subnet_id = "${aws_subnet.zone_a_public.id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  root_block_device {
    volume_size = "${lookup(var.mesos_masters,"volume_size")}"
  }

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_mesos_master"
  }

  depends_on = ["aws_internet_gateway.igw"]
}

resource "aws_instance" "mesos_master_zone_b" {
  count = "${lookup(var.mesos_masters,"${var.region}b")}"
  availability_zone = "${var.region}b"

  ami = "${lookup(var.amis,"${var.region}")}"
  instance_type = "${lookup(var.mesos_masters,"instance_type")}"
  key_name = "${var.key_name}"

  subnet_id = "${aws_subnet.zone_b_public.id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  root_block_device {
    volume_size = "${lookup(var.mesos_masters,"volume_size")}"
  }

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_mesos_master"
  }

  depends_on = ["aws_internet_gateway.igw"]
}

resource "aws_instance" "mesos_master_zone_c" {
  count = "${lookup(var.mesos_masters,"${var.region}c")}"
  availability_zone = "${var.region}c"

  ami = "${lookup(var.amis,"${var.region}")}"
  instance_type = "${lookup(var.mesos_masters,"instance_type")}"
  key_name = "${var.key_name}"

  subnet_id = "${aws_subnet.zone_c_public.id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  root_block_device {
    volume_size = "${lookup(var.mesos_masters,"volume_size")}"
  }

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_mesos_master"
  }

  depends_on = ["aws_internet_gateway.igw"]
}


# Mesos slaves

resource "aws_instance" "mesos_slave_zone_a" {
  count = "${lookup(var.mesos_slaves,"${var.region}a")}"
  availability_zone = "${var.region}a"

  ami = "${lookup(var.amis,"${var.region}")}"
  instance_type = "${lookup(var.mesos_slaves,"instance_type")}"
  key_name = "${var.key_name}"

  subnet_id = "${aws_subnet.zone_a_public.id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  root_block_device {
    volume_size = "${lookup(var.mesos_slaves,"volume_size")}"
  }

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_mesos_slave"
  }

  depends_on = ["aws_internet_gateway.igw"]
}

resource "aws_instance" "mesos_slave_zone_b" {
  count = "${lookup(var.mesos_slaves,"${var.region}b")}"
  availability_zone = "${var.region}b"

  ami = "${lookup(var.amis,"${var.region}")}"
  instance_type = "${lookup(var.mesos_slaves,"instance_type")}"
  key_name = "${var.key_name}"

  subnet_id = "${aws_subnet.zone_b_public.id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  root_block_device {
    volume_size = "${lookup(var.mesos_slaves,"volume_size")}"
  }

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_mesos_slave"
  }

  depends_on = ["aws_internet_gateway.igw"]
}

resource "aws_instance" "mesos_slave_zone_c" {
  count = "${lookup(var.mesos_slaves,"${var.region}c")}"
  availability_zone = "${var.region}c"

  ami = "${lookup(var.amis,"${var.region}")}"
  instance_type = "${lookup(var.mesos_slaves,"instance_type")}"
  key_name = "${var.key_name}"

  subnet_id = "${aws_subnet.zone_c_public.id}"
  security_groups = ["${aws_security_group.clusterwide.id}"]

  root_block_device {
    volume_size = "${lookup(var.mesos_slaves,"volume_size")}"
  }

  associate_public_ip_address = true

  tags {
    Name = "${var.cluster_name}_mesos_slave"
  }

  depends_on = ["aws_internet_gateway.igw"]
}
