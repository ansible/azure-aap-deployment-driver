import React from 'react';
import { FlexItem, Icon, TextVariants, Text, Tooltip, Flex } from '@patternfly/react-core';
import CheckCircleIcon from '@patternfly/react-icons/dist/esm/icons/check-circle-icon';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import HistoryIcon from '@patternfly/react-icons/dist/esm/icons/history-icon';
import { DeploymentStepStatusData, StepStatuses } from '../../apis/types';
import { ErrorInfoPopover } from '../../ErrorInfo/ErrorInfo';

interface IStepStatusProps {
  stepStatusData: DeploymentStepStatusData
}

export const StepStatus = ({ stepStatusData }: IStepStatusProps) => {

  const hasStarted = (stepStatusData.status === StepStatuses.STARTED)
  const hasSucceeded = (stepStatusData.status === StepStatuses.SUCCEEDED)
  const hasFailed = (stepStatusData.status === StepStatuses.FAILED)
  const hasRestartPending = (stepStatusData.status === StepStatuses.RESTART_PENDING)
  const hasBeenCanceled = (stepStatusData.status === StepStatuses.CANCELED)

  return (
    <Flex align={{ default: 'alignRight' }} className='deployment-info'>
      <FlexItem>
        { (hasSucceeded || hasFailed) && <Text className='timeTaken' component={TextVariants.h5}>({stepStatusData.duration})</Text> }
      </FlexItem>
      <FlexItem>
        { (hasFailed || hasRestartPending) && <Text className='attempt' component={TextVariants.h5}>({stepStatusData.attempts} { stepStatusData.attempts === 1 ? 'attempt' : 'attempts' })</Text> }
      </FlexItem>
      <FlexItem className="statusTooltip" align={{ default: 'alignRight' }}>
        { hasBeenCanceled && <Tooltip removeFindDomNode={true} content={<div>{"Cancelled"}</div>}><Icon className='icon1' status="warning"><ExclamationCircleIcon /></Icon></Tooltip> }
        { hasRestartPending && <Tooltip removeFindDomNode={true} content={<div>{"Restart pending"}</div>}><Icon className='icon1' status="warning"><HistoryIcon /></Icon></Tooltip> }
        { hasSucceeded && <Tooltip removeFindDomNode={true} content={<div>Success</div>}><Icon className='icon1' status="success"><CheckCircleIcon /></Icon></Tooltip> }
        { hasFailed && <ErrorInfoPopover stepStatusData={stepStatusData}></ErrorInfoPopover> }
        { hasStarted && <Icon className='icon1' isInProgress={true}><CheckCircleIcon /></Icon> }
      </FlexItem>
    </Flex>
  )
};
