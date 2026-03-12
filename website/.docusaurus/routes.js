import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/core/',
    component: ComponentCreator('/core/', 'ef7'),
    routes: [
      {
        path: '/core/',
        component: ComponentCreator('/core/', '0b7'),
        routes: [
          {
            path: '/core/',
            component: ComponentCreator('/core/', '95d'),
            routes: [
              {
                path: '/core/examples/overview',
                component: ComponentCreator('/core/examples/overview', '387'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/core/modules/auth',
                component: ComponentCreator('/core/modules/auth', '598'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/core/modules/blob',
                component: ComponentCreator('/core/modules/blob', '73f'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/core/modules/notification',
                component: ComponentCreator('/core/modules/notification', '2f1'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/core/modules/permission',
                component: ComponentCreator('/core/modules/permission', '365'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/core/modules/storage',
                component: ComponentCreator('/core/modules/storage', 'ec2'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/core/modules/tenant',
                component: ComponentCreator('/core/modules/tenant', 'edf'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/core/',
                component: ComponentCreator('/core/', 'e91'),
                exact: true,
                sidebar: "tutorialSidebar"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];
