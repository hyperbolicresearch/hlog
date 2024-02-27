# Hlog

> Please, this project is under massive development and one shouldn't take anything written here or in the whole software for definitive facts. We are just conducting experimentations. Keep in mind that the creator(s) is/are just thinking loud (thinking in public) until otherwise said.

## Architectural description

Considering the challenges to be addressed here, We've made choices to maximize availability of the system and query performance. First of all, let's define what all of the buzzwords mean for us:

- `availability` in our context means that a client shouldn't send logs to the system and find out that the system is not available. for some reason.
- `query performance` means that after storing logs to our system, the user should be able to perform experiments on this data by querying billion of records in seconds.

We have then make the following choices:
1. Event-based collection of logs (through messaging system). In our case, we use `kafka`.
2. Column-based storage. We use `ClickHouse`.
3. For queries, we are making our best efforts to allow clients to write SQL queries while we are converting them to the way that we are storing theses semi-structured data to `ClickHouse`.

