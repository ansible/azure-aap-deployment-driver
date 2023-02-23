import React from 'react';
import { FlexItem, Icon, TextVariants, Text, Tooltip, Flex } from '@patternfly/react-core';
import CheckCircleIcon from '@patternfly/react-icons/dist/esm/icons/check-circle-icon';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import { DeploymentStepStatusData, StepStatuses } from '../../apis/types';

interface IStepStatusProps {
  stepStatusData: DeploymentStepStatusData
}

export const StepStatus = ({ stepStatusData }: IStepStatusProps) => {

  const hasStarted = (stepStatusData.status === StepStatuses.STARTED)
  const hasSucceeded = (stepStatusData.status === StepStatuses.SUCCEEDED)
  const hasFailed = (stepStatusData.status === StepStatuses.FAILED)
  const hasBeenCanceled = (stepStatusData.status === StepStatuses.CANCELED)

  return (
    <Flex align={{ default: 'alignRight' }} className='deployment-info'>
      <FlexItem>
        { (hasSucceeded || hasFailed) && <Text className='timeTaken' component={TextVariants.h5}>({stepStatusData.duration})</Text> }
      </FlexItem>
      <FlexItem>
        { hasFailed && <Text className='attempt' component={TextVariants.h5}>({stepStatusData.attempts} attempts)</Text> }
      </FlexItem>
      <FlexItem className="statusTooltip" align={{ default: 'alignRight' }}>
        { hasBeenCanceled && <Tooltip removeFindDomNode={true} content={<div>{"Cancelled"}</div>}><Icon className='icon1' status="warning"><ExclamationCircleIcon /></Icon></Tooltip> }
        { hasSucceeded && <Tooltip removeFindDomNode={true} content={<div>Success</div>}><Icon className='icon1' status="success"><CheckCircleIcon /></Icon></Tooltip> }
        { hasFailed && <Tooltip removeFindDomNode={true} content={stepStatusData.error}><Icon className='icon1' status="danger"><ExclamationCircleIcon /></Icon></Tooltip> }
        { hasStarted && <Icon className='icon1' isInProgress={true}><CheckCircleIcon /></Icon> }
      </FlexItem>
    </Flex>
  )
};
