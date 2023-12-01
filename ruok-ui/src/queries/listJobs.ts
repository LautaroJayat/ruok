import { useQuery } from 'react-query';

export const key = '[listJobsQueryKey]';

export const useListJobs = (limit: number, offset: number) => {
  const _limit = limit || 10;
  const _offset = offset || 0;

  const rows: {
    id: number;
    endpoint: string;
    method: string;
    expression: string;
    lastExecution: string;
    lastStatus: string;
    nextExecution: string;
    createdAt: string;
  }[] = [];

  for (let i = _offset; i <= _limit; i++) {
    rows.push({
      id: i,
      endpoint: 'https://google.com',
      method: 'GET',
      expression: `*/${i} ${i} * * *`,
      lastExecution: `2023-01-0${i} 12:${i}:00`,
      lastStatus: i % 2 === 0 ? 'Success' : 'Error',
      nextExecution: `2023-01-01 12:${i + 5}:00`,
      createdAt: `2023-01-01 12:${i - 5}:00`,
    });
  }
  return useQuery({
    queryKey: [key, _limit, _offset],
    queryFn: () => fetch(`http://localhost:8080/v1/jobs?limit=${_limit}&offset=${_offset}`).then((res) => res.json()),
  });
};
