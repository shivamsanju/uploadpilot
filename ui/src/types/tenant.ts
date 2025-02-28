export type TenantOnboardingRequest = {
  name: string;
  contactEmail: string;
  phone?: string;
  address?: string;
  industry?: string;
  companyName?: string;
  role: string;
};
