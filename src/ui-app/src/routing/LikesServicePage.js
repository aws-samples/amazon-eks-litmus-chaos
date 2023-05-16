import React, { useState, useEffect } from 'react';
import { Grid, Card, CardContent, Typography, Divider } from '@mui/material';
import GraphComponent from "../components/GraphComponent";
import LikeComponent from "../components/LikeComponent";
import LikeService from '../services/LikeService'
import ViewsService from '../services/ViewsService'
import { useInterval } from '../utils';

function LikesServicePage() {
    const [likesError, setLikesError] = useState(false);
    const [viewsError, setViewsError] = useState(false);
    const [likes, setLikes] = useState({ likes: [] });

    const showError = () => {
        return (<h3>API Error</h3>)
    }

    useInterval(() => {
        getLikesResp();
    }, 1000 * 1);

    const getLikesResp = () => {
        // console.log("Executing getLikesResp");
        LikeService.getLikes().then((response) => {
            setLikes(response.data);
            setLikesError(false);
        })
            .catch((error) => {
                if (error.response?.status == 429) {
                    setLikesError(false);
                    setLikes(likes);
                } else {
                    setLikesError(true);
                    setLikes({ likes: [] });
                }
                console.error(`Error: ${error}`);
            })
    }

    useEffect(() => {
        getLikesResp();
    }, []);

    const postLikeResp = (data) => {
        // console.log("Executing postLikeResp", data);
        LikeService.addLike(data).then((response) => {
            setLikes(response.data);
            // setLikesError(false); // TODO: update handling of post to alert
            // console.log(response.data);
        })
            .catch((error) => {
                if (error.response?.status == 429) {
                    // setLikesError(false);
                    setLikes(likes);
                } else {
                    // setLikesError(true);
                    setLikes({ likes: [] });
                }
                console.error(`Error: ${error}`);
            })
    }

    const handle_like = (id) => {
        // console.log(id);
        postLikeResp({ id: id });
    }

    const getTotalLikes = () => {
        if (Object.keys(likes.likes).length > 0) {
            return (
                <>
                    {likes.likes.map(like => like.count).reduce((prev, next) => prev + next)} Likes
                </>
            )
        }
    }

    const [views, setViews] = useState({});

    useInterval(() => {
        getViewsResp();
    }, 1000 * 1);

    const getViewsResp = () => {
        // console.log("Executing getViewsResp");
        ViewsService.getViews().then((response) => {
            setViews(response.data);
            setViewsError(false);
        })
            .catch((error) => {
                if (error.response?.status == 429) {
                    setViewsError(false);
                    setViews(views);
                } else {
                    setViewsError(true);
                    setViews({});
                }
                console.error(`Error: ${error}`);
            })
    }

    useEffect(() => {
        postViewsResp();
    }, []);

    const postViewsResp = (data) => {
        // console.log("Executing postViewsResp", data);
        ViewsService.addView().then((response) => {
            setViews(response.data);
            // setViewsError(false);
        })
            .catch((error) => {
                if (error.response?.status == 429) {
                    // setViewsError(false);
                    setViews(views);
                } else {
                    // setViewsError(true);
                    setViews({});
                }
                console.error(`Error: ${error}`);
            })
    }

    const getLikesHostname = () => {
        if (Object.keys(likes.likes).length > 0) {
            return (
                <>
                    Pod {likes.hostname}
                </>
            )
        }
    }

    const getViewsHostname = () => {
        if (Object.keys(views).length > 0) {
            return (
                <>
                    Pod {views.hostname}
                </>
            )
        }
    }

    const getViewsCount = () => {
        if (Object.keys(views).length > 0) {
            return (
                <>
                    {views.count} Views
                </>
            )
        }
    }

    return (
        <>
            <Typography variant="h4" align="center" style={{ marginTop: '20px' }}>
                Like your most used compute service!
            </Typography>
            <Typography variant="subtitle1" align="center" style={{ marginTop: '0px' }}>
                Which service do you use to run your workloads? Hit that 'like' button!
            </Typography>
            <Grid container spacing={2} style={{ marginTop: '5px', height: '100%' }} justifyContent="center" alignItems="center">
                <Grid item xs={12} sm={12} md={5} lg={5} xl={5}>
                    <Card style={{ borderRadius: '20px', boxShadow: '0px 10px 20px rgba(0, 0, 0, 0.3)' }}>
                        <CardContent>
                            <Grid container direction="column" justifyContent="center" alignItems="center" >
                                <GraphComponent likes={likes.likes} />
                            </Grid>
                            <Typography align="center" variant="h6" color="textSecondary">
                                {likesError ? showError() : getTotalLikes()}
                            </Typography>
                            <Divider style={{ margin: '5px 0', width: '100%' }}></Divider>
                            <Typography align="center" variant="subtitle1" color="textSecondary">
                                {likesError ? showError() : getLikesHostname()}
                            </Typography>
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item xs={12} sm={10} md={6} lg={6} xl={6}>

                    <Grid container spacing={2} direction="row" >

                        {likesError ? showError() : <LikeComponent likes={likes.likes} handler={handle_like} />}

                    </Grid>

                </Grid>
                <Grid item xs={11} sm={11} md={11} lg={11} xl={11}>
                    <Card style={{ borderRadius: '20px', boxShadow: '0px 10px 20px rgba(0, 0, 0, 0.3)', height: '100%' }}>
                        <CardContent>
                            <Typography variant="h5" align="center" style={{ marginTop: '0px' }}>
                                {viewsError ? showError() : getViewsCount()}
                            </Typography>
                            <Divider style={{ margin: '5px 0', width: '100%' }}></Divider>
                            <Typography align="center" variant="subtitle1" color="textSecondary">
                                {viewsError ? showError() : getViewsHostname()}
                            </Typography>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </>
    )
}

export default LikesServicePage;