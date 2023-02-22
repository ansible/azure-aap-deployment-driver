import * as React from 'react';
import { Bullseye, Stack, StackItem, PageSection, PageSectionVariants, TextContent, Text, Title, List } from '@patternfly/react-core';
import { DeploymentStep } from "./Step";
import { DeploymentStepData } from '@app/apis/types';

import './Steps.css'

interface IDeploymentStepsProps {
  stepsData: DeploymentStepData[]
}


export const DeploymentSteps = ({ stepsData }: IDeploymentStepsProps, ) => {
  return (
    <>
      <PageSection variant={PageSectionVariants.light}>
        <TextContent>
          <Text component="h1">Ansible Automation Platform Deployment Engine</Text>
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
                  {stepsData?.map(stepData => (<DeploymentStep key={stepData.id} stepData={stepData}></DeploymentStep>))}
                </List>
              </div>
            </StackItem>
          </Stack>
        </Bullseye>
      </PageSection>
    </>
  )
};
