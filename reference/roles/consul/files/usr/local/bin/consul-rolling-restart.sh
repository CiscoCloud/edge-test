#!/bin/bash

%{ if {{ ansible_hostname }} in group[consul_servers_group]} %}
consul lock -n=1 locks/consul "/bin/bash -c \" 
    sleep 5; 
    /usr/local/bin/consul-wait-for-leader.sh || exit 1; 
    bash -c 'sleep 2; systemctl daemon-reload' & 
\""
{% else %}
service consul restart
{% endif %}

exit 0
