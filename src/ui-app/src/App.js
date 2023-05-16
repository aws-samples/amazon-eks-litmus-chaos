import React from 'react';
import {
  Routes,
  Route,
  Navigate
} from "react-router-dom";
import LikesServicePage from './routing/LikesServicePage';
import NavBarComponent from './components/NavBarComponent';

const App = () => {
  return (
    <>
      <div style={{ backgroundColor: '#e9ebed', minHeight: '100vh' }}>
        <Routes>
          <Route path="/" element={<NavBarComponent />}>
            <Route index element={<Navigate replace to="likes-service" />} />
            <Route path="likes-service" element={<LikesServicePage />} />
          </Route>
        </Routes>
      </div>
    </>
  );
};

export default App;