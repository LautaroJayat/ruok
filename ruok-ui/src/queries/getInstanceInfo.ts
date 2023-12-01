import { useQuery } from 'react-query';

export const key = '[getInstanceInfo]';

export const MockedQuery = async () => {
  return await Promise.resolve({
    appName: 'Application 1',
    dbConnected: true,
    dbURL: 'postgresql://****:****@localhost:5432/db1',
    tlsActive: true,
    tlsVersion: 'TLS v1.3',
    startedAt: 123456,
    upTimeMicro: 123456,
    claimedJobs: 123,
  });
};

export default () =>
  useQuery({
    queryKey: [key],
    cacheTime: 0,
    queryFn: () => fetch(`http://localhost:8080/v1/instance`).then((res) => res.json()),
  });
