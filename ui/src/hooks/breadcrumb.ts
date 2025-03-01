import { useBreadcrumbs } from '../context/BreadcrumbContext';

export const useSetBreadcrumbs = () => {
  const { setBreadcrumbs } = useBreadcrumbs();
  return setBreadcrumbs;
};
