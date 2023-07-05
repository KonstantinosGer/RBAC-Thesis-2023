import Home from "../components/Home";
import Roles from "../components/Roles";
import {Route, Routes} from "react-router-dom";
import React from "react";
import AuthRoute from "../components/AuthRoute";
import Permissions from "../components/Permissions";
import Employees from "../components/Employees";
import Customers from "../components/Customers";
import Users from "../components/Users";


export function MyRoutes() {
    return (

        <Routes>
            <Route path='/' element={
                <AuthRoute><Home/></AuthRoute>
            }/>

            <Route path='/users' element={
                <AuthRoute><Users/></AuthRoute>
            }/>

            <Route path='/employees' element={
                <AuthRoute><Employees/></AuthRoute>
            }/>

            <Route path='/customers' element={
                <AuthRoute><Customers/></AuthRoute>
            }/>

            <Route path='/roles' element={
                <AuthRoute><Roles/></AuthRoute>
            }/>

            <Route path='/permissions' element={
                <AuthRoute><Permissions/></AuthRoute>
            }/>
        </Routes>

    );
}