/**
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Siesta is a low-level Apache Kafka client in Go.

package siesta

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const InvalidOffset int64 = -1

// Connector is an interface that should provide ways to clearly interact with Kafka cluster and hide all broker management stuff from user.
type Connector interface {
	// GetTopicMetadata is primarily used to discover leaders for given topics and how many partitions these topics have.
	// Passing it an empty topic list will retrieve metadata for all topics in a cluster.
	GetTopicMetadata(topics []string) (*TopicMetadataResponse, error)

	// GetAvailableOffset issues an offset request to a specified topic and partition with a given offset time.
	// More on offset time here - https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol#AGuideToTheKafkaProtocol-OffsetRequest
	GetAvailableOffset(topic string, partition int32, offsetTime OffsetTime) (int64, error)

	// Fetch issues a single fetch request to a broker responsible for a given topic and partition and returns a FetchResponse that contains messages starting from a given offset.
	Fetch(topic string, partition int32, offset int64) (*FetchResponse, error)

	// GetOffset gets the offset for a given group, topic and partition from Kafka. A part of new offset management API.
	GetOffset(group string, topic string, partition int32) (int64, error)

	// CommitOffset commits the offset for a given group, topic and partition to Kafka. A part of new offset management API.
	CommitOffset(group string, topic string, partition int32, offset int64) error

	// Tells the Connector to close all existing connections and stop.
	// This method is NOT blocking but returns a channel which will get a single value once the closing is finished.
	Close() <-chan bool
}

// ConnectorConfig is used to pass multiple configuration values for a Connector
type ConnectorConfig struct {
	// BrokerList is a bootstrap list to discover other brokers in a cluster. At least one broker is required.
	BrokerList []string

	// ReadTimeout is a timeout to read the response from a TCP socket.
	ReadTimeout time.Duration

	// WriteTimeout is a timeout to write the request to a TCP socket.
	WriteTimeout time.Duration

	// ConnectTimeout is a timeout to connect to a TCP socket.
	ConnectTimeout time.Duration

	// Sets whether the connection should be kept alive.
	KeepAlive bool

	// A keep alive period for a TCP connection.
	KeepAliveTimeout time.Duration

	// Maximum number of open connections for a connector.
	MaxConnections int

	// Maximum number of open connections for a single broker for a connector.
	MaxConnectionsPerBroker int

	// Maximum fetch size in bytes which will be used in all Consume() calls.
	FetchSize int32

	// The minimum amount of data the server should return for a fetch request. If insufficient data is available the request will block
	FetchMinBytes int32

	// The maximum amount of time the server will block before answering the fetch request if there isn't sufficient data to immediately satisfy FetchMinBytes
	FetchMaxWaitTime int32

	// Number of retries to get topic metadata.
	MetadataRetries int

	// Backoff value between topic metadata requests.
	MetadataBackoff time.Duration

	// Number of retries to commit an offset.
	CommitOffsetRetries int

	// Backoff value between commit offset requests.
	CommitOffsetBackoff time.Duration

	// Number of retries to get consumer metadata.
	ConsumerMetadataRetries int

	// Backoff value between consumer metadata requests.
	ConsumerMetadataBackoff time.Duration

	// Client id that will be used by a connector to identify client requests by broker.
	ClientId string
}

// Returns a new ConnectorConfig with sane defaults.
func NewConnectorConfig() *ConnectorConfig {
	return &ConnectorConfig{
		ReadTimeout:             5 * time.Second,
		WriteTimeout:            5 * time.Second,
		ConnectTimeout:          5 * time.Second,
		KeepAlive:               true,
		KeepAliveTimeout:        1 * time.Minute,
		MaxConnections:          5,
		MaxConnectionsPerBroker: 5,
		FetchSize:               1024000,
		FetchMaxWaitTime:        1000,
		MetadataRetries:         5,
		MetadataBackoff:         200 * time.Millisecond,
		CommitOffsetRetries:     5,
		CommitOffsetBackoff:     200 * time.Millisecond,
		ConsumerMetadataRetries: 15,
		ConsumerMetadataBackoff: 500 * time.Millisecond,
		ClientId:                "siesta",
	}
}

//Validates this ConnectorConfig. Returns a corresponding error if the ConnectorConfig is invalid and nil otherwise.
func (this *ConnectorConfig) Validate() error {
	if this == nil {
		return errors.New("Please provide a ConnectorConfig.")
	}

	if len(this.BrokerList) == 0 {
		return errors.New("BrokerList must have at least one broker.")
	}

	if this.ReadTimeout < time.Millisecond {
		return errors.New("ReadTimeout must be at least 1ms.")
	}

	if this.WriteTimeout < time.Millisecond {
		return errors.New("WriteTimeout must be at least 1ms.")
	}

	if this.ConnectTimeout < time.Millisecond {
		return errors.New("ConnectTimeout must be at least 1ms.")
	}

	if this.KeepAliveTimeout < time.Millisecond {
		return errors.New("KeepAliveTimeout must be at least 1ms.")
	}

	if this.MaxConnections < 1 {
		return errors.New("MaxConnections cannot be less than 1.")
	}

	if this.MaxConnectionsPerBroker < 1 {
		return errors.New("MaxConnectionsPerBroker cannot be less than 1.")
	}

	if this.FetchSize < 1 {
		return errors.New("FetchSize cannot be less than 1.")
	}

	if this.MetadataRetries < 0 {
		return errors.New("MetadataRetries cannot be less than 0.")
	}

	if this.MetadataBackoff < time.Millisecond {
		return errors.New("MetadataBackoff must be at least 1ms.")
	}

	if this.CommitOffsetRetries < 0 {
		return errors.New("CommitOffsetRetries cannot be less than 0.")
	}

	if this.CommitOffsetBackoff < time.Millisecond {
		return errors.New("CommitOffsetBackoff must be at least 1ms.")
	}

	if this.ConsumerMetadataRetries < 0 {
		return errors.New("ConsumerMetadataRetries cannot be less than 0.")
	}

	if this.ConsumerMetadataBackoff < time.Millisecond {
		return errors.New("ConsumerMetadataBackoff must be at least 1ms.")
	}

	if this.ClientId == "" {
		return errors.New("ClientId cannot be empty.")
	}

	return nil
}

// A default (and only one for now) Connector implementation for Siesta library.
type DefaultConnector struct {
	config         ConnectorConfig
	leaders        map[string]map[int32]*brokerLink
	links          []*brokerLink
	bootstrapLinks []*brokerLink
	lock           sync.Mutex

	//offset coordination part
	offsetCoordinators map[string]int32
}

// Creates a new DefaultConnector with a given ConnectorConfig. May return an error if the passed config is invalid.
func NewDefaultConnector(config *ConnectorConfig) (*DefaultConnector, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	leaders := make(map[string]map[int32]*brokerLink)
	connector := &DefaultConnector{
		config:             *config,
		leaders:            leaders,
		offsetCoordinators: make(map[string]int32),
	}

	return connector, nil
}

// Returns a string representation of this DefaultConnector.
func (this *DefaultConnector) String() string {
	return "Default Connector"
}

// GetTopicMetadata is primarily used to discover leaders for given topics and how many partitions these topics have.
// Passing it an empty topic list will retrieve metadata for all topics in a cluster.
func (this *DefaultConnector) GetTopicMetadata(topics []string) (*TopicMetadataResponse, error) {
	for i := 0; i <= this.config.MetadataRetries; i++ {
		if metadata, err := this.getMetadata(topics); err == nil {
			return metadata, nil
		}

		Debugf(this, "GetTopicMetadata for %s failed after %d try", topics, i)
		time.Sleep(this.config.MetadataBackoff)
	}

	return nil, errors.New(fmt.Sprintf("Could not get topic metadata for %s after %d retries", topics, this.config.MetadataRetries))
}

// GetAvailableOffset issues an offset request to a specified topic and partition with a given offset time.
func (this *DefaultConnector) GetAvailableOffset(topic string, partition int32, offsetTime OffsetTime) (int64, error) {
	request := new(OffsetRequest)
	request.AddPartitionOffsetRequestInfo(topic, partition, offsetTime, 1)
	response, err := this.sendToAllAndReturnFirstSuccessful(request, this.offsetValidator)
	if response != nil {
		return response.(*OffsetResponse).Offsets[topic][partition].Offsets[0], err
	} else {
		return -1, err
	}
}

// Fetch issues a single fetch request to a broker responsible for a given topic and partition and returns a FetchResponse that contains messages starting from a given offset.
func (this *DefaultConnector) Fetch(topic string, partition int32, offset int64) (*FetchResponse, error) {
	link := this.getLeader(topic, partition)
	if link == nil {
		leader, err := this.tryGetLeader(topic, partition, this.config.MetadataRetries)
		if err != nil {
			return nil, err
		}
		link = leader
	}

	request := new(FetchRequest)
	request.MinBytes = this.config.FetchMinBytes
	request.MaxWaitTime = this.config.FetchMaxWaitTime
	request.AddFetch(topic, partition, offset, this.config.FetchSize)
	bytes, err := this.syncSendAndReceive(link, request)
	if err != nil {
		this.removeLeader(topic, partition)
		return nil, err
	}

	decoder := NewBinaryDecoder(bytes)
	response := new(FetchResponse)
	decodingErr := response.Read(decoder)
	if decodingErr != nil {
		this.removeLeader(topic, partition)
		Errorf(this, "Could not decode a FetchResponse. Reason: %s", decodingErr.Reason())
		return nil, decodingErr.Error()
	}

	return response, nil
}

// GetOffset gets the offset for a given group, topic and partition from Kafka. A part of new offset management API.
func (this *DefaultConnector) GetOffset(group string, topic string, partition int32) (int64, error) {
	coordinator, err := this.getOffsetCoordinator(group)
	if err != nil {
		return InvalidOffset, err
	}

	request := NewOffsetFetchRequest(group)
	request.AddOffset(topic, partition)
	bytes, err := this.syncSendAndReceive(coordinator, request)
	if err != nil {
		return InvalidOffset, err
	}
	response := new(OffsetFetchResponse)
	decodingErr := this.decode(bytes, response)
	if decodingErr != nil {
		Errorf(this, "Could not decode an OffsetFetchResponse. Reason: %s", decodingErr.Reason())
		return InvalidOffset, decodingErr.Error()
	}

	if topicOffsets, exist := response.Offsets[topic]; !exist {
		return InvalidOffset, fmt.Errorf("OffsetFetchResponse does not contain information about requested topic")
	} else {
		if offset, exists := topicOffsets[partition]; !exists {
			return InvalidOffset, fmt.Errorf("OffsetFetchResponse does not contain information about requested partition")
		} else if offset.Error != NoError {
			return InvalidOffset, offset.Error
		} else {
			return offset.Offset, nil
		}
	}
}

// CommitOffset commits the offset for a given group, topic and partition to Kafka. A part of new offset management API.
func (this *DefaultConnector) CommitOffset(group string, topic string, partition int32, offset int64) error {
	for i := 0; i <= this.config.CommitOffsetRetries; i++ {
		if err := this.tryCommitOffset(group, topic, partition, offset); err == nil {
			return nil
		}

		Debugf(this, "Failed to commit offset %d for group %s, topic %s, partition %d after %d try", offset, group, topic, partition, i)
		time.Sleep(this.config.CommitOffsetBackoff)
	}

	return errors.New(fmt.Sprintf("Could not get commit offset %d for group %s, topic %s, partition %d after %d retries", offset, group, topic, partition, this.config.CommitOffsetRetries))
}

//func (this *DefaultConnector) Produce(message Message) error {
//	//TODO keep in mind: If RequiredAcks == 0 the server will not send any response (this is the only case where the server will not reply to a request)
//	panic("Not implemented yet")
//}

// Tells the Connector to close all existing connections and stop.
// This method is NOT blocking but returns a channel which will get a single value once the closing is finished.
func (this *DefaultConnector) Close() <-chan bool {
	closed := make(chan bool)
	go func() {
		this.closeBrokerLinks()
		for _, link := range this.bootstrapLinks {
			link.stop <- true
		}
		this.bootstrapLinks = nil
		this.links = nil
		closed <- true
	}()

	return closed
}

func (this *DefaultConnector) closeBrokerLinks() {
	for _, link := range this.links {
		link.stop <- true
	}
}

func (this *DefaultConnector) refreshMetadata(topics []string) {
	if len(this.bootstrapLinks) == 0 {
		for i := 0; i < len(this.config.BrokerList); i++ {
			broker := this.config.BrokerList[i]
			hostPort := strings.Split(broker, ":")
			if len(hostPort) != 2 {
				panic(fmt.Sprintf("incorrect broker connection string: %s", broker))
			}

			port, err := strconv.Atoi(hostPort[1])
			if err != nil {
				panic(fmt.Sprintf("incorrect port in broker connection string: %s", broker))
			}

			this.bootstrapLinks = append(this.bootstrapLinks, newBrokerLink(&Broker{NodeId: -1, Host: hostPort[0], Port: int32(port)},
				this.config.KeepAlive,
				this.config.KeepAliveTimeout,
				this.config.MaxConnectionsPerBroker))
		}
	}

	response, err := this.sendToAllLinks(this.links, NewTopicMetadataRequest(topics), this.topicMetadataValidator(topics))
	if err != nil {
		Warnf(this, "Could not get topic metadata from all known brokers, trying bootstrap brokers...")
		if response, err = this.sendToAllLinks(this.bootstrapLinks, NewTopicMetadataRequest(topics), this.topicMetadataValidator(topics)); err != nil {
			Errorf(this, "Could not get topic metadata from all known brokers")
			return
		}
	}
	this.refreshLeaders(response.(*TopicMetadataResponse))
}

func (this *DefaultConnector) refreshLeaders(response *TopicMetadataResponse) {
	brokers := make(map[int32]*brokerLink)
	for _, broker := range response.Brokers {
		brokers[broker.NodeId] = newBrokerLink(broker, this.config.KeepAlive, this.config.KeepAliveTimeout, this.config.MaxConnectionsPerBroker)
	}

	if len(brokers) != 0 && len(response.TopicMetadata) != 0 {
		this.closeBrokerLinks()
		this.links = make([]*brokerLink, 0)
	}

	for _, metadata := range response.TopicMetadata {
		for _, partitionMetadata := range metadata.PartitionMetadata {
			if leader, exists := brokers[partitionMetadata.Leader]; exists {
				this.putLeader(metadata.TopicName, partitionMetadata.PartitionId, leader)
			} else {
				Warnf(this, "Topic Metadata response has no leader present for topic %s, parition %d", metadata.TopicName, partitionMetadata.PartitionId)
				//TODO: warn about incomplete broker list
			}
		}
	}
}

func (this *DefaultConnector) getMetadata(topics []string) (*TopicMetadataResponse, error) {
	response, err := this.sendToAllAndReturnFirstSuccessful(NewTopicMetadataRequest(topics), this.topicMetadataValidator(topics))
	if response != nil {
		return response.(*TopicMetadataResponse), err
	} else {
		return nil, err
	}
}

func (this *DefaultConnector) tryGetLeader(topic string, partition int32, retries int) (*brokerLink, error) {
	for i := 0; i <= retries; i++ {
		this.refreshMetadata([]string{topic})
		if link := this.getLeader(topic, partition); link != nil {
			return link, nil
		}
		time.Sleep(this.config.MetadataBackoff)
	}

	return nil, errors.New(fmt.Sprintf("Could not get leader for %s:%d after %d retries", topic, partition, retries))
}

func (this *DefaultConnector) getLeader(topic string, partition int32) *brokerLink {
	leadersForTopic, exists := this.leaders[topic]
	if !exists {
		return nil
	}

	return leadersForTopic[partition]
}

func (this *DefaultConnector) putLeader(topic string, partition int32, leader *brokerLink) {
	Tracef(this, "putLeader for topic %s, partition %d - %s", topic, partition, leader.broker)
	if _, exists := this.leaders[topic]; !exists {
		this.leaders[topic] = make(map[int32]*brokerLink)
	}

	exists := false
	for _, link := range this.links {
		if *link.broker == *leader.broker {
			exists = true
			break
		}
	}

	if !exists {
		this.links = append(this.links, leader)
	}

	this.leaders[topic][partition] = leader
}

func (this *DefaultConnector) removeLeader(topic string, partition int32) {
	if leadersForTopic, exists := this.leaders[topic]; exists {
		delete(leadersForTopic, partition)
	}
}

func (this *DefaultConnector) refreshOffsetCoordinator(group string) error {
	for i := 0; i <= this.config.ConsumerMetadataRetries; i++ {
		if err := this.tryRefreshOffsetCoordinator(group); err == nil {
			return nil
		}

		Debugf(this, "Failed to get consumer coordinator for group %s after %d try", group, i)
		time.Sleep(this.config.ConsumerMetadataBackoff)
	}

	return fmt.Errorf("Could not get consumer coordinator for group %s after %d retries", group, this.config.ConsumerMetadataRetries)
}

func (this *DefaultConnector) tryRefreshOffsetCoordinator(group string) error {
	request := NewConsumerMetadataRequest(group)

	response, err := this.sendToAllAndReturnFirstSuccessful(request, this.consumerMetadataValidator)
	if err != nil {
		Infof(this, "Could not get consumer metadata from all known brokers")
		return err
	}
	this.offsetCoordinators[group] = response.(*ConsumerMetadataResponse).CoordinatorId

	return nil
}

func (this *DefaultConnector) getOffsetCoordinator(group string) (*brokerLink, error) {
	coordinatorId, exists := this.offsetCoordinators[group]
	if !exists {
		err := this.refreshOffsetCoordinator(group)
		if err != nil {
			return nil, err
		}
		coordinatorId = this.offsetCoordinators[group]
	}

	Debugf(this, "Offset coordinator for group %s: %d", group, coordinatorId)

	var brokerLink *brokerLink
	for _, link := range this.links {
		if link.broker.NodeId == coordinatorId {
			brokerLink = link
			break
		}
	}

	if brokerLink == nil {
		return nil, fmt.Errorf("Could not find broker with node id %d", coordinatorId)
	}

	return brokerLink, nil
}

func (this *DefaultConnector) tryCommitOffset(group string, topic string, partition int32, offset int64) error {
	coordinator, err := this.getOffsetCoordinator(group)
	if err != nil {
		return err
	}

	request := NewOffsetCommitRequest(group)
	request.AddOffset(topic, partition, offset, time.Now().Unix(), "")

	bytes, err := this.syncSendAndReceive(coordinator, request)
	if err != nil {
		return err
	}

	response := new(OffsetCommitResponse)
	decodingErr := this.decode(bytes, response)
	if decodingErr != nil {
		Errorf(this, "Could not decode an OffsetCommitResponse. Reason: %s", decodingErr.Reason())
		return decodingErr.Error()
	}

	if topicErrors, exist := response.Offsets[topic]; !exist {
		return fmt.Errorf("OffsetCommitResponse does not contain information about requested topic")
	} else {
		if partitionError, exist := topicErrors[partition]; !exist {
			return fmt.Errorf("OffsetCommitResponse does not contain information about requested partition")
		} else if partitionError != NoError {
			return partitionError
		}
	}

	return nil
}

func (this *DefaultConnector) decode(bytes []byte, response Response) *DecodingError {
	decoder := NewBinaryDecoder(bytes)
	decodingErr := response.Read(decoder)
	if decodingErr != nil {
		Errorf(this, "Could not decode a response. Reason: %s", decodingErr.Reason())
		return decodingErr
	}

	return nil
}

func (this *DefaultConnector) sendToAllAndReturnFirstSuccessful(request Request, check func([]byte) Response) (Response, error) {
	if len(this.links) == 0 {
		this.refreshMetadata(nil)
	}

	response, err := this.sendToAllLinks(this.links, request, check)
	if err != nil {
		response, err = this.sendToAllLinks(this.bootstrapLinks, request, check)
	}

	return response, err
}

func (this *DefaultConnector) sendToAllLinks(links []*brokerLink, request Request, check func([]byte) Response) (Response, error) {
	if len(links) == 0 {
		return nil, errors.New("Empty broker list")
	}

	responses := make(chan *rawResponseAndError, len(links))
	for i := 0; i < len(links); i++ {
		link := links[i]
		go func() {
			bytes, err := this.syncSendAndReceive(link, request)
			responses <- &rawResponseAndError{bytes, link, err}
		}()
	}

	var response *rawResponseAndError
	for i := 0; i < len(links); i++ {
		response = <-responses
		if response.err == nil {
			if checkResult := check(response.bytes); checkResult != nil {
				return checkResult, nil
			} else {
				response.err = errors.New("Check result did not pass")
			}
		}

		Infof(this, "Could not process request with broker %s:%d", response.link.broker.Host, response.link.broker.Port)
	}

	return nil, response.err
}

func (this *DefaultConnector) syncSendAndReceive(link *brokerLink, request Request) ([]byte, error) {
	id, conn, err := link.getConnection()
	if err != nil {
		link.failed()
		return nil, err
	}

	if err := this.send(id, conn, request); err != nil {
		link.failed()
		return nil, err
	}

	bytes, err := this.receive(conn)
	if err != nil {
		link.failed()
		return nil, err
	}

	link.succeeded()
	link.connectionPool.Return(conn)
	return bytes, err
}

func (this *DefaultConnector) send(correlationId int32, conn *net.TCPConn, request Request) error {
	writer := NewRequestWriter(correlationId, this.config.ClientId, request)
	bytes := make([]byte, writer.Size())
	encoder := NewBinaryEncoder(bytes)
	writer.Write(encoder)

	conn.SetWriteDeadline(time.Now().Add(this.config.WriteTimeout))
	_, err := conn.Write(bytes)
	return err
}

func (this *DefaultConnector) receive(conn *net.TCPConn) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(this.config.ReadTimeout))
	header := make([]byte, 8)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}

	decoder := NewBinaryDecoder(header)
	length, err := decoder.GetInt32()
	if err != nil {
		return nil, err
	}
	response := make([]byte, length-4)
	_, err = io.ReadFull(conn, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (this *DefaultConnector) topicMetadataValidator(topics []string) func(bytes []byte) Response {
	return func(bytes []byte) Response {
		response := new(TopicMetadataResponse)
		err := this.decode(bytes, response)
		if err != nil {
			return nil
		}

		if len(topics) > 0 {
			for _, topic := range topics {
				var topicMetadata *TopicMetadata
				for _, topicMetadata = range response.TopicMetadata {
					if topicMetadata.TopicName == topic {
						break
					}
				}

				if topicMetadata.Error != NoError {
					Infof(this, "Topic metadata err: %s", topicMetadata.Error)
					return nil
				}

				for _, partitionMetadata := range topicMetadata.PartitionMetadata {
					if partitionMetadata.Error != NoError && partitionMetadata.Error != ReplicaNotAvailable {
						Infof(this, "Partition metadata err: %s", partitionMetadata.Error)
						return nil
					}
				}
			}
		}

		return response
	}
}

func (this *DefaultConnector) consumerMetadataValidator(bytes []byte) Response {
	response := new(ConsumerMetadataResponse)
	err := this.decode(bytes, response)
	if err != nil || response.Error != NoError {
		return nil
	}

	return response
}

func (this *DefaultConnector) offsetValidator(bytes []byte) Response {
	response := new(OffsetResponse)
	err := this.decode(bytes, response)
	if err != nil {
		return nil
	}
	for _, offsets := range response.Offsets {
		for _, offset := range offsets {
			if offset.Error != NoError {
				return nil
			}
		}
	}

	return response
}

type brokerLink struct {
	broker                    *Broker
	connectionPool            *connectionPool
	lastConnectTime           time.Time
	lastSuccessfulConnectTime time.Time
	failedAttempts            int
	correlationIds            chan int32
	stop                      chan bool
}

func newBrokerLink(broker *Broker, keepAlive bool, keepAliveTimeout time.Duration, maxConnectionsPerBroker int) *brokerLink {
	brokerConnect := fmt.Sprintf("%s:%d", broker.Host, broker.Port)
	correlationIds := make(chan int32)
	stop := make(chan bool)

	go correlationIdGenerator(correlationIds, stop)

	return &brokerLink{
		broker:         broker,
		connectionPool: newConnectionPool(brokerConnect, maxConnectionsPerBroker, keepAlive, keepAliveTimeout),
		correlationIds: correlationIds,
		stop:           stop,
	}
}

func (this *brokerLink) failed() {
	this.lastConnectTime = time.Now()
	this.failedAttempts++
}

func (this *brokerLink) succeeded() {
	timestamp := time.Now()
	this.lastConnectTime = timestamp
	this.lastSuccessfulConnectTime = timestamp
}

func (this *brokerLink) getConnection() (int32, *net.TCPConn, error) {
	correlationId := <-this.correlationIds
	conn, err := this.connectionPool.Borrow()
	return correlationId, conn, err
}

func correlationIdGenerator(out chan int32, stop chan bool) {
	var correlationId int32 = 0
	for {
		select {
		case out <- correlationId:
			correlationId++
		case <-stop:
			return
		}
	}
}

type rawResponseAndError struct {
	bytes []byte
	link  *brokerLink
	err   error
}
