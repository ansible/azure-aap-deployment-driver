/* eslint-disable @typescript-eslint/no-unused-vars */
import * as React from 'react';
import {
    Button,
    Card,
    CardBody,
    CardTitle,
    Progress,
    Text
  } from '@patternfly/react-core';

  export const ProgressBar = (props) =>
  {

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
                  <div>
        {props.data1 === 100 ? <Text className="SuccessMessage" >Your Ansible Automation Platform deployment is now complete.</Text>: <Button className='cancelButton' variant="secondary" onClick={() => handleClick(props.data['ID'])}>Cancel Deployment</Button>}
        </div>
            </CardTitle>
            <CardBody>
              <br/>
            </CardBody>
          </Card>
    </>
  )};

function jsonData(jsonData: any): BodyInit | null | undefined {
  throw new Error('Function not implemented.');
}
