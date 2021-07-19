# BIT Test Results Exporter

'BIT Test Results Exporter' microservice pushes collected data asynchronously to 'Indexer' service for it to pass it on and store it under the frameworks' storage. The service implements 'sidecar' design pattern for each network segment in order to ensure network traffic and enable prioritization.

The 'Exporter' uses a thread-safe priority queue in order to store all received test results. It posts them through the indexer to the storage under strict bandwidth limitation pre-defined under systems' configurations.
The service responds to end-points with acknowledgement that the data has received, and uses separate asynchronous routines (threads) in order to periodically pushes the reports forward, under the bandwidth limitation and prioritization.

The purpose of this type of design pattern is to locate the cohesive task of receiving and pushing test results, with the primary application, but place it inside its own process, providing a homogeneous interface for platform services across languages.

Further explanation about internal functions can be found in the source code.

**Resources and open source packages:** https://github.com/bgadrian/data-structures/tree/master/priorityqueue (MIT License).