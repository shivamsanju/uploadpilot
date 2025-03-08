import { Box } from '@mantine/core';
import { useEffect } from 'react';
import { PricingSection } from '../../components/Pricing';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';

const BillingsPage = () => {
  const setBreadcrumbs = useSetBreadcrumbs();

  useEffect(() => {
    setBreadcrumbs([{ label: 'Workspaces', path: '/' }, { label: 'Billing' }]);
  }, [setBreadcrumbs]);

  return (
    <Box mb={50}>
      <PricingSection />
    </Box>
  );
};

export default BillingsPage;
