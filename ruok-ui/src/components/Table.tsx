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
    <Sheet
      className="OrderTableContainer"
      variant="outlined"
      sx={{
        display: { xs: 'none', sm: 'initial' },
        width: '100%',
        borderRadius: 'sm',
        flexShrink: 1,
        overflow: 'auto',
        minHeight: 0,
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
        <thead>
          <Headers />
        </thead>
        <tbody>{Rows}</tbody>
        {Foot && Foot}
      </Table>
    </Sheet>
  );
};
