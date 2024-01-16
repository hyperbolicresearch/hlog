// The storage package basically answers to the question of how do
// we store the logs and where do we do it.
// It is intended to be built in an extensive and integrable way
// that would basically allow any other database system to be easily
// added provided that they implement the interfaces necessary.
// For now, the intent is to integrate PostgreSQL and Clickhouse.
// The idea of using two different storage comes from having two
// different needs:
// 	1. We need to run real-time analytics on the data and provide
//     the user of the log software with insights about the
//     stucture of the logs, statistical information and others. For
//     this reason, we believe that an OLAP storage system is needed
//     Hence the choice of Clickhouse.
//	2. We need to persist the data for arbitrary time in order.
//     While it's possible to store it on Clickhouse, we believe
//     that this requires a more traditional OLTP storage system,
//     where operations such as SELECT * FROM table can be performed
//     in an optimized way.

package storage
