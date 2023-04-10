import React from 'react';
import { Popover, Icon, Tooltip } from '@patternfly/react-core';
import { TableComposable, Tr, Th, Td, Tbody } from '@patternfly/react-table';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import { DeploymentStepStatusData } from '@app/apis/types';

import './ErrorInfo.css';

interface IErrorInfoProps {
  stepStatusData: DeploymentStepStatusData;
}

export const ErrorInfoPopover = ({ stepStatusData }: IErrorInfoProps) => {
  const data = (
      <TableComposable>
        <Tbody>
          <Tr>
            <Th>Error:</Th>
            <Td className="errorCell">{stepStatusData.error}</Td>
          </Tr>
          <Tr>
            <Th>Details:</Th>
            <Td className="errorDetailCell">
              <div className="errorDetailDiv">{stepStatusData.errorDetails}</div>
            </Td>
          </Tr>
          <Tr>
            <Th modifier="nowrap">Correlation ID:</Th>
            <Td>{stepStatusData.correlationId}</Td>
          </Tr>
        </Tbody>
      </TableComposable>
  );
  return (
    <Tooltip content={"Click for error info"}>
      <Popover minWidth="40rem" bodyContent={data} showClose={false} >
        <Icon className="icon1" status="danger">
          <ExclamationCircleIcon />
        </Icon>
      </Popover>
    </Tooltip>
  );
};
