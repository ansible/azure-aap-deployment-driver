import * as React from 'react';
import { useEffect } from 'react';
import { useDispatch } from "react-redux";
import axios from 'axios';
import { setDeploymentSteps } from '../redux/actions/deploymentActions';
import { DeploymentSteps } from './DeploymentStep';

export const Installer: React.FunctionComponent = () =>
{
  const dispatch = useDispatch();
  var [error, seterr] = React.useState(null)
  const fetchDeploymentSteps = async () => {
    const response = await axios
    .get("http://localhost:9090/step")
    .catch((err) => {
      seterr(err)
    });
    if (response) {
      seterr(null)
      dispatch(setDeploymentSteps(response,error))
    }
    else {
      dispatch(setDeploymentSteps([],error))
    }
  };

  useEffect(() => {
    const id = setInterval(fetchDeploymentSteps, 4000);
    return () => clearInterval(id);
 });

    return (
      <>
        <DeploymentSteps></DeploymentSteps>
      </>
  )}
