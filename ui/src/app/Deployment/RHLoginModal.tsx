import React from 'react';
import { Button, Modal, ModalVariant } from '@patternfly/react-core';
import ExternalLinkSquareAltIcon from '@patternfly/react-icons/dist/esm/icons/external-link-square-alt-icon';

interface IRHLoginProps {
	isModalShown: boolean
}

export const RHLoginModal = ({isModalShown}: IRHLoginProps) => {
	return (
		<Modal
        title="A valid subscription for Ansible Automation Platform in your Red Hat account is required"
				titleIconVariant="warning"
        isOpen={isModalShown}
        showClose={false}
				variant={ModalVariant.medium}
				actions={[
					// TODO Set the URL in href below programmatically
          <Button
					key="login" variant="primary"
					icon={<ExternalLinkSquareAltIcon />} iconPosition="right"
					component="a" href="https://sso.redhat.com/" target="_blank">Log in with Red Hat account</Button>,
        ]}
      ><p>Your Ansible Automation Platform deployment is underway.</p>
			<br /><p>To use Ansible Automation Platform on Azure, you MUST have a valid subscription for Ansible Automation Platform in your Red Hat account.</p>
			<br />
			<p>You can set up your Ansible Automation Platform subscription and your Red Hat account by clicking the button below. You will be redirected back to this page upon successful log in or account creation.</p>
		</Modal>
	)
}
