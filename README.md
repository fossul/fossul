![](../images/fossul_logo.png)
# Welcome to the Fossul Framework
Fossul is a container-native backup and recovery framework. It aims to provide an application consistent backup, restore as well as recorvery for container-native applications and databases. It's goal is to provide enterprise backup and recovery capabilities enjoyed in traditional world, container-native, in order to increase container adoption and make the migration from traditial to container world that much easier. The challenge of backup, especially in the container-driven world is signifigant. Each application or database has it's own tools, methods and procedures. Applications are much more fluid, dynamic and provide greater abstractions inside containers. In other words backup and restore is even more of a challenge than in traditional world. In addition there is an enormous variation of application and storage vendor capabilities integrated into container platforms (for example snapshots that impact backup as well as restore process).

Fossul addresses the challanges by building a true modular plugin framework where every workflow function is executed inside a plugin. Backup and recovery like anything else at a high-level is always the same process (quiesce or dump data, back it up, unquiesce, archiving, etc). The only thing that changes is the specific tools and lower-level proceadures within those high-level functions. The fossul belief is that the backup/recovery process can be standardized or democratized using a dynamic plugin-driven workflow and framework. Welcome to the Fossul Framework where you can not only imortalize your applications but their data as well!

# How it works at 10000ft
Fossul has a server microservice for handling messaging, providing workflow-engine, scheduling, state, API and more. The server microservice triggers a workflow, via scheduler or manual which then executes a series of steps. Each step is a function executed in a plugin that is exposed through an additional service API. Currently there are two plugin services storage for storage / archive plugins and app for applications. A service can expose more than one type of plugin. While today the Fossul Framework is focused on backup / restore it could be extended to support virtually any process as a workflow where integration through plugins is favorable. 

If we look at a very, very simple backup example workflow it is typically something like ```quiesce application -> backup data -> unquiesce application -> archive data to secondary storage```. In Fossul the workflow is executed on the server but the implementation happens in plugins running on plugin services. Application functions like quiesce and unquiesce is handled by app service and a plugin that knows how to do this for say mysql (could be any database). Storage functions like backup data is handled by the storage service and a plugin that can backup the data given a specific storage, say Rook (could be any storage). Archive data is hadled by storage service as well and a plugin that can archive data to a secondary storage, perhaps S3 (could be any archive storage). The key to making this all work is a simple yet powerful model to share information between the services. If your interested read more in the documentation below.

# Contribution
If you want to provide additional features, please feel free to contribute via pull requests and please document your pull request thouroughly. We are happy to track and discuss ideas, topics and requests via 'Issues'. It is recommended you should start by writing plugins if you want to make development contributions, as you can jump-in and write them with almost no learning curve.

# Documentation
Additional documentation and detail

## General
If you are new it is highly recommend to read more about Fossul to understand the various components and architecture.
* [General Information and Architecture](docs/GENERAL.md)

## Plugins
More details about the plugin architecture and specific plugins. Read the general information before diving into plugins.
* [Plugin Information and Details](docs/PLUGINS.md)

## Build and Development Environment
If you want to be a contributor the first step is getting a development environment up and running.
* [Setup Development Environment](docs/BUILD.md)

## Using Fossul
Follow the deployment guide to get Fossul up and running.
* [Deploy Fossul and Get It Running!](docs/DEPLOY.md)

Once Fossul is operational start using it by following some of the examples in the getting started guide.
* [Getting Started](docs/GETTING_STARTED.md)





