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
		"message": "To manage or setup new subscription, visit the", // followed by a link
		"linkTitle": "Red Hat Hybrid Cloud Console",
		"linkURL":"https://console.redhat.com/",
		"linkInLine": true
	},
	"pending": {
		"variant": AlertVariant.info,
		"title": "Your Ansible Automation Platform subscription is pending",
		"message": "You do not have an Ansible Automation Platform subscription or your subscription is being entitled. " +
			"Click on the following link to enable your Ansible Automation Platform subscription and to access Red Hat support. This is a required step in order to access the Ansible Automation Platform once it is deployed.", //followed by a link
		"linkTitle": "Red Hat Hybrid Cloud Console",
		"linkURL":"https://console.redhat.com/?azure-ansible-activation",
		"linkInLine": false
	},
	"error": {
		"variant": AlertVariant.danger,
		"title": "We're temporarily unable to fetch your subscription information",
		"message": "Click on the following link to enable your Ansible Automation Platform subscription and to access Red Hat support. This is a required step in order to access the Ansible Automation Platform once it is deployed.", //followed by a link
		"linkTitle": "Red Hat Hybrid Cloud Console",
		"linkURL":"https://console.redhat.com/",
		"linkInLine": false
	}
}

export const EntitlementsInfo = ({entitlementsCount}:IEntitlementsInfoProps) => {
	let alertContent: any

	if (entitlementsCount.error) {
		alertContent = AlertContent.error
	} else {
		alertContent = (entitlementsCount.count > 0) ? AlertContent.existing : AlertContent.pending
	}

	return(
		<Alert variant={alertContent.variant} isInline={true} title={alertContent.title} className='entitlements-info'>
			<Text component={TextVariants.p}>{alertContent.message} {!alertContent.linkInLine && <br/>}<Button variant="link" component="a" isInline={alertContent.linkInLine}
				icon={<ExternalLinkSquareAltIcon />} iconPosition="right"
				href={alertContent.linkURL} target="_blank">{alertContent.linkTitle}</Button>
			</Text>
		</Alert>
	)
}
