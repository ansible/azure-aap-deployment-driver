import * as React from 'react';
import { CubesIcon } from '@patternfly/react-icons';
import {
  PageSection,
  Title,
  Button,
  EmptyState,
  EmptyStateVariant,
  EmptyStateIcon,
  EmptyStateBody
} from '@patternfly/react-core';

export const Documentation = () => {
  return(
    <PageSection>
    <EmptyState variant={EmptyStateVariant.full}>
      <EmptyStateIcon icon={CubesIcon} />
      <Title headingLevel="h1" size="lg">
        Documentation
        </Title>
      <EmptyStateBody>
      Documentation links below:
        </EmptyStateBody>
        <Button component="a" href="https://docs.google.com/document/d/1jdr3QHJa8sMYzmYD2itliYtlhswjryhR2ZhaAZVBVOw/edit#heading=h.dzey4iizy1h" target="_blank" variant="primary">
      Link to core docs</Button>{' '}
    </EmptyState>
  </PageSection>
  )
}

