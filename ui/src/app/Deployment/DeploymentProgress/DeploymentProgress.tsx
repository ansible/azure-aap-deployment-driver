import React, { Dispatch, SetStateAction } from 'react';
import { RestartStep } from './RestartStep';
import { ProgressBar } from './ProgressBar';
import { CancelDeployment } from './CancelDeployment';
import { DeploymentProgressData} from '../../apis/types';
import './DeploymentProgress.css'
import { PageSection, Bullseye, Stack } from '@patternfly/react-core';

interface IDeploymentProgressProps {
  progressData: DeploymentProgressData
  setCancelled: Dispatch<SetStateAction<boolean>>
}


export const DeploymentProgress = ({ progressData, setCancelled}: IDeploymentProgressProps ) => {

  // render restart for the first failed step
  const restartStep = (progressData.failedStepIds.length > 0 ?
    <RestartStep stepExId={progressData.failedExId} stepName={progressData.failedStepNames[0]} /> :
    <></>
  )

  // render progress bar only if no failed steps
  const progressBar = (progressData.failedStepIds.length === 0 ?
    <ProgressBar progressPercent={progressData.progress} isComplete={progressData.isComplete}></ProgressBar> :
    <></>
  )

  const cancelButton = (!progressData.isComplete ? <CancelDeployment setCancelled={setCancelled}/> : <></>)

  return (

    <div >
      <PageSection>
        <Bullseye>
          <Stack hasGutter className='deployProgress'>
            {restartStep}
            {progressBar}
            {cancelButton}
          </Stack>
        </Bullseye>
      </PageSection>
    </div>
  )
}
