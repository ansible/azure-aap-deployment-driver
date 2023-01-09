/* eslint-disable @typescript-eslint/no-unused-vars */
import * as React from 'react';
import {
    Button,
    Card,
    CardBody,
    CardTitle,
    Progress,
    StackItem,
    Text
  } from '@patternfly/react-core';
import { ApiService } from 'src/Services/apiService';
import './ProgressBar.css'

  const apiService = new ApiService();

  export const ProgressBar = (props) =>
  {
    function handleClick() {
      apiService.cancelDeployment();
   }

  return (
    <>
      <StackItem className='progress'>
        <Progress value={props.data1} title="Overall progress" />
        <br></br>
        <div>
          {props.data1 === 100 ? <Text className="SuccessMessage" >Your Ansible Automation Platform deployment is now complete.</Text> :
            <Button className='cancelButton' variant="secondary" onClick={() => handleClick()}>Cancel Deployment</Button>}
        </div>
      </StackItem>
    </>
  )};
