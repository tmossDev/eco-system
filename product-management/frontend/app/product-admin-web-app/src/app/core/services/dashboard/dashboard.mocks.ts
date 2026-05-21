import { DashboardSummary } from './dashboard.models';

export const MOCK_DASHBOARD_SUMMARY: DashboardSummary = {
  stats: [
    {
      label: 'Total products',
      value: '4',
      caption: 'Seed catalog items',
    },
    {
      label: 'Active products',
      value: '2',
      caption: 'Visible for sale',
    },
    {
      label: 'Inventory units',
      value: '1145',
      caption: 'Available stock',
    },
    {
      label: 'Draft products',
      value: '1',
      caption: 'Need review',
    },
  ],
  recentActivity: [
    'Everyday Ceramic Mug was marked active',
    'Starter Gift Kit was archived',
    'Digital Buying Guide is waiting for review',
    'Catalog settings were reviewed',
  ],
};
