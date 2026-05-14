import { DashboardSummary } from './dashboard.models';

export const MOCK_DASHBOARD_SUMMARY: DashboardSummary = {
  stats: [
    {
      label: 'Total users',
      value: '128',
      caption: '12 added this month',
    },
    {
      label: 'Active users',
      value: '96',
      caption: '75% active rate',
    },
    {
      label: 'Pending invites',
      value: '8',
      caption: 'Awaiting response',
    },
    {
      label: 'Admin users',
      value: '5',
      caption: 'Privileged accounts',
    },
  ],
  recentActivity: [
    'New user account created for Alex Morgan',
    'Priya Shah updated her profile details',
    'Admin permissions changed for Jordan Lee',
    'System settings were reviewed',
  ],
};
