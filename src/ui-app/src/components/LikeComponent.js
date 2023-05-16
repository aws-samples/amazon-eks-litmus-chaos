import { Grid, Card, CardActions, Button, CardContent, Typography } from '@mui/material';
import FavoriteIcon from '@mui/icons-material/Favorite';

export default function LikeComponent(props) {

  // TODO: Optimize
  const displayLikes = () => {
    if (Object.keys(props.likes).length > 0) {
      return (
        <>
              {props.likes.map((like, index) => (
                <Grid item xs={12} sm={6} md={6} lg={6} xl={6} >
                  <Card style={{ borderRadius: '20px', boxShadow: '0px 10px 20px rgba(0, 0, 0, 0.3)', height: '245px' }} key={like.id}>
                    <CardContent style={{ marginTop: '20px' }}>

                      <Grid container spacing={2} justifyContent="space-evenly"
                        alignItems="stretch">
                        <Grid item xs={3} sm={4}>
                          <img src={like.image} alt="Placeholder" style={{ maxWidth: '100%', maxHeight: '100%' }} />
                        </Grid>
                        <Grid item xs={3} sm={4}>
                          <Typography variant="subtitle1" gutterBottom>
                          {like.name}
                          </Typography>
                          <Typography variant="h4" style={{ display: 'inline-block' }} gutterBottom>
                          {like.count}
                          </Typography>
                          <Typography variant="caption" style={{ display: 'inline-block' }} gutterBottom>
                            &nbsp;likes
                          </Typography>
                        </Grid>
                      </Grid>
                    </CardContent>
                    <CardActions style={{ justifyContent: 'center' }}>
                      <Button style={{ width: '80%' }} color="squid" variant="contained" startIcon={<FavoriteIcon />} onClick={() => {
                props.handler(like.id);
              } }>
                        Like
                      </Button>
                    </CardActions>
                  </Card>
                </Grid>

              ))}
          </>
      )
    } 
  }

  return (
    <>
      {displayLikes()}
    </>
  );
};