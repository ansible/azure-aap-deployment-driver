import React from 'react';
import { Page, PageHeader } from '@patternfly/react-core';
import logo from '../bgimages/Ansible.svg';
import { Outlet } from 'react-router-dom';

const PlainLayout = () => {
  return (
    <Page
      mainContainerId="primary-app-container"
      header={<PageHeader logo={<img src={logo} alt="Ansible Logo" />} />}
      className="pf-m-full-height" >
      <Outlet />
    </Page>
  );
};

export { PlainLayout };
