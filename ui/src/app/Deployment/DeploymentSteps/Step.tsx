import React from 'react';
import { ListItem, Flex } from '@patternfly/react-core';
import { StepStatus } from './StepStatus';
import { DeploymentStepData } from '../../apis/types';

interface IStepProps {
  stepData: DeploymentStepData
  isCancelled:boolean
}


export const DeploymentStep = ({ stepData, isCancelled }: IStepProps) => {

  return (
    <ListItem className='service-list pf-u-box-shadow-md'>
      <Flex className='step-name'>{stepData.name}
        {stepData.status && <StepStatus stepStatusData={stepData.status} isCancelled={isCancelled}/>}
      </Flex>
    </ListItem>
  )
};
