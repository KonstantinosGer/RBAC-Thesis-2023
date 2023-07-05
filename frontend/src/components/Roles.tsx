import React, {useEffect, useRef, useState} from "react";
import {Col, ConfigProvider, notification, Row} from 'antd';
import type {ActionType, ProColumns} from '@ant-design/pro-table';
import {EditableProTable} from '@ant-design/pro-table';
import enUSIntl from "antd/lib/locale/en_US";
import {PageContainer} from "@ant-design/pro-components";
import axiosApiInstance from "../api/axiosClient";
import {UserEmail} from "./CustomerUserRelationship";

type Props = {};

type User = {
    id: string
    full_name: string
    email: string
    role: string
};

export type Role = {
    role: string
    description: string
};


const Roles = (props: Props) => {

    //
    // Initialize State
    //
    const [editableKeys, setEditableKeys] = useState<React.Key[]>()
    const [dataSource, setDataSource] = useState<User[]>();
    const refUsersTable = useRef<ActionType>();

    const [roles, setRoles] = useState<UserEmail[]>();

    const getRoles = async () => {
        //data has key and value
        //must be on .json format
        try {
            const res = await axiosApiInstance.get<Role[]>('/api/roles/')
            let data = res.data || []
            let selectOptions = data.map((r: Role) => (
                {label: r.role, value: r.role}
            ))

            setRoles(selectOptions)

            return selectOptions

        } catch (e: any) {
            notification.error({message: e.response.data.message})
            //Το return να είναι της μορφής αυτού που όντως επιστρέφεται
            return []
        }
    }


    const onSaveUserRole = async (key: any, row: User & { index?: number | undefined },
                                  newLineConfig: User & { index?: number | undefined }) => {
        try {
            await axiosApiInstance.put('/api/roles/', {
                ...row
            })
            notification.success({message: 'Success'})
        } catch (e: any) {
            notification.error({message: e.response.data.message})
            refUsersTable.current?.reload()
        }
    }


    const columns: ProColumns<User>[] = [
        {title: 'User id', dataIndex: 'id', editable: false, align: "center"},
        {title: 'Full name', dataIndex: 'full_name', editable: false, align: "center"},
        {title: 'Email', dataIndex: 'email', editable: false, align: "center"},
        {
            title: 'Role', dataIndex: 'role', valueType: 'select', align: "center",
            // request: getRoles,
            fieldProps: {
                // showSearch: true,
                // placeholder: "Select a role",
                optionFilterProp: "children",
                filterOption: (input: string, option: { label: any; }) => (option?.label ?? '').toLowerCase().includes(input.toLowerCase()),
                options: roles,
            }
        },

        {
            title: 'Action',
            valueType: 'option',
            align: "center",
            render: (text, record, _, action) => [
                <a
                    key="editable"
                    onClick={() => {
                        action?.startEditable?.(record.id);
                    }}
                >
                    Edit
                </a>,
            ],
        },
    ];

    // It's called everytime the values inside the array change (the array at the end, after the function)
    // Now that the array is empty, useEffect is called only once, at the start
    useEffect(() => {
        getRoles()
    }, [])


    return (
        <PageContainer>
            <Row>
                <Col span={24}>
                    <ConfigProvider locale={enUSIntl}>
                        <EditableProTable<User>
                            request={async (params, sort, filter) => {

                                try {
                                    const res = await axiosApiInstance.get('/api/users/', {
                                        params: {
                                            keyword: params.keyword
                                        }
                                    })
                                    return {data: res.data, success: true, total: res.data.length}
                                } catch (e: any) {
                                    notification.error({message: e.response.data.message})
                                    //Το return να είναι της μορφής αυτού που όντως επιστρέφεται
                                    return {data: [], success: false, total: 0}
                                }

                            }}

                            actionRef={refUsersTable}
                            columns={columns}
                            rowKey="id"
                            controlled={true}
                            value={dataSource}
                            onChange={(dataSource) => setDataSource(dataSource as any)}

                            //Το βάζω false αν δε θέλω να προσθέτει γραμμή
                            recordCreatorProps={false}

                            pagination={{pageSize: 8, hideOnSinglePage: true, showQuickJumper: true}}
                            editable={{
                                type: 'single',
                                editableKeys: editableKeys,
                                actionRender: (row, config, defaultDoms) => {
                                    // return [defaultDoms.save, defaultDoms.delete || defaultDoms.cancel];
                                    return [defaultDoms.save, defaultDoms.cancel];
                                },
                                onChange: (editableKeys) => setEditableKeys(editableKeys),
                                onSave: onSaveUserRole,
                                onCancel: async (key, record, originRow, newRow) => refUsersTable.current?.reload(),
                                //deletePopconfirmMessage: 'Delete this line?',
                                onlyOneLineEditorAlertMessage: 'Only one line can be edited at the same time',
                                //onlyAddOneLineAlertMessage: 'Only add one line'
                            }}
                            options={{
                                search: {placeholder: 'Please enter keyword', allowClear: true},
                            }}
                            bordered
                        />
                    </ConfigProvider>
                </Col>
            </Row>

        </PageContainer>
    );

};

export default Roles;