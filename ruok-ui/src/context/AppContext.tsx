import React from 'react';

export const JobsResultNamespace = '[jobsResultNamespace]';

export const initialContext = {
  [JobsResultNamespace]: {
    jobId: '',
    jobName: '',
  },
};

export const AppContext = React.createContext(initialContext);
