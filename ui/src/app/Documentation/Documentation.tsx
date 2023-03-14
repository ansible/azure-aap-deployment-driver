import * as React from 'react';
import { CubesIcon } from '@patternfly/react-icons';
import {
  PageSection,
  Title,
  Button,
  EmptyState,
  EmptyStateVariant,
  EmptyStateIcon,
  EmptyStateBody,
} from '@patternfly/react-core';

export const Documentation = () => {
  return (
    <PageSection>
      <EmptyState variant={EmptyStateVariant.full}>
        <EmptyStateIcon icon={CubesIcon} />
        <Title headingLevel="h1" size="lg">
          Ansible Automation Platform on Azure documentation
        </Title>
        <EmptyStateBody>(opens in a new tab)</EmptyStateBody>
        <Button
          component="a"
          href="https://access.redhat.com/documentation/en-us/ansible_on_clouds/2.x/html/red_hat_ansible_automation_platform_on_microsoft_azure_guide/index"
          target="_blank"
          variant="primary"
        >
          AAP on Azure Guide
        </Button>{' '}
      </EmptyState>
    </PageSection>
  );
};
