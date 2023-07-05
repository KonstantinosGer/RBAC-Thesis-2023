import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import GlobalStateProvider from "./context/GlobalContext";
import {UserAuthContextProvider} from "./context/AuthContext";
import {ConfigProvider} from "antd";
import enUS from "antd/es/locale/en_US";
import {BrowserRouter} from "react-router-dom";

ReactDOM.render(
    <React.StrictMode>
        <GlobalStateProvider>
            <UserAuthContextProvider>
                <ConfigProvider locale={enUS}>
                    <BrowserRouter>
                        <App/>
                    </BrowserRouter>
                </ConfigProvider>
            </UserAuthContextProvider>
        </GlobalStateProvider>

    </React.StrictMode>,
    document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
