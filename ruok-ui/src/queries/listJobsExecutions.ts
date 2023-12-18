import { useQuery } from 'react-query';

export const key = '[listJobResultsQueryKey]';

export const useListJobResults = (id: string | number, limit: number, offset: number) => {
  const _limit = limit || 10;
  const _offset = offset || 0;

  return useQuery({
    queryKey: [key, _limit, _offset],
    queryFn: () =>
      fetch(`http://localhost:8080/v1/jobs/${id}?limit=${_limit}&offset=${_offset}`).then((res) => res.json()),
  });
};
