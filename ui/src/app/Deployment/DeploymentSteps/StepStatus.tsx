import * as React from 'react';
import { FlexItem, Icon, TextVariants, Text, Tooltip, Flex } from '@patternfly/react-core';
import CheckCircleIcon from '@patternfly/react-icons/dist/esm/icons/check-circle-icon';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import { DeploymentStepStatusData, StepStatuses } from '../../apis/types';

interface IStepStatusProps {
  stepStatusData: DeploymentStepStatusData
}

export const StepStatus = ({ stepStatusData }: IStepStatusProps) => {

  const startState = (stepStatusData.status === StepStatuses.STARTED ?
    <Icon className='icon1' isInProgress={true}><CheckCircleIcon /></Icon> :
    <></>
  );

  const stepDuration = ((stepStatusData.status === StepStatuses.SUCCEEDED || stepStatusData.status === StepStatuses.FAILED) ?
    <Text className='timeTaken' component={TextVariants.h5}>({stepStatusData.duration})</Text> :
    <></>
  );


  const attemptsText = (stepStatusData.status === StepStatuses.FAILED ?
    <Text className='attempt' component={TextVariants.h5}>({stepStatusData.attempts} attempts)</Text> :
    <></>
  );

  const statusTooltip = (stepStatusData.status === StepStatuses.SUCCEEDED ?
    <Tooltip content={<div>Success</div>}><Icon className='icon1' status="success"><CheckCircleIcon /></Icon></Tooltip> :
    stepStatusData.status === StepStatuses.FAILED ?
      <Tooltip content={<div>{stepStatusData.message}</div>}><Icon className='icon1' status="warning"><ExclamationCircleIcon /></Icon></Tooltip> :
      <></>
  )

  return (
    <Flex align={{ default: 'alignRight' }} className='deployment-info'>
      <FlexItem>{stepDuration}</FlexItem>
      <FlexItem>{attemptsText}</FlexItem>
      <FlexItem className="statusTooltip" align={{ default: 'alignRight' }} >{statusTooltip}{startState}</FlexItem>
    </Flex>
  )
};
