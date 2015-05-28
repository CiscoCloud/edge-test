# Security groups and rules

resource "aws_security_group" "clusterwide" {
  name = "${var.cluster_name}"
  description = "${var.cluster_name} security group"
  vpc_id = "${aws_vpc.vpc.id}"

  tags {
      Name = "${var.cluster_name}"
  }

  depends_on = ["aws_vpc.vpc"]
}

resource "aws_security_group_rule" "ingress_inter_node" {
  security_group_id = "${aws_security_group.clusterwide.id}"
  type = "ingress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  self = true
}

resource "aws_security_group_rule" "egress_inter_node" {
  security_group_id = "${aws_security_group.clusterwide.id}"
  type = "egress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  self = true
}

resource "aws_security_group_rule" "ingress_global_tcp" {
  security_group_id = "${aws_security_group.clusterwide.id}"

  count = "${length(split(",",var.ingress_tcp_ports))}"

  type = "ingress"
  from_port = "${element(split(",",var.ingress_tcp_ports), count.index)}"
  to_port = "${element(split(",",var.ingress_tcp_ports), count.index)}"
  protocol = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "egress_global_all" {
  security_group_id = "${aws_security_group.clusterwide.id}"

  type = "egress"
  from_port = "0"
  to_port = "0"
  protocol = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}
