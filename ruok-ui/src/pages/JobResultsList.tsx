import { useMemo, useState } from 'react';
import GenericTable from '../components/Table';
import { Box, Chip, ColorPaletteProp, Sheet, Stack, Typography } from '@mui/joy';
import Loading from '../components/Loading';
import CheckRoundedIcon from '@mui/icons-material/CheckRounded';
import BlockIcon from '@mui/icons-material/Block';
import FormControl from '@mui/joy/FormControl';
import FormLabel from '@mui/joy/FormLabel';
import IconButton from '@mui/joy/IconButton';
import Tooltip from '@mui/joy/Tooltip';
import Select from '@mui/joy/Select';
import Option from '@mui/joy/Option';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardDoubleArrowLeftIcon from '@mui/icons-material/KeyboardDoubleArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import { useListJobResults } from '../queries/listJobsExecutions';
import { useLocation } from 'react-router-dom';

type rowData = {
  key?: number;
  id: number;
  jobId: number;
  endpoint: string;
  method: string;
  lastStatusCode: number;
  succeeded: 'ok' | 'error';
  lastResponseAt: string;
};

const ChipColor: Record<string, ColorPaletteProp> = {
  ok: 'success',
  error: 'danger',
};

const ChipIcon = (succeed: 'ok' | 'error') => {
  if (succeed === 'ok') {
    return <CheckRoundedIcon color="success" />;
  } else if (succeed == 'error') {
    return <BlockIcon color="error" />;
  }
};

const StatusChip = ({ lastStatus }: { lastStatus: 'ok' | 'error' }) => {
  return (
    <>
      <Chip variant="soft" size="lg" startDecorator={ChipIcon(lastStatus)}>
        <Typography color={ChipColor[lastStatus]}>{lastStatus}</Typography>
      </Chip>
    </>
  );
};

const Row = ({ id, jobId, endpoint, method, succeeded, lastStatusCode, lastResponseAt }: rowData) => {
  return (
    <tr style={{ width: '100%' }}>
      <td>{id}</td>
      <td>{jobId}</td>

      <td>
        <Tooltip title={endpoint} variant="outlined">
          <span>{endpoint.substring(0, 20)}...</span>
        </Tooltip>
      </td>
      <td>{method}</td>
      <td>
        <StatusChip lastStatus={succeeded} />
      </td>
      <td>{lastStatusCode}</td>
      <td>{new Date(lastResponseAt).toLocaleString()}</td>
    </tr>
  );
};

const Headers = () => {
  return (
    <tr>
      <th style={{ width: 80, minWidth: 80, padding: '12px 6px' }}>Id</th>
      <th style={{ width: 50, minWidth: 50, padding: '12px 6px' }}>Job Id</th>
      <th style={{ minWidth: 180, padding: '12px 6px' }}>Endpoint</th>
      <th style={{ width: 80, minWidth: 80, padding: '12px 6px' }}>Method</th>
      <th style={{ minWidth: 140, padding: '12px 6px' }}>Succeeded</th>
      <th style={{ minWidth: 140, padding: '12px 6px' }}>Last Status Code</th>
      <th style={{ minWidth: 140, padding: '12px 6px' }}>Response At</th>
    </tr>
  );
};

const Foot = (props: {
  currentPageLength: number;
  page: number;
  rowsPerPage: number;
  handleChangeRowsPerPage: (event: unknown, newValue: number | null) => void;
  handleChangePage: (page: number) => void;
}) => {
  return (
    <thead>
      <tr>
        <td colSpan={8}>
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 2,
              justifyContent: 'flex-end',
            }}
          >
            <FormControl orientation="horizontal" size="sm">
              <FormLabel>Rows per page:</FormLabel>
              <Select onChange={props.handleChangeRowsPerPage} value={props.rowsPerPage}>
                <Option value={5}>5</Option>
                <Option value={10}>10</Option>
                <Option value={25}>25</Option>
                <Option value={50}>50</Option>
                <Option value={100}>100</Option>
              </Select>
            </FormControl>
            <Box sx={{ display: 'flex', gap: 1 }}>
              <IconButton
                size="sm"
                color="neutral"
                variant="outlined"
                disabled={props.page === 0}
                onClick={() => props.handleChangePage(0)}
                sx={{ bgcolor: 'background.surface' }}
              >
                <KeyboardDoubleArrowLeftIcon />
              </IconButton>
              <IconButton
                size="sm"
                color="neutral"
                variant="outlined"
                disabled={props.page === 0}
                onClick={() => props.handleChangePage(props.page - 1)}
                sx={{ bgcolor: 'background.surface' }}
              >
                <KeyboardArrowLeftIcon />
              </IconButton>
              <IconButton
                size="sm"
                color="neutral"
                variant="outlined"
                disabled={props.currentPageLength < props.rowsPerPage}
                onClick={() => props.handleChangePage(props.page + 1)}
                sx={{ bgcolor: 'background.surface' }}
              >
                <KeyboardArrowRightIcon />
              </IconButton>
            </Box>
          </Box>
        </td>
      </tr>
    </thead>
  );
};

const JobResultsList = () => {
  const location = useLocation();

  const id = useMemo(() => {
    const splitPath = location.pathname.split('/');
    return splitPath[splitPath.length - 1];
  }, [location.pathname]);

  const [pageSize, setPageSize] = useState(10);
  const [pageNumber, setPageNumber] = useState(0);
  const { data, error, isLoading } = useListJobResults(id, pageSize, pageNumber * pageSize);
  return (
    <>
      <Stack spacing={4}>
        <Typography style={{ marginBottom: '1rem' }} level="h1">
          Executions
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
              Couldn't query Job Results
            </Typography>
          </Sheet>
        )}
        {isLoading && <Loading />}
        {data && (
          <GenericTable
            title="Registered Jobs"
            Headers={Headers}
            Rows={
              //@ts-ignore
              data?.jobResults.map((e, i) => {
                return (
                  <Row
                    key={i}
                    id={e.id}
                    jobId={e.jobId}
                    endpoint={e.endpoint}
                    method={e.httpmethod}
                    succeeded={e.succeeded}
                    lastStatusCode={e.lastStatusCode}
                    lastResponseAt={e.lastResponseAt}
                  />
                );
              }) || []
            }
            Foot={
              <Foot
                page={pageNumber}
                rowsPerPage={pageSize}
                handleChangePage={(n: number) => {
                  setPageNumber(n);
                }}
                handleChangeRowsPerPage={(_event: unknown, newValue: number | null) => {
                  setPageSize(parseInt(newValue!.toString(), 10));
                }}
                currentPageLength={data?.jobResults || 0}
              />
            }
          />
        )}
      </Stack>
    </>
  );
};

export default JobResultsList;
