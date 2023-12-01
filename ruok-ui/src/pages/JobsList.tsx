import { useState } from 'react';
import GenericTable from '../components/Table';
import { Box, Chip, ColorPaletteProp, Sheet, Stack, Typography } from '@mui/joy';
import { useListJobs } from '../queries/listJobs';
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
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';

type rowData = {
  key?: number;
  id: number;
  endpoint: string;
  method: string;
  expression: string;
  lastExecution: string;
  lastStatus: string;
  nextExecution: string;
  createdAt: string;
};

const StatusChip = ({ lastStatus }: { lastStatus: string }) => {
  return (
    <>
      <Chip
        variant="soft"
        size="sm"
        startDecorator={
          {
            Success: <CheckRoundedIcon />,
            Error: <BlockIcon />,
          }[lastStatus]
        }
        color={
          {
            Success: 'success',
            Error: 'danger',
          }[lastStatus as 'Error' | 'Success'] as ColorPaletteProp
        }
      >
        {lastStatus}
      </Chip>
    </>
  );
};

const Row = ({ id, endpoint, method, expression, lastExecution, lastStatus, nextExecution, createdAt }: rowData) => {
  return (
    <tr style={{ width: '100%' }}>
      <td>{id}</td>
      <td>
        <Tooltip title={endpoint} variant="outlined">
          <span>{endpoint.substring(0, 20)}...</span>
        </Tooltip>
      </td>
      <td>{method}</td>
      <td>{expression}</td>
      <td>{nextExecution}</td>
      <td>
        <StatusChip lastStatus={lastStatus} />
      </td>
      <td>{lastExecution}</td>
      <td>{createdAt}</td>
    </tr>
  );
};

const Headers = () => {
  return (
    <tr>
      <th style={{ width: 80, minWidth: 80, padding: '12px 6px' }}>Job Id</th>
      <th style={{ minWidth: 180, padding: '12px 6px' }}>Endpoint</th>
      <th style={{ width: 80, minWidth: 80, padding: '12px 6px' }}>Method</th>
      <th style={{ minWidth: 140, padding: '12px 6px' }}>CronExpression</th>
      <th style={{ minWidth: 140, padding: '12px 6px' }}>Last Execution</th>
      <th style={{ minWidth: 140, padding: '12px 6px' }}>Last Status</th>
      <th style={{ minWidth: 140, padding: '12px 6px' }}>Next Execution</th>
      <th style={{ minWidth: 140, padding: '12px 6px' }}>Created At</th>
    </tr>
  );
};

const Foot = (props: {
  totalRows: number;
  page: number;
  rowsPerPage: number;
  handleChangeRowsPerPage: (evemt: unknown, newValue: number | null) => void;
  handleChangePage: (page: number) => void;
}) => {
  return (
    <tfoot>
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
            <Typography textAlign="center" sx={{ minWidth: 80 }}>
              {props.page * props.rowsPerPage + 1} to{' '}
              {Math.min(props.page * props.rowsPerPage + props.rowsPerPage, props.totalRows)} of {props.totalRows}
            </Typography>
            <Box sx={{ display: 'flex', gap: 1 }}>
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
                disabled={props.totalRows !== -1 ? props.totalRows <= (props.page + 1) * props.rowsPerPage : false}
                onClick={() => props.handleChangePage(props.page + 1)}
                sx={{ bgcolor: 'background.surface' }}
              >
                <KeyboardArrowRightIcon />
              </IconButton>
            </Box>
          </Box>
        </td>
      </tr>
    </tfoot>
  );
};

const JobList = () => {
  const [pageSize, setPageSize] = useState(10);
  const [pageNumber, setpageNumber] = useState(0);
  const { data, error, isLoading } = useListJobs(pageSize, pageNumber * pageSize);
  return (
    <>
      <Stack spacing={4}>
        <Typography style={{ marginBottom: '1rem' }} level="h1">
          Registered Jobs
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
              Couldn't query Jobs
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
              data?.jobs.map((e, i) => {
                return (
                  <Row
                    key={i}
                    id={e.id}
                    endpoint={e.endpoint}
                    method={e.httpmethod}
                    expression={e.cronexp}
                    lastExecution={e.lastExecution}
                    lastStatus={e.lastStatus}
                    nextExecution={e.shouldExecuteAt}
                    createdAt={e.createdAt}
                  />
                );
              }) || []
            }
            Foot={
              <Foot
                page={pageNumber}
                rowsPerPage={pageSize}
                handleChangePage={(n: number) => {
                  setpageNumber(n);
                }}
                handleChangeRowsPerPage={(_event: unknown, newValue: number | null) => {
                  setPageSize(parseInt(newValue!.toString(), 10));
                }}
                totalRows={data.claimedJobs}
              />
            }
          />
        )}
      </Stack>
    </>
  );
};

export default JobList;
