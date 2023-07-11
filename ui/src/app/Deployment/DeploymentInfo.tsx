import React from 'react';
import { Alert, AlertActionLink } from '@patternfly/react-core';

interface IDeploymentInfoProps {
  closeHandler: () => void;
}

export const DeploymentInfo = ({closeHandler}: IDeploymentInfoProps) => {

	return (
		//<Alert variant="info" title="Info alert title" />
		<Alert
      isInline
      variant="info"
      title="Ansible Automation Platform requirements"
      //actionClose={<AlertActionCloseButton onClose={} />}
      actionLinks={
        <React.Fragment>
          {/* <AlertActionLink onClick={() => alert('Clicked on View details')}>Go to console.redhat.com</AlertActionLink> */}
          <AlertActionLink onClick={closeHandler}>Dismiss, I already have it all setup</AlertActionLink>
        </React.Fragment>
      }
    >
      <p>To use Ansible Automation Platform on Azure, you must have a valid subscription for Ansible Automation Platform in your Red Hat account.</p>
      <p>You can click on <a href='https://console.redhat.com' target='_blank' rel="noreferrer">this link</a> to setup your Ansible Automation Platform subscription as well as your Red Hat, account if needed.</p>
    </Alert>
	)
}
