import React from 'react';
import { RestartStep } from './RestartStep';
import { ProgressBar } from './ProgressBar';
import { CancelDeployment } from './CancelDeployment';
import { DeploymentProgressData, DeploymentStepData} from '../../apis/types';
import './DeploymentProgress.css'
import { PageSection, Bullseye, Stack } from '@patternfly/react-core';

interface IDeploymentProgressProps {
  progressData: DeploymentProgressData
  stepsData: DeploymentStepData[]
}


export const DeploymentProgress = ({ progressData, stepsData }: IDeploymentProgressProps) => {

  // render restart for the first failed step
  const restartStep = (progressData.failedStepIds.length > 0 ?
    <RestartStep stepExId={progressData.failedExId} stepName={progressData.failedStepNames[0]} /> :
    <></>
  )

  // render progress bar only if no failed steps
  const progressBar = (progressData.failedStepIds.length < 0 ?
    <ProgressBar progressPercent={progressData.progress} isComplete={progressData.isComplete}></ProgressBar> :
    <></>
  )

  return (

    <div >
      <PageSection>
        <Bullseye>
          <Stack hasGutter className='deployProgress'>
            {restartStep}
            {progressBar}
            <CancelDeployment />
          </Stack>
        </Bullseye>
      </PageSection>
    </div>
  )
}
