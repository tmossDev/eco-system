export const iconNames = [
  'check',
  'close',
  'info',
  'warning',
  'arrow-right',
  'arrow-left',
  'arrow-up',
  'arrow-down',
  'chevron-right',
  'chevron-left',
  'chevron-up',
  'chevron-down',
  'plus',
  'minus',
  'search',
  'settings',
  'user',
  'home',
  'trash',
  'edit',
  'calendar',
  'menu',
  'heart',
  'star',
  'download',
  'upload',
] as const;

export type IconName = (typeof iconNames)[number];

export interface IconDefinition {
  path: string;
  viewBox?: string;
}

export const iconRegistry: Record<IconName, IconDefinition> = {
  check: {
    path: 'M9 16.17 4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z',
  },
  close: {
    path: 'M18.3 5.71 12 12l6.3 6.29-1.41 1.41L10.59 13.41 4.29 19.71 2.88 18.3 9.17 12 2.88 5.71 4.29 4.29l6.3 6.3 6.29-6.3z',
  },
  info: {
    path: 'M11 17h2v-6h-2zm1-14a9 9 0 1 0 0 18 9 9 0 0 0 0-18zm0 16a7 7 0 1 1 0-14 7 7 0 0 1 0 14zm-1-10h2V7h-2z',
  },
  warning: {
    path: 'M1 21h22L12 2zm12-3h-2v-2h2zm0-4h-2v-4h2z',
  },
  'arrow-right': {
    path: 'M12 4l1.41 1.41L8.83 10H20v2H8.83l4.58 4.59L12 18l-8-7z',
  },
  'arrow-left': {
    path: 'M12 20l-1.41-1.41L15.17 14H4v-2h11.17l-4.58-4.59L12 6l8 7z',
  },
  'arrow-up': {
    path: 'M4 12l1.41 1.41L10 8.83V20h2V8.83l4.59 4.58L18 12l-7-8z',
  },
  'arrow-down': {
    path: 'M20 12l-1.41-1.41L14 15.17V4h-2v11.17l-4.59-4.58L6 12l7 8z',
  },
  'chevron-right': {
    path: 'M9.29 6.71a1 1 0 0 1 1.42 0L16 12l-5.29 5.29a1 1 0 1 1-1.42-1.42L13.17 12 9.29 8.12a1 1 0 0 1 0-1.41z',
  },
  'chevron-left': {
    path: 'M14.71 6.71a1 1 0 0 0-1.42 0L8 12l5.29 5.29a1 1 0 1 0 1.42-1.42L10.83 12l3.88-3.88a1 1 0 0 0 0-1.41z',
  },
  'chevron-up': {
    path: 'M6.71 14.71a1 1 0 0 1 0-1.42L12 8l5.29 5.29a1 1 0 1 1-1.42 1.42L12 10.83l-3.88 3.88a1 1 0 0 1-1.41 0z',
  },
  'chevron-down': {
    path: 'M6.71 9.29a1 1 0 0 1 1.42 0L12 13.17l3.88-3.88a1 1 0 1 1 1.42 1.42L12 16 6.71 10.71a1 1 0 0 1 0-1.42z',
  },
  plus: {
    path: 'M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6z',
  },
  minus: {
    path: 'M5 11h14v2H5z',
  },
  search: {
    path: 'M9.5 3a6.5 6.5 0 0 1 5.17 10.44l4.45 4.44-1.42 1.42-4.44-4.45A6.5 6.5 0 1 1 9.5 3zm0 2a4.5 4.5 0 1 0 0 9 4.5 4.5 0 0 0 0-9z',
  },
  settings: {
    path: 'M19.43 12.98c.04-.32.07-.65.07-.98s-.02-.66-.07-.98l2.11-1.65-2-3.46-2.49 1a7.28 7.28 0 0 0-1.69-.98L15 3h-4l-.36 2.93c-.61.24-1.18.56-1.69.98l-2.49-1-2 3.46 2.11 1.65c-.04.32-.07.65-.07.98s.02.66.07.98l-2.11 1.65 2 3.46 2.49-1c.51.4 1.08.73 1.69.98L11 21h4l.36-2.93c.61-.24 1.18-.56 1.69-.98l2.49 1 2-3.46zm-6.43 2.52A3.5 3.5 0 1 1 13 8a3.5 3.5 0 0 1 0 7z',
  },
  user: {
    path: 'M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8zm0 2c-4.42 0-8 2.24-8 5v1h16v-1c0-2.76-3.58-5-8-5z',
  },
  home: {
    path: 'M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z',
  },
  trash: {
    path: 'M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6zm3-9h2v8H9zm4 0h2v8h-2zM15.5 4l-1-1h-5l-1 1H5v2h14V4z',
  },
  edit: {
    path: 'M3 17.25V21h3.75L17.81 9.94l-3.75-3.75zM20.71 7.04a1 1 0 0 0 0-1.41l-2.34-2.34a1 1 0 0 0-1.41 0l-1.83 1.83 3.75 3.75z',
  },
  calendar: {
    path: 'M19 3h-1V1h-2v2H8V1H6v2H5c-1.11 0-2 .9-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5c0-1.1-.9-2-2-2zm0 16H5V8h14z',
  },
  menu: {
    path: 'M3 6h18v2H3zm0 5h18v2H3zm0 5h18v2H3z',
  },
  heart: {
    path: 'M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5A5.45 5.45 0 0 1 7.5 3c1.74 0 3.41.81 4.5 2.09A6.02 6.02 0 0 1 16.5 3 5.45 5.45 0 0 1 22 8.5c0 3.78-3.4 6.86-8.55 11.54z',
  },
  star: {
    path: 'M12 17.27 18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z',
  },
  download: {
    path: 'M5 20h14v-2H5zm14-9h-4V3H9v8H5l7 7z',
  },
  upload: {
    path: 'M5 20h14v-2H5zM9 16h6V8h4l-7-7-7 7h4z',
  },
};
