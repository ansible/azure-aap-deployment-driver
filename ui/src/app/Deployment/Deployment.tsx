import React, { useState, useEffect } from 'react';
import { DeploymentSteps } from './DeploymentSteps/Steps';
import { DeploymentProgress } from './DeploymentProgress/DeploymentProgress';
import { getSteps } from '../apis/deployment';
import { getEntitlementsCount } from '../apis/entitlements';
import { DeploymentStepData, DeploymentProgressData, EntitlementsCount } from '../apis/types';
import { RHLoginModal } from './RHLoginModal';
import { EntitlementsInfo } from './EntitlementsInfo';

const SHOW_RH_LOGIN_SESSION_STORAGE_KEY = "showRHLogin"

export const Deployment = () => {

  const [stepsData, setStepsData] = useState<DeploymentStepData[]>()
  const [progressData, setProgressData] = useState<DeploymentProgressData>()
  const [showRHLogin, setShowRHLogin] = useState<boolean>(()=>{
    // this function gets initial state value from session storage
    const storageItem = sessionStorage.getItem(SHOW_RH_LOGIN_SESSION_STORAGE_KEY);
    if (storageItem === null) {
      // store initial value of true and return it
      sessionStorage.setItem(SHOW_RH_LOGIN_SESSION_STORAGE_KEY,String(true));
      return true;
    } else {
      // return stored value
      return storageItem.toLowerCase() === 'true'
    }
  })
  const [entitlementsCount, setEntitlementsCount] = useState<EntitlementsCount>()

  const fetchData = async () => {
    try {
      const data = await getSteps()
      setStepsData(data.steps)
      setProgressData(data.progress)
    } catch (error) {
      console.log("Could not fetch steps data.", error)
    }
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
    // only supporting changing the value to false
    if (showRHLogin === false) {
      sessionStorage.setItem(SHOW_RH_LOGIN_SESSION_STORAGE_KEY,String(false));
    }
  },[showRHLogin])

  const handleRHLoginModalAction = (loginOpened: boolean) => {
    //place to add more logic if needed
    setShowRHLogin(false)
  }

  useEffect(()=>{
    fetchEntitlementsData()
  },[])

  return (
    <>
      { <RHLoginModal isModalShown={showRHLogin} actionHandler={handleRHLoginModalAction}/> }
      {/* TODO Add some place holder for case when data is not available */}

      {entitlementsCount && <EntitlementsInfo entitlementsCount={entitlementsCount}></EntitlementsInfo> }

      {stepsData &&<DeploymentSteps stepsData={stepsData}></DeploymentSteps>}
      {progressData  && <DeploymentProgress progressData={progressData}></DeploymentProgress>}
    </>
  )
}
