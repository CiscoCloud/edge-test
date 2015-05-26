#!/bin/bash -Eux

pushd /opt
  # install oracle java 7 and set as default java
  curl -LO --cookie oraclelicense=accept-securebackup-cookie \
    http://download.oracle.com/otn-pub/java/jdk/7u79-b15/jdk-7u79-linux-x64.rpm
  yum localinstall -y jdk-7u79-linux-x64.rpm
  alternatives --install /usr/bin/java java /usr/java/jdk1.7.0_79/bin/java 2
  alternatives --set java /usr/java/jdk1.7.0_79/bin/java

  # install maven so we can build exhibitor
  curl -LO http://mirror.nexcess.net/apache/maven/maven-3/3.3.3/binaries/apache-maven-3.3.3-bin.tar.gz
  tar xzf apache-maven-3.3.3-bin.tar.gz -C /usr/local
  ln -s /usr/local/apache-maven-3.3.3 /usr/local/maven
  echo 'export M2_HOME=/usr/local/maven' > /etc/profile.d/maven.sh
  echo 'export PATH=${M2_HOME}/bin:${PATH}' >> /etc/profile.d/maven.sh
popd
