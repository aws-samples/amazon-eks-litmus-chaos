import React, { useState } from 'react';
import { Typography, AppBar, Toolbar } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import {
    Link,
    Outlet
} from "react-router-dom";

function NavBarComponent() {
    const [anchorEl, setAnchorEl] = useState(null);
    const open = Boolean(anchorEl);
    const handleMenuClick = (event) => {
      setAnchorEl(event.currentTarget);
    };
    const handleMenuClose = () => {
      setAnchorEl(null);
    };
  
    return (
<>
<AppBar position="static" style={{ backgroundColor: '#0f1b2a' }}>
                <Toolbar>
                    <IconButton
            size="large"
            edge="start"
            color="inherit"
            aria-label="open drawer"
            sx={{ mr: 2 }}
            id="basic-button"
            aria-controls={open ? 'basic-menu' : undefined}
            aria-haspopup="true"
            aria-expanded={open ? 'true' : undefined}
            onClick={handleMenuClick}
          >
            <MenuIcon />
          </IconButton>
          <Menu
        id="basic-menu"
        anchorEl={anchorEl}
        open={open}
        onClose={handleMenuClose}
        MenuListProps={{
          'aria-labelledby': 'basic-button',
        }}
      >
        <MenuItem onClick={handleMenuClose} component={Link} to="likes-service">Likes Service | Microservices | REST APIs</MenuItem>
        {/* <MenuItem onClick={handleMenuClose}>Shopping Cart Checkout | Microservices | Event Driven | Queues | gRPC</MenuItem>
        <MenuItem onClick={handleMenuClose}>Chat Room | Microservices | Streaming</MenuItem> */}
      </Menu>
                    <img src='/fire.png' alt="logo" style={{ width: '48px', height: '48px' }} />
                    <Typography variant="h6" style={{ marginLeft: '10px' }}>
                    Chaos Demo
                    </Typography>
                </Toolbar>
            </AppBar>
            <Outlet />
</>
    );
  }

export default NavBarComponent;
