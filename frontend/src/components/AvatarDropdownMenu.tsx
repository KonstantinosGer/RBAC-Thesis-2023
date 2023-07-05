import {Menu} from "antd";
import {LogoutOutlined} from "@ant-design/icons";
import * as React from "react";
import {useNavigate} from "react-router-dom";
import {UserAuth} from "../context/AuthContext";

export const AvatarDropdownMenu = () => {

    const {user, logout} = UserAuth();
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            await logout();
            navigate('/login');
            console.log('You are logged out')
        } catch (err: any) {
            //TODO handle error
            console.log(err.message);
        }
    };

    return <>
        <Menu
            items={[
                {
                    key: 'user.name',
                    // icon: <UserOutlined/>,
                    label: user?.displayName,
                    style: {pointerEvents: "none"}
                },
                {
                    key: 'user.email',
                    // icon: <UserOutlined/>,
                    label: user?.email,
                    style: {pointerEvents: "none", color: '#9d9a9a', fontSize: 14}
                },
                {
                    type: 'divider' as const,
                },
                {
                    key: 'logout',
                    icon: <LogoutOutlined/>,
                    label: 'Logout',
                    onClick: handleLogout
                },
            ]}
        />
    </>
}