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
        'modules/notification',
        'modules/blob',
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
