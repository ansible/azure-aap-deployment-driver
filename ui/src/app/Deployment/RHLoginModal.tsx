import React from 'react';
import { Button, Icon, Modal, ModalVariant, Title, TitleSizes } from '@patternfly/react-core';
import ExternalLinkSquareAltIcon from '@patternfly/react-icons/dist/esm/icons/external-link-square-alt-icon';
import ExclamationTriangleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-triangle-icon';
import './RHLoginModal.css';

interface IRHLoginProps {
	isModalShown: boolean
}

const DialogHeader = <Title headingLevel="h1" size={TitleSizes['2xl']}>
	<Icon status="warning" isInline><ExclamationTriangleIcon /></Icon>A valid subscription for Ansible Automation Platform in your Red Hat account is required
</Title>;

export const RHLoginModal = ({isModalShown}: IRHLoginProps) => {
	return (
		<Modal
				header={DialogHeader}
        isOpen={isModalShown}
        showClose={false}
				variant={ModalVariant.medium}
				actions={[
					// TODO Set the URL in href below programmatically
          <Button
					key="login" variant="primary"
					icon={<ExternalLinkSquareAltIcon />} iconPosition="right"
					component="a" href="/sso" target="_self">Log in with Red Hat account</Button>,
        ]}
      ><p>Your Ansible Automation Platform deployment is underway.</p>
			<br /><p>To use Ansible Automation Platform on Azure, you MUST have a valid subscription for Ansible Automation Platform in your Red Hat account.</p>
			<br />
			<p>You can set up your Ansible Automation Platform subscription and your Red Hat account by clicking the button below. You will be redirected back to this page upon successful log in or account creation.</p>
		</Modal>
	)
}
