// The transport package is responsible to answer the question of
// how the clients (external world) is communication with the log
// aggregator core.
// Effectively, we make the choice of having different transports
// in order to give flexibility to the clients' implementations.
// Here is the motivation behind this choice:
// 	1. The primary and preferable way is through messaging
//     mechanism, since we have a by definition better availability
//     guarantee. In this regard, we prioritize Kafka.
//  2. We give the option to the users who may want to use a more
//     traditional way of communicating through HTTP.
// Both of these should implement the same interface, this will be
// a blueprint for further integrations.

package transport
