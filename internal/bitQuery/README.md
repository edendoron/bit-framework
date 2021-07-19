# BIT Query

'BIT Query' microservice Receives user queries through standardizes API, fetches data directly from framework's 'Storage' service based on these queries and reports results back to the user.

Fetched data includes 'BIT status' and 'test results' (raw reports).
Query types for fetching test reports include:
- Query by time
- Query by field
- Query by tag-set

A query for fetching available user groups from storage is also provided.

Further explanation about internal functions can be found in the source code.