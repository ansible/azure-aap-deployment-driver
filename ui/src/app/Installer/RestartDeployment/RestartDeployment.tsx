import * as React from 'react';
import {
  Button,
} from '@patternfly/react-core';
import { ApiService } from 'src/Services/apiService';
const apiService = new ApiService();
function handleClick(id) {
    apiService.restartStep(id);
}

export const RestartDeployment = (props) => {
  return (
    <>
      <div>
        {props.data['executions'][props.data['executions'].length - 1]['provisioningState'] === 'Failed' ? <Button variant="primary" onClick={() => handleClick(props.data['ID'])}>Retry step</Button> : <></>}
      </div>
    </>
  );
};
