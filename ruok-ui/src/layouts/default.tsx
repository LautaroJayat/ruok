import { Outlet } from 'react-router-dom';
import { Box } from '@mui/joy';

import Header from '../components/Header';
import Sidebar from '../components/Sidebar';

export default () => {
  return (
    <Box sx={{ display: 'flex', minHeight: '100dvh', maxHeight: '100dvh', overflow: 'hidden' }}>
      <Header />
      <Sidebar />
      <Box
        component="main"
        className="MainContent"
        sx={{
          px: { xs: 2, md: 6 },
          pt: {
            xs: 'calc(12px + var(--Header-height))',
            sm: 'calc(12px + var(--Header-height))',
            md: 3,
          },
          pb: { xs: 2, sm: 2, md: 3 },
          flex: 1,
          display: 'flex',
          flexDirection: 'column',
          minWidth: 0,
          height: '100dvh',
          minHeight: '100dvh',
          maxHeight: '100dvh',
          gap: 1,
          overflow: 'hidden',
        }}
      >
        <Outlet />
      </Box>
    </Box>
  );
};
