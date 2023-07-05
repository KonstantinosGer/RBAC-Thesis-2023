import React, {useContext, useEffect, useState} from "react";
import {ConfigProvider, notification, Popconfirm, Switch} from 'antd';
import type {ProColumns} from '@ant-design/pro-table';
import {EditableProTable} from '@ant-design/pro-table';
import enUSIntl from "antd/lib/locale/en_US";
import {GlobalStateContext} from "../context/GlobalContext";
import axiosApiInstance from "../api/axiosClient";

type Props = { customer_id: string };

export type CustomerUser = {
    id: string
    email: string
    has_performance_access: boolean
    has_financial_access: boolean
};

export type UserEmail = {
    value: string
    label: string
};


const CustomerUserRelationship = (props: Props) => {
    // From GlobalStateContext
    const {
        refCustomerUserTable,
    } = useContext(GlobalStateContext);

    //
    // Initialize State
    //
    const [editableKeys, setEditableKeys] = useState<React.Key[]>()
    const [dataSource, setDataSource] = useState<CustomerUser[]>();
    const [userEmails, setUserEmails] = useState<UserEmail[]>();

    const onAddCustomerUserAssociation = async (key: any, row: CustomerUser & { index?: number | undefined },
                                                newLineConfig: CustomerUser & { index?: number | undefined }) => {
        try {
            await axiosApiInstance.post<CustomerUser>('/api/customers/associations/', {
                id: +props.customer_id,
                email: row.email,
            })
            notification.success({message: 'Success'})
            refCustomerUserTable?.current?.reload()

        } catch (e: any) {
            notification.error({message: e.response.data.message})
            refCustomerUserTable?.current?.reload()
        }

    }


    const onDeleteCustomerUserAssociation = async (customerUser: CustomerUser) => {
        try {
            await axiosApiInstance.delete('/api/customers/associations/', {
                data: {
                    customer_id: +props.customer_id,
                    user_id: customerUser.id
                }
            })
            notification.success({message: 'Success'})
            refCustomerUserTable?.current?.reload()
        } catch (e: any) {
            notification.error({message: e.response.data.message})
            refCustomerUserTable?.current?.reload()
        }
    }

    const getUserEmails = async () => {

        try {
            const res = await axiosApiInstance.get('/api/users/emails')
            let data = res.data || []
            let selectOptions = data.map((email: string) => (
                {label: email, value: email}
            ))

            setUserEmails(selectOptions)

            return selectOptions

        } catch (e: any) {
            notification.error({message: e.response.data.message})
            //Το return να είναι της μορφής αυτού που όντως επιστρέφεται
            return []
        }
    }

    const onToggleAccessSwitch = async (customerUser: CustomerUser, performanceOrFinance: string, checked: boolean) => {
        try {
            await axiosApiInstance.put('/api/customers/associations/', {
                customer_id: +props.customer_id,
                user_id: customerUser.id,
                access_object: performanceOrFinance,
                has_access: checked
            })
            notification.success({message: 'Success'})
            refCustomerUserTable?.current?.reload()
        } catch (e: any) {
            notification.error({message: e.response.data.message})
        }
    }

    const columns: ProColumns<CustomerUser>[] = [
        {title: 'User id', dataIndex: 'id', editable: false, width: '26%', align: "center"},
        {
            title: 'Email', dataIndex: 'email', valueType: 'select', align: "center", width: '26%',
            fieldProps: {
                showSearch: true,
                placeholder: "Select an email",
                optionFilterProp: "children",
                filterOption: (input: string, option: { label: any; }) => (option?.label ?? '').toLowerCase().includes(input.toLowerCase()),
                options: userEmails,
            }
        },
        {
            title: 'Performance Access', dataIndex: 'has_performance_access', editable: false, valueType: "switch",
            align: "center", width: '16%',
            render: (text, record) => (
                <Switch checked={record.has_performance_access}
                        onClick={async (checked) => {
                            onToggleAccessSwitch(record, "performance", checked)
                        }}
                />)
        },
        {
            title: 'Financial Access', dataIndex: 'has_financial_access', editable: false, valueType: "switch",
            align: "center", width: '16%',
            render: (text, record) => (
                <Switch checked={record.has_financial_access}
                        onClick={async (checked) => {
                            onToggleAccessSwitch(record, "finance", checked)
                        }}
                />)
        },

        {
            title: 'Action',
            valueType: 'option',
            width: '16%',
            align: "center",
            render: (text, record, _, action) => [
                <Popconfirm title={'Delete this line?'} onConfirm={() => onDeleteCustomerUserAssociation(record)}>
                    <a key="delete">Delete</a>
                </Popconfirm>
            ],
        },
    ];

    // It's called everytime the values inside the array change (the array at the end, after the function)
    // Now that the array is empty, useEffect is called only once, at the start
    useEffect(() => {
        getUserEmails()
    }, [])

    return (
        <ConfigProvider locale={enUSIntl}>
            <EditableProTable<CustomerUser>
                request={async (params, sort, filter) => {
                    try {
                        const res = await axiosApiInstance.get<CustomerUser[]>('/api/customers/associations/' + props.customer_id)
                        const data = res.data || []
                        return {data, success: true, total: data.length}
                    } catch (e: any) {
                        notification.error({message: e.response.data.message})
                        return {data: [], success: false, total: 0}
                    }
                }}

                actionRef={refCustomerUserTable}
                columns={columns}
                rowKey="id"
                controlled={true}
                value={dataSource}
                onChange={(dataSource) => setDataSource(dataSource as any)}
                recordCreatorProps={{
                    newRecordType: 'dataSource',
                    record: (index) => ({
                        id: '-',
                        email: '',
                        has_performance_access: false,
                        has_financial_access: false,
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
                    onSave: onAddCustomerUserAssociation,
                    onCancel: async (key, record, originRow, newRow) => refCustomerUserTable?.current?.reload(),
                    deletePopconfirmMessage: 'Delete this line?',
                    onlyOneLineEditorAlertMessage: 'Only one line can be edited at the same time',
                    onlyAddOneLineAlertMessage: 'Only add one line'
                }}
                bordered
            />
        </ConfigProvider>
    );

};

export default CustomerUserRelationship;
