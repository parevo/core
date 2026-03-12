// @ts-check
// `@type` JSDoc annotations allow editor autocompletion and type checking
const config = {
  title: 'Parevo Core',
  tagline: 'Framework-agnostic Go library for auth, tenant, and permission management',
  favicon: 'img/favicon.svg',
  url: 'https://parevo.github.io',
  baseUrl: '/core/',
  organizationName: 'parevo',
  projectName: 'core',
  onBrokenLinks: 'warn',
  markdown: {
    mermaid: false,
  },
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          routeBasePath: '/',
          sidebarPath: './sidebars.js',
          editUrl: 'https://github.com/parevo/core/tree/main/website/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      },
    ],
  ],
  themeConfig: {
    navbar: {
      title: 'Parevo Core',
      logo: {
        alt: 'Parevo',
        src: 'img/logo.svg',
      },
      hideOnScroll: false,
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          href: 'https://github.com/parevo/core',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Introduction',
              to: '/',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/parevo/core',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Parevo. Built with Docusaurus.`,
    },
  },
};

export default config;
