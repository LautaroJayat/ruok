import { lazy } from 'react';

const JobsList = lazy(() => import('../pages/JobsList'));
const InstanceInfo = lazy(() => import('../pages/IntanceInfo'));

export default [
  {
    path: '/jobs',
    title: 'Registered Jobs',
    component: JobsList,
  },
  /*{
    path: '/jobs/:jobId',
    title: 'Job Executions',
  },*/
  {
    path: '/instance',
    title: 'Instance Info',
    component: InstanceInfo,
  },
];
