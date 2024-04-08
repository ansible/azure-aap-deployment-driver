import React from 'react';
import { RestartStep } from './RestartStep';
import { ProgressBar } from './ProgressBar';
import { CancelDeployment } from './CancelDeployment';
import { DeploymentProgressData, EngineConfiguration} from '../../apis/types';
import './DeploymentProgress.css'
import { PageSection, Bullseye, Stack, StackItem, Text, TextVariants, Title, HelperText, HelperTextItem, Level, LevelItem } from '@patternfly/react-core';
import { EngineConfigurationInfo } from '../EngineConfigurationInfo';

interface IDeploymentProgressProps {
  progressData: DeploymentProgressData
  engineConfig: EngineConfiguration | undefined
}


export const DeploymentProgress = ({ progressData, engineConfig }: IDeploymentProgressProps ) => {

  // set a deployment message based on the status
  let progressTitle = "This is the status of your deployment and deployment engine configuration limits that could prevent your successful deployment."
  let deploymentMessage = ""
  let deploymentStatusIcon: "default"|"indeterminate"|"warning"|"success"|"error"|undefined = 'default'
  let deploymentCancelMessage = "You can cancel your deployment at any time."
  const cleanupMessage = "You still need to delete the managed application from your Azure subscription. " +
  "For more information about deleting resources, refer to the documentation linked to the left."
  if (progressData.isComplete) {
    progressTitle = "This is the status of your deployment and deployment engine configuration limits."
    deploymentMessage = "Your Ansible Automation Platform deployment is now complete."
    deploymentStatusIcon = "success"
    deploymentCancelMessage = ""
  } else if (progressData.isCanceled) {
    progressTitle = "This is the status of your deployment and deployment engine configuration limits."
    deploymentMessage = `Your Ansible on Azure deployment is cancelled. ${cleanupMessage}`
    deploymentStatusIcon = "default"
    deploymentCancelMessage = ""
  } else if (progressData.isPermanentlyFailed) {
    progressTitle = "This is the status of your deployment and deployment engine configuration limits."
    deploymentMessage = `The maximum number of retries has been reached and your Ansible on Azure deployment has failed. Please reinstall. ${cleanupMessage}`
    deploymentStatusIcon = "error"
    deploymentCancelMessage = ""
  } else if (progressData.failedStepIds.length > 0) {
    deploymentMessage = `Deployment step "${progressData.failedStepNames[0]}" failed. Press the Restart button below to restart it.`
    deploymentStatusIcon = 'warning'
  } else {
    deploymentMessage = ""
  }

  const deploymentStatus = <HelperTextItem hasIcon variant={deploymentStatusIcon}>{deploymentMessage}</HelperTextItem>

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
    <PageSection isFilled>
      <Bullseye>
        <Stack className='deploymentProgressContainer' >
          <StackItem>
            <Title headingLevel="h2">Deployment Status</Title>
          </StackItem>
          <StackItem className='deploymentProgress'>
            <Stack hasGutter>
              <StackItem>
                <Text component={TextVariants.h1}>{progressTitle} {deploymentCancelMessage}</Text>
              </StackItem>
              { deploymentMessage && <StackItem><HelperText>{ deploymentStatus }</HelperText></StackItem>}
              <StackItem>
                <Bullseye>
                  <Level hasGutter>
                    <LevelItem>{ restartStep }</LevelItem>
                    <LevelItem>{ cancelButton }</LevelItem>
                  </Level>
                </Bullseye>
              </StackItem>
              { progressBar }
              { engineConfig && <EngineConfigurationInfo engineConfig={engineConfig}></EngineConfigurationInfo>}
            </Stack>
          </StackItem>
        </Stack>
      </Bullseye>
    </PageSection>
  )
}
