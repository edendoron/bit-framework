# BIT Config

'BIT Config' microservice extracts all system configuration from a designated path or repository, and stores it under the framework's storage, managed by the 'Storage' service. Service is used only during configuration stage.

System configuration includes:
- **Failures definitions and rules**, which currently can be found under `./configs/config_failures`. Config failures
  defines the exact conditions that should be met in order to be considered as a vulnerability. Failures are read by the handler in the beginning of it's run, and are used by it to analyze a specific, or a sequence of reports that may meet the conditions. Further details about failures configuration can be found in the section below.
- **User groups Filtering rules**, which currently can be found under `./configs/config_user_froups_filtering`. User groups filtering defines the rules that apply upon different user groups. Every failure test ID may be ignored by a specific user group. In this case it will be listed under the corresponding user group's masked test ID's.

### Failures

Failures definition and explanation can be found mainly in the protobuf generated go file located in `./configs/rafael.com/bina/bit/bit.pb.go` and in the guidance presentation.
To complete all definitions and assumptions we took, we will detail more about specific fields of a configured failure.

- `Failure.FailureExaminationRule.MatchingTag` indicates the tag set that needed to be present in the test report. One of the test report tag sets should contain the failure specified tag set in order for it to be considered as this failure violation.
- `Failure.FailureExaminationRule.FailureValueCriteria.ThresholdMode` - either 'WITHIN' or 'OUTOF', report will be considered violation of the failure if the tested field values are within or out of the range respectively.
- `Failure.FailureExaminationRule.FailureValueCriteria.ExceedingType and ExceedingValue` - report will be considered violation of the failure if the tested field values are within or out of the range according to threshold mode, with a deviation of some value or percentage from the range allowed. For example: in our usage example, a voltage value of 12 is of course out of the range 1-7 so will be applied as violation. Voltage value of 5 is with-in the range so will not be considered violation, but also a value of 7.6 is under the limitations of 10% deviation of the configured range, and so wil also not be considered as violation. similar approach applies on 'WITHIN' exceeding type.
- `Failure.FailureTimeCriteria` - fields 'windowSize' and 'failuresCount' are only relevant for 'SLIDING_WINDOW' type failures. 'NO_WINDOW' failures count will be documented under 'BITStatus.ReportedFailure.exactCount' of the corresponding BITStatus report failure, if found.


Further explanation about internal functions can be found in the source code.