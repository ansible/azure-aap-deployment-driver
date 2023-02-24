import React from 'react';
import { Popover, Icon, Tooltip } from '@patternfly/react-core';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import { DeploymentStepStatusData } from '@app/apis/types';

interface IErrorInfoProps {
  stepStatusData: DeploymentStepStatusData;
}

export const ErrorInfoPopover = ({ stepStatusData }: IErrorInfoProps) => {
  const borders = {
    border: '1px solid rgb(0,0,0)',
  };

  let data = (
    <>
      <table>
        <tr style={borders}>
          <th>Error:</th>
          <td>{stepStatusData.error}</td>
        </tr>
        <tr style={borders}>
          <th>Details:</th>
          <td>
            <div style={{ height: '20rem', overflowY: 'scroll', whiteSpace: 'pre-wrap' }}>{stepStatusData.errorDetails}</div>
          </td>
        </tr>
        <tr style={borders}>
          <th style={{width: "8rem"}}>Correlation ID:</th>
          <td>{stepStatusData.correlationId}</td>
        </tr>
      </table>
    </>
  );
  return (
    <Tooltip content={<div>Click for error info</div>}>
      <Popover minWidth="40rem" bodyContent={<div>{data}</div>}>
        <Icon className="icon1" status="danger">
          <ExclamationCircleIcon />
        </Icon>
      </Popover>
    </Tooltip>
  );
};
