import React from 'react';
import { Button, Modal, ModalVariant } from '@patternfly/react-core';
import ExternalLinkSquareAltIcon from '@patternfly/react-icons/dist/esm/icons/external-link-square-alt-icon';

interface IRHLoginProps {
	isModalShown: boolean
	actionHandler: (boolean) => void;
}

export const RHLoginModal = ({isModalShown, actionHandler}: IRHLoginProps) => {
	return (
		<Modal
        title="Ansible Automation Platform requirements"
				titleIconVariant="info"
        isOpen={isModalShown}
        showClose={false}
				variant={ModalVariant.medium}
				actions={[
          <Button
					key="login" variant="primary"
					icon={<ExternalLinkSquareAltIcon />} iconPosition="right"
					component="a" href="https://console.redhat.com/" target="_blank"  onClick={()=>{actionHandler(true)}}
					>Red Hat Hybrid Cloud Console</Button>,
          <Button
					key="dismiss" variant="link"
					onClick={()=>{actionHandler(false)}}
					>Dismiss: my account and subscription are set up</Button>
        ]}
      ><p>To use Ansible Automation Platform on Azure, you MUST have a valid subscription for Ansible Automation Platform in your Red Hat account.</p>
			<br />
			<p>You can set up your Ansible Automation Platform subscription and your Red Hat account by clicking the button below.</p>
		</Modal>
	)
}
