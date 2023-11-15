import React from 'react';
import { Alert, AlertVariant, Button, Text, TextVariants } from '@patternfly/react-core';
import ExternalLinkSquareAltIcon from '@patternfly/react-icons/dist/esm/icons/external-link-square-alt-icon';
import { EntitlementsCount } from '../apis/types';

import "./EntitlementsInfo.css"

interface IEntitlementsInfoProps {
	entitlementsCount: EntitlementsCount
}

const AlertContent = {
	"existing": {
		"variant": AlertVariant.info,
		"title":"You currently have a subscription to Ansible Automation Platform",
		"content": "To manage or setup new subscription, visit the" // followed by a link
	},
	"pending": {
		"variant": AlertVariant.info,
		"title": "Your Ansible Automation Platform subscription is pending",
		"content": "Your subscription is being entitled and deployed, and will be ready for use shortly. " +
		"In the meantime, you can manage your subscription from the" //followed by a link
	},
	"error": {
		"variant": AlertVariant.danger,
		"title": "We're temporarily unable to fetch your subscription information",
		"content": "In the meantime, you can manage your subscription from the" //followed by a link
	}
}

export const EntitlementsInfo = ({entitlementsCount}:IEntitlementsInfoProps) => {
	let alert: any

	if (entitlementsCount.error) {
		alert = AlertContent.error
	} else {
		alert = (entitlementsCount.count > 0) ? AlertContent.existing : AlertContent.pending
	}

	return(
		<Alert variant={alert.variant} isInline={true} title={alert.title} className='entitlements-info'>
			<Text component={TextVariants.p}>{alert.content} <Button variant="link" component="a" isInline
				icon={<ExternalLinkSquareAltIcon />} iconPosition="right"
				href="https://console.redhat.com/" target="_blank">Red Hat Hybrid Cloud Console</Button>
			</Text>
		</Alert>
	)
}
