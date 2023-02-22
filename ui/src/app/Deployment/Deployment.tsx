import React, { useState, useEffect } from 'react';
import { DeploymentSteps } from './DeploymentSteps/Steps';
import { DeploymentProgress } from './DeploymentProgress/DeploymentProgress';
import { getSteps } from '../apis/deployment';
import { DeploymentStepData, DeploymentProgressData } from '../apis/types';

export const Deployment = () => {

  const [stepsData, setStepsData] = useState<DeploymentStepData[]>()
  const [progressData, setProgressData] = useState<DeploymentProgressData>()

  const fetchData = async () => {
    try {
      const data = await getSteps()
      setStepsData(data.steps)
      setProgressData(data.progress)
    } catch (error) {
      console.log("Could not fetch steps data.", error)
    }

  }

  useEffect(() => {
    fetchData()
    const intervalId = setInterval(fetchData, 3000)
    // returning function to clear interval
    return () => {
      clearInterval(intervalId)
    }
  }, [])

  return (
    <>
      {/* TODO Add some place holder for case when data is not available */}
      {stepsData &&<DeploymentSteps stepsData={stepsData}></DeploymentSteps>}
      {progressData  && <DeploymentProgress progressData={progressData}></DeploymentProgress>}
    </>
  )
}
