import * as React from 'react';
import {
  Bullseye, Stack, StackItem, Icon, List, ListItem, PageSection, PageSectionVariants, TextContent, Text, Title, Tooltip, TextVariants, Flex, FlexItem,
} from '@patternfly/react-core';
import { RestartDeployment } from '../RestartDeployment/RestartDeployment';
import CheckCircleIcon from '@patternfly/react-icons/dist/esm/icons/check-circle-icon';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import { useSelector } from 'react-redux';
import { RootState } from '@app/store/reducers';
import { ProgressBar } from '../Progressbar/ProgressBar';
import './DeploymentStep.css'

export const DeploymentSteps = () => {
  const deploymentSteps = useSelector((state: RootState) => state.deployment.deploymentSteps);
  const error = useSelector((state: RootState) => state.deployment.err);
  var percent = 0
  var dataLength = 0
  if (deploymentSteps) {
    dataLength = deploymentSteps.length
  }
  const ProgressChangeHandler = (data) => {
    percent = percent + 100 / dataLength;
    return <Text className='timeTaken' component={TextVariants.h5}>({data['executions'][0]['duration']})</Text>
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
                  {error != null ? <Text component="h1">Unable to reach servers</Text> : <>
                    {deploymentSteps?.map(data => (
                      <ListItem className='service-list pf-u-box-shadow-md'>
                        <Flex className='step-name'>{data['name']}
                          {data['executions'].length ? <FlexItem align={{ default: 'alignRight' }} className='deployment-info'>
                            {data['executions'][data['executions'].length - 1]['provisioningState'] === 'Succeeded' ? ProgressChangeHandler(data) : <></>}
                            {data['executions'][data['executions'].length - 1]['provisioningState'] === 'Succeeded' ? <Tooltip
                              content={
                                <div>Success</div>
                              }
                            ><Icon className='icon1' status="success">
                                <CheckCircleIcon />
                              </Icon></Tooltip> : <Tooltip
                                content={
                                  <div>{data['executions'][data['executions'].length - 1]['message']}</div>
                                }
                              ><Icon className='icon1' status="warning">
                                <ExclamationCircleIcon />
                              </Icon></Tooltip>}
                            {data['executions'][data['executions'].length - 1]['provisioningState'] === 'Failed' ? <Text className='attempt' component={TextVariants.h5}>({data['executions'].length} attempts)</Text> : <></>}
                            <RestartDeployment data={data}></RestartDeployment>
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
      <PageSection>
      </PageSection>
    </>
  )
};
