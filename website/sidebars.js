const sidebars = {
  tutorialSidebar: [
    'index',
    {
      type: 'category',
      label: 'Modules',
      items: [
        'modules/auth',
        'modules/tenant',
        'modules/permission',
        'modules/storage',
        'modules/query',
        'modules/notification',
        'modules/blob',
        'modules/cache',
        'modules/health',
        'modules/lock',
        'modules/billing',
        'modules/job',
        'modules/search',
        'modules/export',
        'modules/validation',
        'modules/geo',
      ],
    },
    {
      type: 'category',
      label: 'Examples',
      items: ['examples/overview'],
    },
  ],
};

export default sidebars;
