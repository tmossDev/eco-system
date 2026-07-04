import type { StorybookConfig } from '@storybook/angular';

const config: StorybookConfig = {
  "stories": [
    "../../../package/design-system/primitive/*/src/**/*.stories.@(js|jsx|mjs|ts|tsx)",
    "../../../package/design-system/shared/*/src/**/*.stories.@(js|jsx|mjs|ts|tsx)"
  ],
  "addons": [
    "@storybook/addon-a11y",
    "@storybook/addon-docs",
    "@storybook/addon-onboarding"
  ],
  "framework": "@storybook/angular",
  "refs": {
    "admin-web-app": {
      "title": "Admin Web App",
      "url": "/stories/admin-web-app"
    }
  }
};
export default config;
