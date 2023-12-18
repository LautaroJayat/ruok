import { Stack, Typography } from '@mui/joy';
import Sheet from '@mui/joy/Sheet';
import Simplecard from '../components/SimpleCard';
import getInstanceInfo from '../queries/getInstanceInfo';

import Loading from '../components/Loading';
import { useState, useEffect } from 'react';

const second = 1000000;
const minute = 60 * second;
const hour = 60 * minute;
const Uptime = ({ micro }: { micro: number }) => {
  const [upTimeString, setUptimeString] = useState('');

  useEffect(() => {
    let up = micro;
    let h = Math.floor(micro / hour);
    let m = Math.floor((micro - hour * h) / minute);
    let s = Math.floor((micro - (hour * h + m * minute)) / second);
    let str = `${h} h - ${m} m - ${s} s`;
    setUptimeString(str);
    const interval = setInterval(() => {
      up += second;
      h = Math.floor(up / hour);
      m = Math.floor((up - hour * h) / minute);
      s = Math.floor((up - (hour * h + m * minute)) / second);
      str = `${h} h - ${m} m - ${s} s`;
      setUptimeString(str);
      console.log('tick', str);
    }, 1000);
    return () => clearInterval(interval);
  }, []);
  return <Simplecard variant="outlined" title="Uptime" description={upTimeString} />;
};

const InstanceInfo = () => {
  const { data, error, isLoading } = getInstanceInfo();

  if (error) return <div>Request Failed</div>;
  return (
    <>
      <Stack spacing={4}>
        <Typography style={{ marginBottom: '1rem' }} level="h1">
          Instance Info
        </Typography>
        {!!error && (
          <Sheet
            style={{
              minWidth: '20rem',
              minHeight: '20rem',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
            color="danger"
            variant="soft"
          >
            <Typography color="danger" level="body-lg">
              Couldn't query for instance info
            </Typography>
          </Sheet>
        )}

        {isLoading && <Loading />}
        {data && (
          <Stack direction="row" spacing={2} justifyContent="center" alignItems="stretch">
            <Stack direction="column" spacing={2}>
              <Simplecard variant="outlined" title="Application Name" description={data?.appName || 'unknown'} />
              <Simplecard variant="outlined" title="Database Connected" description={`${data?.dbConnected}`} />
              <Simplecard variant="outlined" title="Database URL" description={data?.dbURL || 'unknown'} />
              <Simplecard variant="outlined" title="TLS Active" description={`${data?.tlsActive}`} />
              <Simplecard variant="outlined" title="TLS Version" description={data?.tlsVersion || 'unknown'} />
            </Stack>

            <Stack direction="column" spacing={2}>
              <Simplecard
                variant="outlined"
                title="Started At"
                description={`${new Date(data?.startedAtMicro / 1000).toLocaleString()}` || 'unknown'}
              />
              <Uptime micro={data?.upTimeMicro || 0} />
              <Simplecard variant="outlined" title="Claimed Jobs" description={`${data?.claimedJobs}` || 'unknown'} />
              <Simplecard variant="outlined" title="Max Jobs" description={`${data?.maxJobs}` || 'unknown'} />
            </Stack>
          </Stack>
        )}
      </Stack>
    </>
  );
};

export default InstanceInfo;
