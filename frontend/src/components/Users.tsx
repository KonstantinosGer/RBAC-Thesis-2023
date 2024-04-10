import React, {useState, useContext} from "react";
import {Button, ConfigProvider, message, notification, Popconfirm, Row} from 'antd';
import type {ProColumns} from '@ant-design/pro-table';
import {EditableProTable} from '@ant-design/pro-table';
import enUSIntl from "antd/lib/locale/en_US";
import {ModalForm, PageContainer, ProFormGroup, ProFormText} from "@ant-design/pro-components";
import {PlusOutlined} from "@ant-design/icons";
import {GlobalStateContext} from "../context/GlobalContext";
import axiosApiInstance from "..//api/axiosClient";

type Props = {};

export type User = {
    id: string
    email: string
};

export type UserToAdd = {
    id: string
    password: string
};

const Users = (props: Props) => {
    // From GlobalStateContext
    const {
        refUsersTable,
    } = useContext(GlobalStateContext);

    //
    // Initialize State
    //
    const [dataSource, setDataSource] = useState<User[]>();

    const onDeleteUser = async (user: User) => {
        try {
            await axiosApiInstance.delete('/api/firebase/' + user.id);
            notification.success({message: 'Success'});
            refUsersTable?.current?.reload();
        } catch (e: any) {
            notification.error({message: e.response.data.message});
        }
    }

    const columns: ProColumns<User>[] = [
        {title: 'Id (Firebase Id)', dataIndex: 'id', align: "center", editable: false},
        {title: 'Email', dataIndex: 'email', align: "center"},

        {
            title: 'Action',
            valueType: 'option',
            width: 200,
            align: "center",
            render: (text, record, _, action) => [
                <Popconfirm title={'Delete this line?'} onConfirm={() => onDeleteUser(record)}>
                    <a key="delete">Delete</a>
                </Popconfirm>
            ],
        },
    ];


    return (
        <PageContainer>

            <ConfigProvider locale={enUSIntl}>
                <Row justify={"end"}>
                    <ModalForm<{ email: string, password: string }>
                        title="Add new user to firebase"
                        layout={"vertical"}
                        trigger={
                            <Button>
                                <PlusOutlined/>
                                Add new user
                            </Button>
                        }
                        autoFocusFirstInput
                        submitter={{resetButtonProps: false}}
                        modalProps={{okText: 'Add user'}}

                        // Add new user
                        onFinish={async (values) => {
                            try {
                                await axiosApiInstance.post<UserToAdd>('/api/firebase/', {
                                    ...values
                                });
                                notification.success({message: 'Success'});
                                refUsersTable?.current?.reload();
                                return true;
                            } catch (e: any) {
                                notification.error({message: e.response.data.message});
                            }
                        }}

                        width={550}
                    >
                        <ProFormGroup align={"center"}>
                            <ProFormText name={'email'} label={'Email'} width="md" rules={[{required: true}]}/>
                            <ProFormText name={'password'} label={'Password'} width="md" rules={[{required: true}]}/>
                        </ProFormGroup>
                    </ModalForm>
                </Row>
                <br/>
                <Row>
                    <EditableProTable<User>
                        request={async (params, sort, filter) => {
                            try {
                                const res = await axiosApiInstance.get<User[]>('/api/firebase/', {
                                    params: {
                                        keyword: params.keyword
                                    }
                                })
                                return {data: res.data, success: true, total: res.data.length}
                            } catch (e: any) {
                                notification.error({message: e.response.data.message})
                                return {data: [], success: false, total: 0}
                            }
                        }}

                        actionRef={refUsersTable}
                        columns={columns}
                        rowKey="id"
                        controlled={true}
                        value={dataSource}
                        onChange={(dataSource) => setDataSource(dataSource as any)}
                        pagination={{pageSize: 8, hideOnSinglePage: false, showQuickJumper: true}}
                        options={{
                            search: {placeholder: 'Please enter keyword', allowClear: true},
                        }}
                        //Το βάζω false αν δε θέλω να προσθέτει γραμμή
                        recordCreatorProps={false}
                        bordered
                    />
                </Row>

            </ConfigProvider>

        </PageContainer>
    );

};

export default Users;