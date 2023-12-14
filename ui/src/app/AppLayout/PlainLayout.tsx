import React from 'react';
import { Brand, Page, PageHeader } from '@patternfly/react-core';
import logo from '../bgimages/Technology_icon-Red_Hat-Ansible_Automation_Platform-Standard-RGB.svg';
import { Outlet } from 'react-router-dom';

const PlainLayout = () => {
  return (
    <Page
      mainContainerId="primary-app-container"
      header={<PageHeader logo={<Brand src={logo} alt="Red Hat Ansible Automation Platform Logo"/>} />}
      className="pf-m-full-height" >
      <Outlet />
    </Page>
  );
};

export { PlainLayout };
