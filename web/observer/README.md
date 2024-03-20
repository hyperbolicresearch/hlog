# Observer

The observer is a web-based observability tool for `hlog` which provides a simple and intuitive user interface for monitoring purposes. It includes different sections the  most important of which are:

## Observability tools
- `Home` which is the default page providing an overview of the current state of the system and live updating section of new incoming data.
- `Dashboard` contains the essential observables (metrics) that are either built-in or user defined customized to the user-specific data.
- `Live tail` which is a generalization of the live tail from the home page, where the user can also query the data using the proposed methods.
- `Metrics` is where the user defines and configures the metrics that will be trackable in the `Dashboard` section.
- `Functions` -- will be defined later, but that's basically the next step of the development.

## Administrative tools
- `Admin` is only visible to users with admin permissions and provides API to handle rights essentially.
- `Settings` provides API for customizable configurations.