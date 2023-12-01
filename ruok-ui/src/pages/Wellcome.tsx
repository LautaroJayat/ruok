import { Container, Box, typographyClasses, AspectRatio, Button, Card, Typography } from '@mui/joy';
import FavoriteRoundedIcon from '@mui/icons-material/FavoriteRounded';
import { useNavigate } from 'react-router-dom';

export const WellcomeMsg = () => {
  const navigate = useNavigate();
  return (
    <>
      <Typography level="h1" fontWeight="xl" fontSize="clamp(1.875rem, 1.3636rem + 2.1818vw, 3rem)">
        Wellcome to "ruOK?"
      </Typography>
      <Typography fontSize="lg" textColor="text.secondary" lineHeight="lg">
        This dashboard aims to be a comprehensive tool for monitoring health metrics, registered jobs, and pertinent
        information associated with the current instance.
      </Typography>
      <Card variant="outlined" color="neutral" orientation="horizontal" sx={{ gap: 2, my: 1, textAlign: 'left' }}>
        {/** @ts-ignore */}
        <FavoriteRoundedIcon style={{ color: '#F75F7B' }} color="success" fontSize="xl3" />
        <div>
          <Typography fontSize="xl" fontWeight="lg" sx={{ mb: 1 }}>
            We appreciate your trust in our tool!
          </Typography>
          <Typography level="body-sm">
            We are committed to providing a seamless experience.
            <br></br>
            This app is currently under development, so please stay tuned for updates
          </Typography>
        </div>
      </Card>
      <Button
        onClick={(e) => {
          e.preventDefault();
          navigate('/instance');
        }}
        variant="outlined"
        size="lg"
      >
        Start by checking this instance info!
      </Button>
    </>
  );
};

export default () => {
  return (
    <Container
      sx={(theme) => ({
        position: 'relative',
        minHeight: '100%',
        display: 'flex',
        alignItems: 'center',
        py: 10,
        gap: 4,
        [theme.breakpoints.up(834)]: {
          flexDirection: 'row',
          gap: 6,
        },
        [theme.breakpoints.up(1199)]: {
          gap: 12,
        },
      })}
    >
      <Box
        sx={(theme) => ({
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: '1rem',
          maxWidth: '50ch',
          textAlign: 'center',
          flexShrink: 999,
          [theme.breakpoints.up(834)]: {
            minWidth: 420,
            alignItems: 'flex-start',
            textAlign: 'initial',
          },
          [`& .${typographyClasses.root}`]: {
            textWrap: 'balance',
          },
        })}
      >
        <WellcomeMsg />
      </Box>
      <AspectRatio
        ratio={3 / 4}
        objectFit="contain"
        variant="plain"
        maxHeight={600}
        sx={(theme) => ({
          minWidth: 300,
          alignSelf: 'stretch',
          [theme.breakpoints.up(834)]: {
            alignSelf: 'initial',
            flexGrow: 1,
            '--AspectRatio-maxHeight': '520px',
            '--AspectRatio-minHeight': '400px',
          },
          flexBasis: '50%',
        })}
      >
        <img src="/big_ruok_logo.svg" alt="ruok_logo" />
      </AspectRatio>
    </Container>
  );
};
