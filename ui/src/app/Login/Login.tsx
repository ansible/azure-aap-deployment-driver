import React from "react";
import { FormLogin } from "./LoginForm";
import {
	LoginPage,
	ListVariant
} from '@patternfly/react-core';

export function Login() {


	const images = {
		lg: '/assets/images/pfbg_1200.jpg',
		sm: '/assets/images/pfbg_768.jpg',
		sm2x: '/assets/images/pfbg_768@2x.jpg',
		xs: '/assets/images/pfbg_576.jpg',
		xs2x: '/assets/images/pfbg_576@2x.jpg'
	};

	return (
		<LoginPage
			footerListVariants={ListVariant.inline}
			backgroundImgSrc={images}
			loginTitle="Deployment Engine"
			loginSubtitle="Please use the administrative credentials for Red Hat Ansible Automation Platform on Microsoft Azure.">
			<FormLogin />
		</LoginPage>
	)
}
