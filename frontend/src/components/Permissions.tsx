import React, {useEffect, useState} from 'react';
import { Button, Card, Checkbox, Col, ConfigProvider, Divider, Layout, Menu, message, notification, Popconfirm, Row, Space, Spin, Typography } from 'antd';
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import {ModalForm, PageContainer, ProFormText} from "@ant-design/pro-components";
import {Role} from "./Roles";
import enUSIntl from "antd/lib/locale/en_US";
import axiosApiInstance from "../api/axiosClient";

const {Content, Sider} = Layout;

type Props = {};

type RolePermissions = {
    [key: string]: PermissionInfo[]
}

type PermissionInfo = {
    id: number
    resource: string
    action: string
    description: string
    has_permission: boolean
}

const Permissions = (props: Props) => {

    //
    // Initialize State
    //
    const [selectedRole, setSelectedRole] = useState<string>();
    const [roles, setRoles] = useState<string[]>();
    const [rolePermissions, setRolePermissions] = useState<RolePermissions>();
    const [loadingRoles, setLoadingRoles] = useState<boolean>(false);
    const [loadingRolePermissions, setLoadingRolePermissions] = useState<boolean>(false);


    const getRoles = async () => {
        setLoadingRoles(true)
        try {
            const res = await axiosApiInstance.get<Role[]>('/api/roles/')
            let data = res.data ? res.data.map((r: any) => r.role) : []

            const selectedRole = data.length > 0 ? data[0] : undefined
            getPermissionsForRole(selectedRole)

            setRoles(data)
            setSelectedRole(selectedRole)
            setLoadingRoles(false)
        } catch (e: any) {
            setLoadingRoles(false)
            notification.error({message: e.response.data.message})
            return []
        }
    }

    const getPermissionsForRole = async (role: string) => {
        setLoadingRolePermissions(true)

        try {
            const res = await axiosApiInstance.get('/api/permissions/', {
                params: {
                    role
                }
            })
            let data = res.data || []

            setRolePermissions(data)
            setLoadingRolePermissions(false)
        } catch (e: any) {
            setLoadingRolePermissions(false)
            notification.error({message: e.response.data.message})
            return []
        }
    }

    const createRole = async (role: Role) => {
        setLoadingRoles(true)
        try {
            await axiosApiInstance.post<Role>('/api/roles/', {
                ...role
            })

            await getRoles()

            setLoadingRoles(false)
        } catch (e: any) {
            setLoadingRoles(false)
            notification.error({message: e.response.data.message})
            return []
        }
    }

    const deleteRole = async () => {
        setLoadingRoles(true)

        try {
            await axiosApiInstance.delete('/api/roles/', {
                data: {
                    role: selectedRole
                }
            })

            await getRoles()

            notification.success({message: 'Success'});
            setLoadingRoles(false)
        } catch (e: any) {
            setLoadingRoles(false)
            // notification.error({message: e.response.data.message})
            return []
        }
    }

    const handlePermissionChange = async (checked: boolean, permissionInfo: PermissionInfo) => {
        // add policy
        if (checked) {
            try {
                await axiosApiInstance.post('/api/permissions/', {
                    newRole: selectedRole,
                    newData: permissionInfo.resource,
                    newPrivilege: permissionInfo.action
                })

                await getPermissionsForRole(selectedRole!)
                notification.success({message: 'Success'})
            } catch (e: any) {
                notification.error({message: e.response.data.message})
            }
        }
        // remove policy
        else {
            try {
                await axiosApiInstance.delete('/api/permissions/', {
                    data: {
                        role: selectedRole,
                        data: permissionInfo.resource,
                        privilege: permissionInfo.action
                    }
                })

                await getPermissionsForRole(selectedRole!)
                notification.success({message: 'Success'})
            } catch (e: any) {
                notification.error({message: e.response.data.message})
            }
        }

    }

    // It's called everytime the values inside the array change (the array at the end, after the function)
    // Now that the array is empty, useEffect is called only once, at the start
    useEffect(() => {
        getRoles()
    }, [])

    return (
        <PageContainer>

            <Row justify={'end'}>

                <Space>
                    <ConfigProvider locale={enUSIntl}>
                        <ModalForm<{ role: string, description: string }>
                            layout={"vertical"}
                            trigger={
                                <Button>
                                    <PlusOutlined/>
                                    Create role
                                </Button>
                            }
                            submitter={{resetButtonProps: false}}
                            modalProps={{okText: 'Create'}}
                            onFinish={async (formData) => {
                                if (!formData.role) return false
                                if (formData.role.trim() == '') return false
                                formData.role = formData.role.trim()
                                createRole(formData)
                                message.success('Success')
                                return true
                            }}
                        >
                            <ProFormText name={'role'} label={'Role name'} required/>
                            <ProFormText name={'description'} label={'Role description'}/>
                        </ModalForm>
                    </ConfigProvider>

                    <Popconfirm title={'Are you sure you want to delete "' + selectedRole + '" role?'}
                                onConfirm={deleteRole}
                    >
                        <Button danger>
                            <DeleteOutlined/>
                            Delete role
                        </Button>
                    </Popconfirm>

                </Space>
            </Row>
            <br/>
            <Card bodyStyle={{paddingLeft: 0, paddingBottom: 0, paddingTop: 0, paddingRight: 0}}>
                <Layout className="roles-management-dashboard">
                    <Sider className="roles-management-dashboard" width={200}>
                        <Spin spinning={loadingRoles}>
                            <Menu
                                mode="inline"
                                selectedKeys={[selectedRole!]}
                                // defaultSelectedKeys={['1']}
                                // defaultOpenKeys={['sub1']}
                                style={{height: '100%', paddingBottom: 16, paddingTop: 16}}
                                onSelect={({item, key, keyPath, selectedKeys, domEvent}) => {
                                    setSelectedRole(key)
                                    getPermissionsForRole(key)
                                }}
                                theme={"dark"}

                            >
                                <Menu.ItemGroup key="roles-group" title="Roles">
                                    {roles?.map(role => {
                                        return <Menu.Item key={role}>{role}</Menu.Item>
                                    })}
                                </Menu.ItemGroup>
                            </Menu>
                        </Spin>
                    </Sider>
                    <Content style={{padding: '24px 0px 24px 24px'}}>

                        <Spin spinning={loadingRolePermissions}>
                            {rolePermissions && Object.keys(rolePermissions!).map((categoryKey) => {
                                return <>

                                    <Row justify={'start'}>
                                        <Typography.Title style={{paddingLeft: 16}}
                                                          level={5}>{categoryKey}</Typography.Title>
                                    </Row>

                                    <Divider style={{margin: "4px 0px 16px 0px"}}/>

                                    <Row style={{paddingBottom: '24px', paddingLeft: '24px'}}>
                                        {rolePermissions![categoryKey].map((info: PermissionInfo) => {
                                            return <>
                                                <Col span={8}>
                                                    <Checkbox value={info.id}
                                                              checked={info.has_permission}
                                                              onChange={e => handlePermissionChange(e.target.checked, info)}
                                                    >{info.description}</Checkbox>
                                                </Col>
                                            </>
                                        })}
                                    </Row>

                                </>
                            })}
                        </Spin>

                    </Content>
                </Layout>
            </Card>

            {/*</Content>*/}
        </PageContainer>
    );
}

export default Permissions;