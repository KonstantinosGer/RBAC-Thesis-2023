import React, {useContext} from 'react';
import './App.css';
import {Route, Routes, useLocation} from 'react-router-dom';
import {MyLayout} from "./components/MyLayout";
import PrivateRoutes from "./components/PrivateRoutes";
import {UserLogin} from "./pages/UserLogin";
import NoFoundPage from "./pages/404";
import {GlobalStateContext} from "./context/GlobalContext";
import {PageLoading} from "@ant-design/pro-layout";
import {UserAuth} from "./context/AuthContext";
import {Access} from "./components/Access";
import Home from "./components/Home";
import Roles from "./components/Roles";
import Permissions from "./components/Permissions";
import Users from "./components/Users";
import Employees from "./components/Employees";
import Customers from "./components/Customers";

function App() {
    const location = useLocation();
    const background = location.state && location.state.background;
    const {authorizing} = useContext(GlobalStateContext);
    const {can, user} = UserAuth();

    if (authorizing) {
        return <PageLoading/>
    }

    return (
        <div className="App">
            <Routes location={background || location}>
                <Route element={<PrivateRoutes/>}>


                    <Route element={<MyLayout/>}>
                        {/*////////////////*/}
                        {/*//// Routes ////*/}
                        {/*////////////////*/}
                        <Route path='/' element={
                            <Access accessible={can('read', 'rbac::data')} fallback={<></>}>
                                <Home/>
                            </Access>
                        }/>

                        <Route path='/roles' element={
                            <Access accessible={can('read', 'rbac::data')} fallback={<></>}>
                                <Roles/>
                            </Access>
                        }/>

                        <Route path='/permissions' element={
                            <Access accessible={can('read', 'rbac::data')} fallback={<></>}>
                                <Permissions/>
                            </Access>
                        }/>

                        <Route path='/users' element={
                            <Access accessible={can('read', 'rbac::data')} fallback={<></>}>
                                <Users/>
                            </Access>
                        }/>

                        <Route path='/employees' element={
                            <Access accessible={can('read', 'rbac::data')} fallback={<></>}>
                                <Employees/>
                            </Access>
                        }/>

                        <Route path='/customers' element={
                            <Access accessible={can('read', 'rbac::data')} fallback={<></>}>
                                <Customers/>
                            </Access>
                        }/>
                    </Route>

                </Route>

                <Route path='/login' element={<UserLogin/>}/>
                <Route path='*' element={<NoFoundPage/>}/>
            </Routes>
        </div>
    );
}

export default App;
