# Description
<!-- Describe the changes introduced in the PR below, include rationale and technical/design decisions -->
This PR ...


<!-- Uncomment following block and add links (surrounded by < and > ) to any resources or documents if those might help explain the PR -->
<!--
## Links
<...>
-->

## Commit message reminder

Commit messages are important because they will be used to build list of items included in Errata when the deployment driver container image is built and QE will use that list to figure out what to test before the image is published.

- For user facing features or fixes make sure Jira item number `AAP-nnnnn` is included in the commit message(s)
- For other changes, like unit tests or code refactoring, do not include Jira item numbers

## PR task checklist
<!-- Mark tasks done by putting x inside [ ]. If any task in the checklist does not apply, mark the task done AND surround the text with ~ ~ (tilda) to strike it out -->
- [ ] I have used appropriate commit messages based on section above
- [ ] I have performed a self-review of my code
- [ ] I have followed code style guidelines of this project
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] I have added thorough tests
- [ ] I have had a successful deployment with the changes in this PR

## Jira
<!--- Put JIRA story/task/bug number in the link below or remove the next line and uncomment one below it. -->
This PR is for JIRA item: <https://issues.redhat.com/browse/AAP-NNNN>
<!-- This PR does not need a corresponding JIRA item. -->

## Testing
<!-- Describe the testing process in set of steps. If testing is not applicable, remove the steps and add a statement explaining why testing isn't applicable. -->
### Steps to test
1. Pull down the PR
2. ...
3. ...

### Expected result
<!-- Describe expected results  -->
With the changes in this PR you will see ...

### Screenshots / Console output / Logs
<!-- Add screenshot and/or console output if applicable and uncomment applicable block below -->
<!--
Screenshots:

-->

<!--
Output or logs:
```sh
Raw text goes here
```
-->
