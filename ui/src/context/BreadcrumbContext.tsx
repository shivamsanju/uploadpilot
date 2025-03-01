import React, { createContext, useContext, useState } from 'react';

interface Breadcrumb {
  label: string;
  path?: string;
}

interface BreadcrumbsContextType {
  breadcrumbs: Breadcrumb[];
  setBreadcrumbs: (items: Breadcrumb[]) => void;
}

const BreadcrumbsContext = createContext<BreadcrumbsContextType | undefined>(
  undefined,
);

export const BreadcrumbsProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [breadcrumbs, setBreadcrumbs] = useState<Breadcrumb[]>([]);

  return (
    <BreadcrumbsContext.Provider value={{ breadcrumbs, setBreadcrumbs }}>
      {children}
    </BreadcrumbsContext.Provider>
  );
};

export const useBreadcrumbs = (): BreadcrumbsContextType => {
  const context = useContext(BreadcrumbsContext);
  if (!context) {
    throw new Error('useBreadcrumbs must be used within a BreadcrumbsProvider');
  }
  return context;
};
