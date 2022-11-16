/* eslint-disable @typescript-eslint/no-unused-vars */
import * as React from 'react';
import {
    Button,
    Card,
    CardBody,
    CardTitle,
    Progress
  } from '@patternfly/react-core';

import { useState } from 'react';

  const options = {
    method : 'POST',
    headers : {
      'Content-Type' : 'application/json'
    }
  }

  export const ProgressBar = (props) =>
  {
  var percent = 30
  console.log(props)

  function handleClick(id) {
    fetch(`http://127.0.0.1:9090/execution/${id}/restart`, {
      method: 'POST',
      mode: 'cors',
      body: JSON.stringify(jsonData)
    })
   }

  return (
    <>

          <Card isHoverable isCompact style={{width:"203%"}}>
            <CardTitle>
                  <Progress value={props.data1} title="Overall progress" />
                  <br></br>
                  <Button className='cancleButton' variant="secondary" onClick={() => handleClick(props.data['ID'])}>Cancle Deployment</Button>
            </CardTitle>
            <CardBody>
              <br />
            </CardBody>
          </Card>
    </>
  )};

function jsonData(jsonData: any): BodyInit | null | undefined {
  throw new Error('Function not implemented.');
}
