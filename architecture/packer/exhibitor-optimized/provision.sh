#!/bin/bash -Eux

# add webupd8 ppa
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys EEA14886
echo deb http://ppa.launchpad.net/webupd8team/java/ubuntu trusty main > /etc/apt/sources.list.d/ppa_launchpad_net_webupd8team_java_ubuntu.list
apt-get update

# accept license and install java
echo debconf shared/accepted-oracle-license-v1-1 select true | sudo debconf-set-selections
apt-get install -y oracle-java7-installer

# necessary to build exhibitor
apt-get install -y maven
