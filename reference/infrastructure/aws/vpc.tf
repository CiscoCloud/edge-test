# VPC

resource "aws_vpc" "vpc" {
    cidr_block = "172.16.0.0/16"

    tags {
        Name = "${var.vpc_name}"
    }
}

resource "aws_internet_gateway" "igw" {
	vpc_id = "${aws_vpc.vpc.id}"

	tags {
		Name = "${var.vpc_name}"
	}
}

# Public subnets

resource "aws_subnet" "zone_a_public" {
	vpc_id = "${aws_vpc.vpc.id}"
	count = "${lookup(var.availability_zones,"${var.region}a")}"

	cidr_block = "172.16.0.0/22"
	availability_zone = "${var.region}a"

	tags {
		Name = "${var.vpc_name}-${var.region}a"
	}
}

resource "aws_subnet" "zone_b_public" {
	vpc_id = "${aws_vpc.vpc.id}"
	count = "${lookup(var.availability_zones,"${var.region}a")}"

	cidr_block = "172.16.4.0/22"
	availability_zone = "${var.region}b"

	tags {
		Name = "${var.vpc_name}-${var.region}b"
	}
}

resource "aws_subnet" "zone_c_public" {
	vpc_id = "${aws_vpc.vpc.id}"
	count = "${lookup(var.availability_zones,"${var.region}a")}"

	cidr_block = "172.16.8.0/22"
	availability_zone = "${var.region}c"

	tags {
		Name = "${var.vpc_name}-${var.region}c"
	}
}

# Routing

resource "aws_route_table" "routes" {
	vpc_id = "${aws_vpc.vpc.id}"

	route {
		cidr_block = "0.0.0.0/0"
		gateway_id = "${aws_internet_gateway.igw.id}"
	}
}

resource "aws_route_table_association" "zone_a_public" {
	subnet_id = "${aws_subnet.zone_a_public.id}"
	route_table_id = "${aws_route_table.routes.id}"

	count = "${lookup(var.availability_zones,"${var.region}a")}"
}

resource "aws_route_table_association" "zone_b_public" {
	subnet_id = "${aws_subnet.zone_b_public.id}"
	route_table_id = "${aws_route_table.routes.id}"

	count = "${lookup(var.availability_zones,"${var.region}b")}"
}

resource "aws_route_table_association" "zone_c_public" {
	subnet_id = "${aws_subnet.zone_c_public.id}"
	route_table_id = "${aws_route_table.routes.id}"

	count = "${lookup(var.availability_zones,"${var.region}c")}"
}
