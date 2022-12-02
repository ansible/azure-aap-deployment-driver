import * as React from 'react';
import { useEffect } from 'react';
import { useDispatch } from "react-redux";
import { setDeploymentSteps } from '../store/actions/deploymentActions';
import { DeploymentSteps } from './DeploymentStep';

export const Installer: React.FunctionComponent = () =>
{
  const dispatch = useDispatch();
  var [error, seterr] = React.useState(null)

  const fetchDeploymentSteps = async() => { const response = await fetch('http://localhost:9090/step', {
  }).then(response => {
    if(response.ok) {
      return response.json()
    }
    throw response.json();
  }).catch(error => {
    seterr(error);
  })
  if (response !== undefined) {
    seterr(null)
    dispatch(setDeploymentSteps(response,error))
  }
  else {
    dispatch(setDeploymentSteps(undefined,error))
  }};


  useEffect(() => {
    const id = setInterval(fetchDeploymentSteps, 4000);
    return () => clearInterval(id);
 });

    return (
      <>
        <DeploymentSteps></DeploymentSteps>
      </>
  )}
