import Card from '@mui/joy/Card';
import CardContent from '@mui/joy/CardContent';
import Typography from '@mui/joy/Typography';

export default ({
  variant = 'outlined',
  title,
  description,
}: {
  variant: 'plain' | 'outlined';
  title: string;
  description: string;
}) => {
  return (
    <Card variant={variant}>
      <CardContent>
        <Typography level="title-lg">{title}</Typography>
        <Typography level="body-lg">{description}</Typography>
      </CardContent>
    </Card>
  );
};
