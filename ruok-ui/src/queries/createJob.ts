export const key = '[createJob]';

export const createJob = (body: string) =>
  fetch('http://localhost:8080/v1/jobs', { method: 'POST', body }).then((res) => res.json());
