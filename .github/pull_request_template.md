### Summary and Scope

EXPLAIN WHY THIS PR IS NECESSARY. WHAT IS IMPACTED?
IS THIS A NEW FEATURE OR CRITICAL BUG FIX? SUMMARIZE WHAT CHANGED.

DOES THIS CHANGE INVOLVE ANY SCHEME CHANGES?  Y/N

REMINDER: HAVE YOU INCREMENTED VERSION NUMBERS? E.G., .spec, Chart.yaml

REMINDER 2: HAVE YOU UPDATED THE COPYRIGHT PER hpe GUIDELINES: Copyright 2014-2021 Hewlett Packard Enterprise Development LP    ? Y/N

NOTE FOR RELEASE BRANCHES: YOU NO LONGER NEED TO UPDATE THE `.version` VERSION ON EVERY PR IF THE PR DOES NOT CHANGE THE CONTENTS OF THE GENERATED CONTAINER IMAGE.  IF YOU DO UPDATE `.version`, YOU ALSO NEED TO UPDATE THE `Chart.yaml` VERSION (if you have one).  IF YOU ARE ONLY UPDATING THE CHART, YOU ONLY NEED TO UPDATE THE `Chart.yaml` VERSION.  THERE MAY STILL BE INSTANCES WHEN THE PR WILL HAVE A GREEN BUILD, BUT WHEN THE PR IS MERGED THE BUILD WILL FAIL DUE TO ADDITIONAL CHECKS.  IF THE CONTAINER IMAGES DIFFER, THE BUILD LOGS WILL HAVE IN THEM A LIST OF WHAT CONTENTS HAVE CHANGED IN THE IMAGE, AND A NEW PR WILL NEED TO BE CREATED TO UPDATE BOTH `.version` AND THE `Chart.yaml` VERSION.

### Issues and Related PRs

LIST AND CHARACTERIZE RELATIONSHIP TO JIRA ISSUES AND OTHER PULL REQUESTS. BE SURE LIST DEPENDENCIES.

* Resolves CASM-XYZ
* Change will also be needed in <insert branch name here>
* Future work required by CASM-ABC
* Merge with <insert PR URL here>
* Merge before <insert PR URL here>
* Merge after <insert PR URL here>

### Testing

LIST THE ENVIRONMENTS IN WHICH THESE CHANGES WERE TESTED.

Tested on:

* <drink system>
* Craystack
* CMS base-box
* Virtual Shasta

Were the install/upgrade based validation checks/tests run?(goss tests/install-validation doc)
Was a fresh Install tested? Y/N   If not, Why?
Was an Upgrade tested?      Y/N   If not, Why?
Was a Downgrade tested?     Y/N.  If not, Why?
If schema changes were part of this change, how were those handled in your upgrade/downgrade testing?

WHAT WAS THE EXTENT OF TESTING PERFORMED? MANUAL VERSUS AUTOMATED TESTS (UNIT/SMOKE/OTHER)
HOW WERE CHANGES VERIFIED TO BE SUCCESSFUL?

### Risks and Mitigations

IF APPLICABLE, HAS A SECURITY AUDIT (via SNYK OR OTHERWISE) BEEN RUN?
ARE THERE KNOWN ISSUES WITH THESE CHANGES?
ANY OTHER SPECIAL CONSIDERATIONS?

INCLUDE THE FOLLOWING ITEMS THAT APPLY. LIST ADDITIONAL ITEMS AND PROVIDE MORE DETAILED INFORMATION AS APPROPRIATE.

Requires:

* Additional testing on bare-metal
* Compute nodes
* V1 system configuration (classic preview)
* V1 system configuration with SSDs
* V2 system configuration
* 3rd party software
* Broader integration testing
* Fresh install
* Platform upgrade
