#!/bin/sh

cd apache-jmeter*/bin
chmod +x *.sh

./jmeter.sh -s -Jserver.rmi.localport=$PORT0
