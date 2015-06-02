# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
# 
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/sh -Eux

echo 'Starting'
if [ -z $1 ]; then
	echo 'USAGE: ./run_single.sh <url>'
	exit
fi

URL=$1
PAYLOAD='{"line": "i am line","source": "generated","tag": null,"logtypeid": 5,"timings": [{"eventName": "key1", "value": 123000},{"eventName": "key2", "value": 124000}]}'
NUM_MSG=0
while [  true ]; do
	if [[ $NUM_MSG%1000 -eq 0 ]]
	then
 		echo produced $NUM_MSG messages
 	fi

 	eval 'curl --data-binary "@test-json" --header "Content-Type: application/json" $URL'
 	let NUM_MSG=NUM_MSG+1 
done