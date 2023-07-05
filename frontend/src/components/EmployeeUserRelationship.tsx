import React, {useContext, useEffect, useState} from "react";
import {ConfigProvider, notification, Popconfirm} from 'antd';
import type {ProColumns} from '@ant-design/pro-table';
import {EditableProTable} from '@ant-design/pro-table';
import enUSIntl from "antd/lib/locale/en_US";
import {GlobalStateContext} from "../context/GlobalContext";
import axiosApiInstance from "../api/axiosClient";

type Props = {
    employee_id: string
};

export type EmployeeUser = {
    id: string
    email: string
};

export type UnassignedEmail = {
    value: string
    label: string
};

const EmployeeUserRelationship = (props: Props) => {
    // From GlobalStateContext
    const {
        refEmployeeUserTable,
    } = useContext(GlobalStateContext);

    //
    // Initialize State
    //
    const [editableKeys, setEditableKeys] = useState<React.Key[]>()
    const [dataSource, setDataSource] = useState<EmployeeUser[]>();
    const [unsignedEmails, setUnsignedEmails] = useState<UnassignedEmail[]>();

    const onAddEmployeeUserAssociation = async (key: any, row: EmployeeUser & { index?: number | undefined },
                                                newLineConfig: EmployeeUser & { index?: number | undefined }) => {
        try {
            await axiosApiInstance.post<EmployeeUser>('/api/employees/associations/', {
                employee_id: +props.employee_id,
                user_email: row.email,
            })
            refEmployeeUserTable?.current?.reload()
            notification.success({message: 'Success'})
        } catch (e: any) {
            notification.error({message: e.response.data.message})
            refEmployeeUserTable?.current?.reload()
        }

    }


    const onDeleteEmployeeUserAssociation = async (employeeUser: EmployeeUser) => {
        try {
            await axiosApiInstance.delete('/api/employees/associations/', {
                data: {
                    employee_id: +props.employee_id,
                    user_id: employeeUser.id
                }
            })
            notification.success({message: 'Success'})
            refEmployeeUserTable?.current?.reload()
        } catch (e: any) {
            notification.error({message: e.response.data.message})
            refEmployeeUserTable?.current?.reload()
        }
    }

    const getUnassignedEmailsRequest = async () => {
        try {
            const res = await axiosApiInstance.get('/api/users/unassigned')
            let data = res.data || []
            let selectOptions = data.map((email: string) => (
                {label: email, value: email}
            ))

            setUnsignedEmails(selectOptions)

            return selectOptions

        } catch (e: any) {
            notification.error({message: e.response.data.message})
            //Το return να είναι της μορφής αυτού που όντως επιστρέφεται
            return []
        }
    }

    const columns: ProColumns<EmployeeUser>[] = [
        {title: 'User id', dataIndex: 'id', editable: false, align: "center"},
        {
            title: 'Email', dataIndex: 'email', valueType: 'select', align: "center",
            fieldProps: {
                showSearch: true,
                placeholder: "Select an email",
                optionFilterProp: "children",
                filterOption: (input: string, option: { label: any; }) => (option?.label ?? '').toLowerCase().includes(input.toLowerCase()),
                options: unsignedEmails,
            }
        },
        {
            title: 'Action',
            valueType: 'option',
            width: 200,
            align: "center",
            render: (text, record, _, action) => [
                <Popconfirm title={'Delete this line?'} onConfirm={() => onDeleteEmployeeUserAssociation(record)}>
                    <a key="delete">Delete</a>
                </Popconfirm>
            ],
        },
    ];

    // It's called everytime the values inside the array change (the array at the end, after the function)
    // Now that the array is empty, useEffect is called only once, at the start
    useEffect(() => {
        getUnassignedEmailsRequest()
    }, [])

    return (
        <ConfigProvider locale={enUSIntl}>
            <EditableProTable<EmployeeUser>
                request={async (params, sort, filter) => {
                    try {
                        const res = await axiosApiInstance.get<EmployeeUser[]>('/api/employees/associations/' + props.employee_id)
                        const data = res.data || []
                        return {data, success: true, total: data.length}
                    } catch (e: any) {
                        notification.error({message: e.response.data.message})
                        return {data: [], success: false, total: 0}
                    }
                }}

                actionRef={refEmployeeUserTable}
                columns={columns}
                rowKey="id"
                controlled={true}
                value={dataSource}
                onChange={(dataSource) => setDataSource(dataSource as any)}
                recordCreatorProps={{
                    newRecordType: 'dataSource',
                    record: (index) => ({
                        id: '-', //(Math.random() * 1000000).toFixed(0)
                        email: '',
                    }),
                    creatorButtonText: 'Add Associated User',
                }}
                pagination={{pageSize: 30, hideOnSinglePage: true}}
                editable={{
                    type: 'single',
                    editableKeys: editableKeys,
                    actionRender: (row, config, defaultDoms) => {
                        return [defaultDoms.save, defaultDoms.cancel || defaultDoms.delete];
                    },
                    onChange: (editableKeys) => setEditableKeys(editableKeys),
                    onSave: onAddEmployeeUserAssociation,
                    onCancel: async (key, record, originRow, newRow) => refEmployeeUserTable?.current?.reload(),
                    deletePopconfirmMessage: 'Delete this line?',
                    onlyOneLineEditorAlertMessage: 'Only one line can be edited at the same time',
                    onlyAddOneLineAlertMessage: 'Only add one line'
                }}
                bordered
            />
        </ConfigProvider>
    );

};

export default EmployeeUserRelationship;