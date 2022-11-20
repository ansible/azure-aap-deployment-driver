import * as React from 'react';
import { useEffect } from 'react';
import { useDispatch } from "react-redux";
import axios from 'axios';
import { setDeploymentSteps } from '../redux/actions/deploymentActions';
import { DeploymentSteps } from './DeploymentStep';

export const Installer: React.FunctionComponent = () =>
{
  const dispatch = useDispatch();
  const fetchProducts = async () => {
    const response = await axios
    .get("http://localhost:9090/step")
    .catch((err) => {
    });
    dispatch(setDeploymentSteps(response))
  };

  useEffect(() => {
    fetchProducts();
  },[])

  return (
  <>
    <DeploymentSteps></DeploymentSteps>
  </>
)};