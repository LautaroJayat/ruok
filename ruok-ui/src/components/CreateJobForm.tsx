import * as React from 'react';
import { useState } from 'react';
import Button from '@mui/joy/Button';
import { useMutation } from 'react-query';
import FormControl from '@mui/joy/FormControl';
import FormLabel from '@mui/joy/FormLabel';
import Input from '@mui/joy/Input';
import ModalDialog from '@mui/joy/ModalDialog';
import DialogTitle from '@mui/joy/DialogTitle';
import DialogContent from '@mui/joy/DialogContent';
import Stack from '@mui/joy/Stack';
import Select from '@mui/joy/Select';
import Option from '@mui/joy/Option';
import HelpOutlineSharpIcon from '@mui/icons-material/HelpOutlineSharp';
import { ModalClose, Tooltip } from '@mui/joy';
import DebounceInput from './DebounceInput';
import { createJob } from '../queries/createJob';

const successStatusesHelp = `
Use a comma separated list of success statuses that will signal an OK backend.
`;

const parseSuccessStatuses = (input: string): [number[], boolean] => {
  try {
    const list = input.split(',').map((e) => parseInt(e));
    const hasError = list.some((e) => isNaN(e));
    if (hasError) {
      return [[], true];
    }
    return [list, false];
  } catch {
    return [[], true];
  }
};

export const CreateJobForm = React.forwardRef(() => {
  const mutation = useMutation(createJob);
  const [cron, setCron] = useState('');
  const [endpoint, setEndpoint] = useState('');
  const [method, setMethod] = useState('');
  const [success, setSuccess] = useState([] as number[]);
  const [successError, setSuccessError] = useState(false);
  if (mutation.isSuccess) {
    return (
      <ModalDialog>
        <ModalClose />
        <DialogTitle>New monitor job created!</DialogTitle>
      </ModalDialog>
    );
  }

  return (
    <ModalDialog>
      <ModalClose />
      <DialogTitle>Create a new monitor job!</DialogTitle>
      <DialogContent>Fill in the fields to create a new monitoring job.</DialogContent>
      <form
        onSubmit={(event: React.FormEvent<HTMLFormElement>) => {
          event.preventDefault();
        }}
      >
        <Stack spacing={2}>
          <FormControl>
            <FormLabel>Cron expression</FormLabel>
            <Input onChange={(event: React.ChangeEvent<HTMLInputElement>) => setCron(event.target.value)} required />
          </FormControl>
          <FormControl>
            <FormLabel>Endpoint</FormLabel>
            <Input
              onChange={(event: React.ChangeEvent<HTMLInputElement>) => setEndpoint(event.target.value)}
              required
            />
          </FormControl>
          <FormControl>
            <FormLabel>HTTP Method</FormLabel>
            <Select
              placeholder="Select a Method"
              onChange={(_: React.SyntheticEvent | null, newValue: string | null) => {
                if (newValue) setMethod(newValue);
              }}
              name="Method"
              required
              sx={{ minWidth: 200 }}
            >
              <Option value="GET">GET</Option>
              <Option value="POST">POST</Option>
              <Option value="PUT">PUT</Option>
              <Option value="DELETE">DELETE</Option>
            </Select>
          </FormControl>
          <FormControl>
            <FormLabel>Success Statuses</FormLabel>
            <Stack spacing={1} direction="row" alignItems="center" justifyContent="space-between">
              <DebounceInput
                error={successError}
                fullWidth
                required
                handleDebounce={(v) => {
                  const [list, err] = parseSuccessStatuses(v);
                  console.log(list);
                  if (!err) {
                    setSuccess(list);
                  }
                  setSuccessError(err);
                }}
                debounceTimeout={200}
              />
              <Tooltip title={successStatusesHelp}>
                <HelpOutlineSharpIcon />
              </Tooltip>
            </Stack>
          </FormControl>
          <Button
            disabled={mutation.isLoading}
            onClick={(e) => {
              e.preventDefault();
              if (!!cron && !!endpoint && !!method && success.length > 0) {
                const body = JSON.stringify({
                  cronexp: cron,
                  maxRetries: 1,
                  endpoint,
                  httpmethod: method,
                  successStatuses: success,
                });
                mutation.mutate(body);
              } else {
                console.log(cron, endpoint, method, success.length);
              }
            }}
            type="submit"
          >
            {!mutation.isLoading ? 'Create!' : 'Creating...'}
          </Button>
        </Stack>
      </form>
    </ModalDialog>
  );
});
