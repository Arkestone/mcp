import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Arkestone MCP Servers',
  tagline: 'MCP servers for GitHub Copilot and AI coding assistants',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  url: 'https://arkestone.github.io',
  baseUrl: '/mcp/',

  organizationName: 'Arkestone',
  projectName: 'mcp',
  trailingSlash: false,

  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/Arkestone/mcp/tree/main/website/',
          routeBasePath: 'docs',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Arkestone MCP',
      logo: {
        alt: 'Arkestone MCP Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          href: 'https://github.com/Arkestone/mcp',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Servers',
          items: [
            {label: 'mcp-instructions', to: '/docs/servers/mcp-instructions'},
            {label: 'mcp-skills', to: '/docs/servers/mcp-skills'},
            {label: 'mcp-prompts', to: '/docs/servers/mcp-prompts'},
            {label: 'mcp-adr', to: '/docs/servers/mcp-adr'},
            {label: 'mcp-memory', to: '/docs/servers/mcp-memory'},
            {label: 'mcp-graph', to: '/docs/servers/mcp-graph'},
          ],
        },
        {
          title: 'Resources',
          items: [
            {label: 'Contributing', to: '/docs/contributing'},
            {label: 'Security', to: '/docs/security'},
            {label: 'Changelog', to: '/docs/changelog'},
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/Arkestone/mcp',
            },
            {
              label: 'Issues',
              href: 'https://github.com/Arkestone/mcp/issues',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Arkestone. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'yaml', 'json', 'go', 'docker'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
