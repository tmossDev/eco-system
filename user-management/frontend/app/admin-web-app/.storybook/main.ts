import type { StorybookConfig } from '@storybook/angular/dist';

const config: StorybookConfig = {
  "stories": [
    "../src/**/*.mdx",
    "../src/**/*.stories.@(js|jsx|mjs|ts|tsx)"
  ],
  "addons": [
    "@storybook/addon-a11y",
    "@storybook/addon-docs",
    "@storybook/addon-onboarding"
  ],
  "framework": "@storybook/angular",
  "webpackFinal": async (config, { configType }) => {
    if (configType === 'PRODUCTION') {
      config.output = {
        ...config.output,
        publicPath: '/stories/admin-web-app/',
      };
    }

    return config;
  }
};
export default config;
