# Observable API

As the name indicates it, this so-called API is intended to serve the `Observer`. 

> You can check what the Observer [here](https://github.com/hyperbolicresearch/hlog/blob/dev/web/observer/README.md)

Through the current, we will have different communication channels with the Observer, such as:
1. `/live`: which is sending new logs in live streaming.
2. `/liveinit`: which is sending the latest k logs ingested in the system.
3. `/genericobservables`: returns the default measurables/observables.
4. `/observe/:observable_id`: returns the data for some user-defined observable. 
5. `/info`: returns generic information about the current system that is running.
