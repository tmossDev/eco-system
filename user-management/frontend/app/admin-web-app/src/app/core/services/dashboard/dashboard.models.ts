export interface DashboardStat {
  label: string;
  value: string;
  caption: string;
}

export interface DashboardSummary {
  stats: DashboardStat[];
  recentActivity: string[];
}
