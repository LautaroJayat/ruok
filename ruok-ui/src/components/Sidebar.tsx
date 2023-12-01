import { useCallback } from 'react';
import GlobalStyles from '@mui/joy/GlobalStyles';
import Box from '@mui/joy/Box';
import IconButton from '@mui/joy/IconButton';
import List from '@mui/joy/List';
import ListItem from '@mui/joy/ListItem';
import ListItemButton, { listItemButtonClasses } from '@mui/joy/ListItemButton';
import ListItemContent from '@mui/joy/ListItemContent';
import Typography from '@mui/joy/Typography';
import Sheet from '@mui/joy/Sheet';
import SettingsRoundedIcon from '@mui/icons-material/SettingsRounded';
import FavoriteRoundedIcon from '@mui/icons-material/FavoriteRounded';
import AccessTimeFilledRoundedIcon from '@mui/icons-material/AccessTimeFilledRounded';
import HelpRoundedIcon from '@mui/icons-material/HelpRounded';
import { useNavigate } from 'react-router-dom';

//import ColorSchemeToggle from './ColorSchemeToggle';
//import { closeSidebar } from '../utils';

type SVGIcon = typeof FavoriteRoundedIcon;

function SidebarButton({
  Icon,
  title,
  path,
  outside = false,
}: {
  Icon: SVGIcon;
  path: string;
  title: string;
  outside: boolean;
}) {
  const navigate = useNavigate();
  const navigateInsideApp = useCallback(() => navigate(path, { state: { key: path } }), []);
  const navigateOutsideApp = useCallback(() => window.open(path, '_blank'), []);

  return (
    <ListItem>
      <ListItemButton onClick={outside ? navigateOutsideApp : navigateInsideApp}>
        <Icon />
        <ListItemContent>
          <Typography fontSize={13} fontWeight={600}>
            {title}
          </Typography>
        </ListItemContent>
      </ListItemButton>
    </ListItem>
  );
}

export default function Sidebar() {
  const navigate = useNavigate();

  return (
    <Sheet
      className="Sidebar"
      sx={{
        position: { xs: 'fixed', md: 'sticky' },
        transform: {
          xs: 'translateX(calc(100% * (var(--SideNavigation-slideIn, 0) - 1)))',
          md: 'none',
        },
        transition: 'transform 0.4s, width 0.4s',
        zIndex: 10000,
        height: '100dvh',
        width: 'var(--Sidebar-width)',
        top: 0,
        p: 2,
        flexShrink: 0,
        display: 'flex',
        flexDirection: 'column',
        gap: 2,
        borderRight: '1px solid',
        borderColor: 'divider',
      }}
    >
      <GlobalStyles
        styles={(theme) => ({
          ':root': {
            '--Sidebar-width': '200px',
            [theme.breakpoints.up('lg')]: {
              '--Sidebar-width': '220px',
            },
          },
        })}
      />
      <Box
        className="Sidebar-overlay"
        sx={{
          position: 'fixed',
          zIndex: 9998,
          top: 0,
          left: 0,
          width: '100vw',
          height: '100vh',
          opacity: 'var(--SideNavigation-slideIn)',
          backgroundColor: 'var(--joy-palette-background-backdrop)',
          transition: 'opacity 0.4s',
          transform: {
            xs: 'translateX(calc(100% * (var(--SideNavigation-slideIn, 0) - 1) + var(--SideNavigation-slideIn, 0) * var(--Sidebar-width, 0px)))',
            lg: 'translateX(-100%)',
          },
        }}
        onClick={() => {}}
      />
      <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
        <IconButton onClick={() => navigate('/')} variant="soft" color="primary" style={{ color: '#F75F7B' }} size="sm">
          <FavoriteRoundedIcon />
        </IconButton>
        <Typography fontFamily="'Exo', sans-serif;" fontWeight={800} level="title-lg">
          ruOK?
        </Typography>
        {/*<ColorSchemeToggle sx={{ ml: 'auto' }} />*/}
      </Box>
      <Box
        sx={{
          minHeight: 0,
          overflow: 'hidden auto',
          flexGrow: 1,
          display: 'flex',
          flexDirection: 'column',
          [`& .${listItemButtonClasses.root}`]: {
            gap: 1.5,
          },
        }}
      >
        <List
          size="sm"
          sx={{
            gap: 1,
            '--List-nestedInsetStart': '30px',
            '--ListItem-radius': (theme) => theme.vars.radius.sm,
          }}
        >
          <SidebarButton outside={false} path="/instance" Icon={SettingsRoundedIcon} title="Instance Info" />
          {/*<SidebarButton outside={false} path="/jobs" Icon={AccessTimeFilledRoundedIcon} title="Jobs By User" />*/}
          <SidebarButton outside={false} path="/jobs" Icon={AccessTimeFilledRoundedIcon} title="Claimed Jobs" />
          <SidebarButton
            outside={true}
            path="https://github.com/back-end-labs/ruok-scheduler"
            Icon={HelpRoundedIcon}
            title="Help"
          />
        </List>
      </Box>
    </Sheet>
  );
}
