import React from 'react';
import { Navigate } from 'react-router-dom';
import { auth } from '../../config/firebase';
import { useLocation } from 'react-router-dom'

export interface IAuthRouteProps {
    // component: React.ComponentType
}

const AuthRoute: React.FunctionComponent<IAuthRouteProps> = (props) => { // { component: RouteComponent}
    const { children } = props;
    const location = useLocation();

    console.log(auth)
    if (!auth.currentUser) {
        console.log('No user detected, redirecting');
        return <Navigate to="/login" replace state={{ path: location.pathname }}/>;
    } else {
        console.log('Passed auth route as a user was found');
        return <div>{children}</div>;
        // return <RouteComponent/>
    }
};

export default AuthRoute;
