import React, { useState, useEffect } from 'react';
import { DeploymentSteps } from './DeploymentSteps/Steps';
import { DeploymentProgress } from './DeploymentProgress/DeploymentProgress';
import { getSteps } from '../apis/deployment';
import { getEntitlementsCount } from '../apis/entitlements';
import { DeploymentStepData, DeploymentProgressData, EntitlementsCount } from '../apis/types';
import { RHLoginModal } from './RHLoginModal';
import { EntitlementsInfo } from './EntitlementsInfo';


interface IDeploymentProps {
  showLoginDialog: boolean
}

export const Deployment = ({showLoginDialog}:IDeploymentProps) => {

  const [stepsData, setStepsData] = useState<DeploymentStepData[]>()
  const [progressData, setProgressData] = useState<DeploymentProgressData>()
  const [entitlementsCount, setEntitlementsCount] = useState<EntitlementsCount>()

  const fetchData = () => {
    // intentionally wrapped inside IIFE to make fetchData function have "void" return
    // because analysis tools need it to correctly use it as a callback in setInterval
    (async ()=>{
      try {
        const data = await getSteps()
        setStepsData(data.steps)
        setProgressData(data.progress)
      } catch (error) {
        console.log("Could not fetch steps data.", error)
      }
    })();
  }

  const fetchEntitlementsData = async () => {
    try {
      const entitlementsData = await getEntitlementsCount()
      setEntitlementsCount(entitlementsData)
    } catch(error) {
      console.log("Could not fetch entitlements data.", error)
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

  useEffect(()=>{
    fetchEntitlementsData()
  },[])

  return (
    <>
      { <RHLoginModal isModalShown={showLoginDialog}/> }
      {/* TODO Add some place holder for case when data is not available */}

      {entitlementsCount && <EntitlementsInfo entitlementsCount={entitlementsCount}></EntitlementsInfo> }

      {stepsData &&<DeploymentSteps stepsData={stepsData}></DeploymentSteps>}
      {progressData  && <DeploymentProgress progressData={progressData}></DeploymentProgress>}
    </>
  )
}
