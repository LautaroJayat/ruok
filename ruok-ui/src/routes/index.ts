import { lazy } from 'react';

const JobsList = lazy(() => import('../pages/JobsList'));
const InstanceInfo = lazy(() => import('../pages/InstanceInfo'));
const JobResultsList = lazy(() => import('../pages/JobResultsList'));

export default [
  {
    path: '/jobs',
    title: 'Registered Jobs',
    component: JobsList,
  },
  {
    path: '/jobs/:jobId',
    title: 'Job Executions',
    component: JobResultsList,
  },
  {
    path: '/instance',
    title: 'Instance Info',
    component: InstanceInfo,
  },
];
