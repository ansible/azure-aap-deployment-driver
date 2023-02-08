/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react';
import { Progress, StackItem, Text } from '@patternfly/react-core';
import './ProgressBar.css'

interface IProgressBarProps {
  progressPercent: number
  isComplete: boolean
}

export const ProgressBar = ({ progressPercent, isComplete }: IProgressBarProps) => {

  return (
    <StackItem className='progress'>
      <Progress value={progressPercent} title="Overall progress" />
      <br></br>
      <div>
        {isComplete ? <Text className="SuccessMessage" >Your Ansible Automation Platform deployment is now complete.</Text> :
          <></>}
      </div>
    </StackItem>
  )
};


