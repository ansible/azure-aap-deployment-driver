import React, { useState, useEffect } from 'react';
import { Bullseye, Stack, StackItem, PageSection, PageSectionVariants, TextContent, Text, Title, List } from '@patternfly/react-core';
import { DeploymentStep } from "./Step";
import { DeploymentStepData } from '@app/apis/types';

import './Steps.css'
import { DeploymentInfo } from '../DeploymentInfo';

interface IDeploymentStepsProps {
  stepsData: DeploymentStepData[]
}


export const DeploymentSteps = ({ stepsData }: IDeploymentStepsProps, ) => {

  // Persist this value browser's session (each tab has its own)
  const [showDeploymentInfo, setShowDeploymentInfo] = useState<Boolean>(()=>{
    let storageItem = sessionStorage.getItem('showDeploymentInfo');
    if (storageItem === null) {
      // store initial value of true and return it
      sessionStorage.setItem('showDeploymentInfo',String(true));
      return true;
    } else {
      // return stored value
      return storageItem.toLowerCase() === 'true'
    }
  })

  useEffect(()=>{
    // update the value in browser session when changed
    sessionStorage.setItem('showDeploymentInfo',String(showDeploymentInfo));
  },[showDeploymentInfo])



  const closeDeploymentInfo = () => {
    setShowDeploymentInfo(!showDeploymentInfo);
  }

  return (
    <>
      <PageSection variant={PageSectionVariants.light}>
        <TextContent>
          <Text component="h1">Ansible Automation Platform Deployment Engine</Text>
        </TextContent>
      </PageSection>
      {showDeploymentInfo && <DeploymentInfo closeHandler={closeDeploymentInfo}/>}
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
