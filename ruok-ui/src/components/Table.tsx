import { Sheet, Table } from '@mui/joy';

export default ({
  title,
  Headers,
  Rows,
  Foot,
}: {
  title: string;
  Headers: () => JSX.Element;
  Rows: JSX.Element[];
  Foot?: JSX.Element;
}) => {
  return (
    <>
      <Sheet
        sx={{
          height: 'auto',
          maxHeight: '78dvh',
          overflowY: 'scroll',
          '--TableCell-height': '40px',
          // the number is the amount of the header rows.
          '--TableHeader-height': 'calc(1 * var(--TableCell-height))',
          background: (theme) => `linear-gradient(${theme.vars.palette.background.surface} 30%, rgba(255, 255, 255, 0)),
            linear-gradient(rgba(255, 255, 255, 0), ${theme.vars.palette.background.surface} 70%) 0 100%,
            radial-gradient(
              farthest-side at 50% 0,
              rgba(0, 0, 0, 0.12),
              rgba(0, 0, 0, 0)
            ),
            radial-gradient(
                farthest-side at 50% 100%,
                rgba(0, 0, 0, 0.12),
                rgba(0, 0, 0, 0)
              )
              0 100%`,
          backgroundSize: '100% 40px, 100% 40px, 100% 14px, 100% 14px',
          backgroundRepeat: 'no-repeat',
          backgroundAttachment: 'local, local, scroll, scroll',
          backgroundPosition: '0 var(--TableHeader-height), 0 100%, 0 var(--TableHeader-height), 0 100%',
          backgroundColor: 'background.surface',
        }}
      >
        <Table
          aria-labelledby={title}
          stickyHeader
          hoverRow
          size="sm"
          sx={{
            '--TableCell-headBackground': 'var(--joy-palette-background-level1)',
            '--Table-headerUnderlineThickness': '1px',
            '--TableRow-hoverBackground': 'var(--joy-palette-background-level1)',
            '--TableCell-paddingY': '4px',
            '--TableCell-paddingX': '8px',
          }}
        >
          <thead style={{ width: '100%' }}>
            <Headers />
          </thead>
          <tbody style={{ width: '100%', overflowX: 'scroll' }}>{Rows}</tbody>
        </Table>
      </Sheet>
      <div
        style={{
          display: 'flex',
          justifyContent: 'flex-end',
          padding: 0,
          marginTop: '0.5rem',
          marginRight: '0.1rem',
          top: 0,
        }}
      >
        {Foot && Foot}
      </div>
    </>
  );
};
