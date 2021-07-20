# BIT Handler

'BIT Handler' microservice analyzes all system data stored in storage, and produce insights and potential vulnerabilities. The service Periodically invokes configured rules (failures) on the stored collected data and decides current BIT status w.r.t failures dependencies. This current BIT status is later written back to storage for other services to read.

BIT Status contains:
- **Reported failures**. A list of all the failures (description and count) found by the analyzer after crosschecking the configuration rules and test reports under the specific periodic check. BIT Trigger is set by default to 1 so statuses are being performed evey one second, but this duration is configurable under services configurations.

Further explanation about internal functions can be found in the source code.
