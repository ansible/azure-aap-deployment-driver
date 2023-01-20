import * as React from 'react';
import { ExclamationTriangleIcon } from '@patternfly/react-icons';
import {
  PageSection,
  Title,
  EmptyState,
  EmptyStateIcon,
  EmptyStateBody,
} from '@patternfly/react-core';

export const NotAvailableYet = () => {

  return (
    <PageSection>
    <EmptyState variant="full" isFullHeight={true}>
      <EmptyStateIcon icon={ExclamationTriangleIcon} />
      <Title headingLevel="h1" size="lg">
      Ansible Automation Platform Deployment driver web UI not available yet.
      </Title>
      <EmptyStateBody>
        The web user interface is not available yet, please see the progress of the installation on your Microsoft Azure portal.
      </EmptyStateBody>
    </EmptyState>
    </PageSection>
  )
};
