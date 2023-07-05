import {Navigate, Outlet, useLocation} from 'react-router-dom'
import {PageLoading} from "@ant-design/pro-layout";
import {UserAuth} from "../context/AuthContext";

//Routes that require users to be authenticated
const PrivateRoutes = () => {
    const {user, permissions, can} = UserAuth();
    const location = useLocation();

    // handle initial state of user (while in the process of authenticating...)
    if (user === undefined) {
        return <PageLoading/>
    }

    return (
        user != null ? <Outlet/> : <Navigate to="/login" replace state={{from: location}}/>
    )
}

export default PrivateRoutes