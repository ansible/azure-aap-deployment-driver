/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react';
import { Progress, StackItem } from '@patternfly/react-core';
import './ProgressBar.css'

interface IProgressBarProps {
  progressPercent: number
}

export const ProgressBar = ({ progressPercent }: IProgressBarProps) => {
  return (
    <StackItem className='progress'>
      <Progress className='infoText' value={progressPercent} title="Overall progress" />
    </StackItem>
  )
};


