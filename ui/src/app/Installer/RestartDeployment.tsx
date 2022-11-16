import * as React from 'react';
import {
    Button,
  } from '@patternfly/react-core';

function handleClick(id) {
    fetch(`http://127.0.0.1:9090/execution/${id}/restart`, {
      method: 'POST',
      mode: 'cors',
      body: JSON.stringify(jsonData)
    })
   }

  export const RestartDeployment = (props) => {
    return(
      <>
        <div>
        {props.data['executions'][props.data['executions'].length-1]['provisioningState'] === 'Failed' ? <Button variant="primary" onClick={() => handleClick(props.data['ID'])}>Retry step</Button>: <></>}
        </div>
      </>
    );
  };

function jsonData(jsonData: any): BodyInit | null | undefined {
    throw new Error('Function not implemented.');
}
