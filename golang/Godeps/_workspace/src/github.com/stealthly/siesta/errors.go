/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package siesta

import "errors"

// Signals that an end of file or stream has been reached unexpectedly.
var EOF = errors.New("End of file reached")

// Happens when a compressed message is empty.
var NoDataToUncompress = errors.New("No data to uncompress")

// A mapping for Kafka error code 0.
var NoError = errors.New("No error - it worked!")

// A mapping for Kafka error code -1.
var Unknown = errors.New("An unexpected server error")

// A mapping for Kafka error code 1.
var OffsetOutOfRange = errors.New("The requested offset is outside the range of offsets maintained by the server for the given topic/partition.")

// A mapping for Kafka error code 2.
var InvalidMessage = errors.New("Message contents does not match its CRC")

// A mapping for Kafka error code 3.
var UnknownTopicOrPartition = errors.New("This request is for a topic or partition that does not exist on this broker.")

// A mapping for Kafka error code 4.
var InvalidMessageSize = errors.New("The message has a negative size")

// A mapping for Kafka error code 5.
var LeaderNotAvailable = errors.New("In the middle of a leadership election and there is currently no leader for this partition and hence it is unavailable for writes.")

// A mapping for Kafka error code 6.
var NotLeaderForPartition = errors.New("You've just attempted to send messages to a replica that is not the leader for some partition. It indicates that the clients metadata is out of date.")

// A mapping for Kafka error code 7.
var RequestTimedOut = errors.New("Request exceeds the user-specified time limit in the request.")

// A mapping for Kafka error code 8.
var BrokerNotAvailable = errors.New("Broker is likely not alive.")

// A mapping for Kafka error code 9.
var ReplicaNotAvailable = errors.New("Replica is expected on a broker, but is not (this can be safely ignored).")

// A mapping for Kafka error code 10.
var MessageSizeTooLarge = errors.New("You've just attempted to produce a message of size larger than broker is allowed to accept.")

// A mapping for Kafka error code 11.
var StaleControllerEpochCode = errors.New("Broker-to-broker communication fault.")

// A mapping for Kafka error code 12.
var OffsetMetadataTooLargeCode = errors.New("You've jsut specified a string larger than configured maximum for offset metadata.")

// A mapping for Kafka error code 13.
var OffsetsLoadInProgressCode = errors.New("Offset loading is in progress. (Usually happens after a leader change for that offsets topic partition).")

// A mapping for Kafka error code 14.
var ConsumerCoordinatorNotAvailableCode = errors.New("Offsets topic has not yet been created.")

// A mapping for Kafka error code 15.
var NotCoordinatorForConsumerCode = errors.New("There is no coordinator for this consumer.")

// Mapping between Kafka error codes and actual error messages.
var BrokerErrors = map[int16]error{
	-1: Unknown,
	0:  NoError,
	1:  OffsetOutOfRange,
	2:  InvalidMessage,
	3:  UnknownTopicOrPartition,
	4:  InvalidMessageSize,
	5:  LeaderNotAvailable,
	6:  NotLeaderForPartition,
	7:  RequestTimedOut,
	8:  BrokerNotAvailable,
	9:  ReplicaNotAvailable,
	10: MessageSizeTooLarge,
	11: StaleControllerEpochCode,
	12: OffsetMetadataTooLargeCode,
	14: OffsetsLoadInProgressCode,
	15: ConsumerCoordinatorNotAvailableCode,
	16: NotCoordinatorForConsumerCode,
}
