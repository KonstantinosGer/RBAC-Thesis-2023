import * as React from 'react';
import {useContext, useState} from 'react';
import {ProLayout} from "@ant-design/pro-layout";
import {
    ClockCircleOutlined,
    HomeOutlined,
    SafetyCertificateOutlined,
    SolutionOutlined,
    TeamOutlined,
    UserOutlined,
} from "@ant-design/icons";
import {ProSettings} from "@ant-design/pro-components";
import enUS from "antd/es/locale/en_US";
import {Avatar, Button, ConfigProvider, Dropdown, notification} from "antd";
import {Link, Outlet, useLocation} from "react-router-dom";
import {ReactComponent as Logo} from '../assets/dm_logo_long.svg';
import {UserAuth} from "../context/AuthContext";
import {AvatarDropdownMenu} from "./AvatarDropdownMenu";
import {GlobalStateContext} from "../context/GlobalContext";
import axiosApiInstance from "../api/axiosClient";

type Props = {};

export const MyLayout = (props: Props) => {
    // From GlobalStateContext
    const {
        refUsersTable,
        refEmployeesTable,
        refCustomersTable,
    } = useContext(GlobalStateContext);

    //
    // Initialize State
    //
    const [settings, setSetting] = useState<Partial<ProSettings> | undefined>({
        fixSiderbar: true,
        layout: 'mix',
        splitMenus: true,
    });

    const location = useLocation();
    const {can, user} = UserAuth();
    const [syncing, setSyncing] = useState<boolean>(false);

    const onSyncWithFirebase = async () => {
        try {
            setSyncing(true)
            notification.info({message: 'Started update'})
            await axiosApiInstance.get('/api/users/sync')
            notification.success({message: 'Successfully updated', duration: 0})
            setSyncing(false)

            //Refresh all tables when sync button is pressed
            refUsersTable?.current?.reload()
            refEmployeesTable?.current?.reload()
            refCustomersTable?.current?.reload()
        } catch (e: any) {
            notification.error({message: e.response.data.message})
        }
    }

    return (
        <div id="test-pro-layout">
            <ConfigProvider locale={enUS}>

                <ProLayout
                    location={location}
                    menuItemRender={(item: any, defaultDom: any) => {
                        return <Link to={item.path}> {defaultDom} </Link>

                    }}
                    subMenuItemRender={(item: any, defaultDom: any) => <Link to={item.path}> {defaultDom} </Link>}
                    route={{
                        path: '/',
                        // name: 'Home',
                        routes: [
                            {
                                path: '/',
                                name: 'Home',
                                icon: <HomeOutlined/>,
                                // component: './Welcome',
                            },
                            {
                                path: '/users',
                                name: 'Users',
                                icon: <TeamOutlined/>,
                            },
                            {
                                path: '/employees',
                                name: 'Employees',
                                icon: <TeamOutlined/>,
                            },
                            {
                                path: '/customers',
                                name: 'Customers',
                                icon: <TeamOutlined/>,
                            },
                            {
                                path: '/roles',
                                name: 'Roles',
                                icon: <SolutionOutlined/>,
                            },
                            {
                                path: '/permissions',
                                name: 'Permissions',
                                icon: <SafetyCertificateOutlined/>,
                            },
                        ]
                    }}
                    appList={
                        [
                            {
                                icon: <></>,
                                title: 'Portal',
                                desc: 'Customer and Internal Portal',
                                url: 'http://localhost:3000/',
                                target: '_blank',
                            },
                        ]
                    }
                    siderMenuType="group"
                    menu={{
                        collapsedShowGroupTitle: true,
                    }}
                    actionsRender={(props) => {
                        if (props.isMobile) return [];


                        console.log(user?.photoURL)
                        // let avatar: JSX.Element
                        let avatar = <Avatar style={{backgroundColor: '#006d75'}} icon={<UserOutlined/>}/>

                        if (user?.photoURL != null) {
                            avatar = <Avatar src={user?.photoURL}/>
                        }

                        return [

                            <Button icon={<ClockCircleOutlined/>} disabled={syncing} onClick={() => {
                                //Refresh table's data from google drive
                                onSyncWithFirebase()
                            }}>Sync with Firebase</Button>,

                            <div style={{
                                display: 'flex',
                                alignItems: 'center',
                                marginInlineEnd: 32,
                            }}>
                                <Dropdown overlay={<AvatarDropdownMenu/>}>
                                    {avatar}
                                </Dropdown>
                            </div>
                        ];
                    }}

                    // footerRender={() => {
                    //     return (
                    //         <div
                    //             style={{
                    //                 textAlign: 'center',
                    //                 paddingBlockStart: 12,
                    //                 // flexShrink: 0,
                    //             }}
                    //         >
                    //             <div>© 2023 Digital Minds</div>
                    //         </div>
                    //     );
                    // }}


                    // menuFooterRender={(props) => {
                    //     if (props?.collapsed) return undefined;
                    //     return (
                    //         <div
                    //             style={{
                    //                 textAlign: 'center',
                    //                 paddingBlockStart: 12,
                    //             }}
                    //         >
                    //             <div>© 2023 Digital Minds</div>
                    //         </div>
                    //     );
                    // }}


                    logo={
                        <Link to="/">
                            <Logo fill='#006d75' width={170}
                                  style={{marginBottom: -20, marginLeft: -20, marginRight: -4}}
                            />
                        </Link>
                    }

                    title={'RBAC'}

                    // style={{display: "flex", flexDirection: "column", height: "100vh", flexGrow: 1}}

                    {...settings}
                >

                    <Outlet/>

                </ProLayout>

            </ConfigProvider>

        </div>
    );
};

