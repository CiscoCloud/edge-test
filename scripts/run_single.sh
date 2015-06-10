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
if [ -z $2 ]; then
	echo 'USAGE: ./run_single.sh <url> <num_per_batch>'
	exit
fi

URL=$1
NUM_PER_BATCH=$2
NUM_MSG=0
BATCH_URL=$URL
for i in $(seq 2 1 $NUM_PER_BATCH)
do
	BATCH_URL+=" $URL"
done
while [  true ]; do
	if [[ $NUM_MSG%10 -eq 0 ]]
	then
 		echo produced $NUM_MSG messages
 	fi

 	eval 'curl --silent -d "{\"line\": \"i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line i am line\",\"source\": \"generated\",\"tag\": null,\"logtypeid\": 5,\"timings\": []}" --header "Content-Type: application/json" -v $BATCH_URL 2> /dev/null'
 	let NUM_MSG=NUM_MSG+1 
done