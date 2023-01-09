import * as React from 'react';
import {
  Bullseye, Stack, StackItem, Icon, List, ListItem, PageSection, PageSectionVariants, TextContent, Text, Title, Tooltip, TextVariants, Flex, FlexItem,
} from '@patternfly/react-core';
import { RestartDeployment } from '../RestartDeployment/RestartDeployment';
import CheckCircleIcon from '@patternfly/react-icons/dist/esm/icons/check-circle-icon';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import { ProgressBar } from '../Progressbar/ProgressBar';
import './DeploymentStep.css'
import { ApiService } from 'src/Services/apiService';

export const DeploymentSteps = () => {
  var percent = 0
  var dataLength = 0
  const apiService = new ApiService();
  const [deploymentSteps, setList] = React.useState([] as any);
  React.useEffect(() => {
    let mounted = true;
    apiService.getSteps()
      .then(items => {
        if (mounted) {
          setList(items)
        }
      })
    return () => {
      mounted = false;
    };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  dataLength = deploymentSteps.length;
  const ProgressChangeHandler = (items) => {
    percent = percent + 100 / dataLength;
    return <Text className='timeTaken' component={TextVariants.h5}>({items['executions'][0]['duration']})</Text>
  };

  return (
    <>
      <PageSection variant={PageSectionVariants.light}>
        <TextContent>
          <Text component="h1">Ansible Automation Platform Installer</Text>
        </TextContent>
      </PageSection>
      <PageSection>
        <Bullseye>
          <Stack hasGutter className='deploymentStepsCont'>
            <StackItem isFilled>
              <Title headingLevel="h2">
                Deployment Steps
              </Title>
              <div className='deploy-step pf-u-box-shadow-md'>
                <List isPlain isBordered >
                  {deploymentSteps == null ? <Text component="h1">Unable to reach servers</Text> : <>
                    {deploymentSteps?.map(items => (
                      <ListItem className='service-list pf-u-box-shadow-md' key={items['ID'].toString()}>
                        <Flex className='step-name'>{items['name']}
                          {items['executions'].length ? <FlexItem align={{ default: 'alignRight' }} className='deployment-info'>
                            {items['executions'][items['executions'].length - 1]['provisioningState'] === 'Succeeded' ? ProgressChangeHandler(items) : <></>}
                            {items['executions'][items['executions'].length - 1]['provisioningState'] === 'Succeeded' ? <Tooltip
                              content={
                                <div>Success</div>
                              }
                            ><Icon className='icon1' status="success">
                                <CheckCircleIcon />
                              </Icon></Tooltip> : <Tooltip
                                content={
                                  <div>{items['executions'][items['executions'].length - 1]['message']}</div>
                                }
                              ><Icon className='icon1' status="warning">
                                <ExclamationCircleIcon />
                              </Icon></Tooltip>}
                            {items['executions'][items['executions'].length - 1]['provisioningState'] === 'Failed' ? <Text className='attempt' component={TextVariants.h5}>({items['executions'].length} attempts)</Text> : <></>}
                            <RestartDeployment data={items}></RestartDeployment>
                          </FlexItem> : <></>}
                        </Flex>
                      </ListItem>
                    ))} </>}
                </List>
              </div>
            </StackItem>
            <ProgressBar data={deploymentSteps} data1={percent}></ProgressBar>
          </Stack>
        </Bullseye>
      </PageSection>
    </>
  )
};
