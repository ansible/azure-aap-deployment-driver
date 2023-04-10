import React from 'react';
import { RestartStep } from './RestartStep';
import { ProgressBar } from './ProgressBar';
import { CancelDeployment } from './CancelDeployment';
import { DeploymentProgressData} from '../../apis/types';
import './DeploymentProgress.css'
import { PageSection, Bullseye, Stack, StackItem, Text, TextVariants } from '@patternfly/react-core';

interface IDeploymentProgressProps {
  progressData: DeploymentProgressData
}


export const DeploymentProgress = ({ progressData }: IDeploymentProgressProps ) => {

  // set a deployment message based on the status
  let deploymentMessage = ""
  const cleanupMessage = "You still need to delete the managed application from your Azure subscription. " +
  "For more information about deleting resources, refer to the documentation linked to the left."
  if (progressData.isComplete) {
    deploymentMessage = "Your Ansible Automation Platform deployment is now complete."
  } else if (progressData.isCanceled) {
    deploymentMessage = `Your Ansible on Azure deployment is cancelled. ${cleanupMessage}`
  } else if (progressData.isPermanentlyFailed) {
    deploymentMessage = `The maximum number of retries has been reached and your Ansible on Azure deployment has failed. Please reinstall. ${cleanupMessage}`
  } else if (progressData.failedStepIds.length > 0) {
    deploymentMessage = `Deployment step "${progressData.failedStepNames[0]}" failed. Press the Restart button below to restart it.`
  } else {
    deploymentMessage = ""
  }

  // render restart for the first failed step
  const restartStep = (progressData.failedStepIds.length > 0 && !progressData.isCanceled && !progressData.isPermanentlyFailed ?
    <RestartStep stepExId={progressData.failedExId} /> :
    <></>
  )

  // render progress bar only if no failed steps
  const progressBar = (progressData.failedStepIds.length === 0 && !progressData.isPermanentlyFailed ?
    <ProgressBar progressPercent={progressData.progress} ></ProgressBar> :
    <></>
  )

  const cancelButton = ( !progressData.isComplete && !progressData.isCanceled && !progressData.isPermanentlyFailed ? <CancelDeployment/> : <></>)

  return (
    <>
      <PageSection>
        <Bullseye>
          <Stack hasGutter className='deployProgress'>
            { deploymentMessage.length > 0 && <StackItem><Text component={TextVariants.h1}>{ deploymentMessage }</Text></StackItem>}
            { restartStep }
            { progressBar }
            { cancelButton }
          </Stack>
        </Bullseye>
      </PageSection>
    </>
  )
}
