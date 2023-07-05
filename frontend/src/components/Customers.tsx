import React, {useContext, useState} from "react";
import {ConfigProvider, notification, Popconfirm} from 'antd';
import type {ProColumns} from '@ant-design/pro-table';
import {EditableProTable} from '@ant-design/pro-table';
import enUSIntl from "antd/lib/locale/en_US";
import CustomerUserRelationship from "./CustomerUserRelationship";
import {PageContainer} from "@ant-design/pro-components";
import {GlobalStateContext} from "../context/GlobalContext";
import axiosApiInstance from "../api/axiosClient";

type Props = {};

export type Customer = {
    id: string
    full_name: string
};

const Customers = (props: Props) => {
    // From GlobalStateContext
    const {
        refCustomersTable,
    } = useContext(GlobalStateContext);

    //
    // Initialize State
    //
    const [editableKeys,setEditableKeys] = useState<React.Key[]>()
    const [dataSource, setDataSource] = useState<Customer[]>();
    const [editable, setEditable] = useState<boolean>(false);

    const onSaveCustomer = async (key: any, row: Customer & { index?: number | undefined },
                            newLineConfig: Customer & { index?: number | undefined }) => {
        console.log(key, row, newLineConfig)
        setEditable(false)
        if (newLineConfig.id == '-') {
            // Create new employee
            try {
                await axiosApiInstance.post<Customer>('/api/customers/', {
                    id: +row.id,
                    full_name: row.full_name
                });
                notification.success({message: 'Success'});
                refCustomersTable?.current?.reload();
            } catch (e: any) {
                notification.error({message: e.response.data.message});
                refCustomersTable?.current?.reload();
            }

        } else {
            // Update existing employee
            try {
                await axiosApiInstance.put<Customer>('/api/customers/', {
                    ...row
                })
                notification.success({message: 'Success'});
                refCustomersTable?.current?.reload();
            } catch (e: any) {
                notification.error({message: e.response.data.message});
                refCustomersTable?.current?.reload();
            }
        }
    }


    const onDeleteCustomer = async (customer: Customer) => {
        try {
            await axiosApiInstance.delete('/api/customers/' + customer.id)
            notification.success({message: 'Success'})
            refCustomersTable?.current?.reload()
        } catch (e: any) {
            notification.error({message: e.response.data.message})
        }
    }

    const columns: ProColumns <Customer>[] = [
        {title: 'Id (Podio  Id)', dataIndex: 'id', align: "center", width: '25%', editable: () => editable},
        {title: 'Full Name', dataIndex: 'full_name', align: "center", width: '50%'},

        {
            title: 'Action',
            valueType: 'option',
            width: '25%',
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
                <Popconfirm title={'Delete this line?'} onConfirm={() => onDeleteCustomer(record)}>
                    <a key="delete">Delete</a>
                </Popconfirm>
            ],
        },
    ];


    return (
        <PageContainer>

            <ConfigProvider locale={enUSIntl}>
                <EditableProTable<Customer>
                    request={async (params, sort, filter) => {
                        try {
                            const res = await axiosApiInstance.get<Customer[]>('/api/customers/', {
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

                    actionRef={refCustomersTable}
                    columns={columns}
                    rowKey="id"
                    controlled={true}
                    value={dataSource}
                    pagination={{pageSize: 8, hideOnSinglePage: false, showQuickJumper: true}}

                    // Update dataSource
                    onDataSourceChange={(dataSource) => {setDataSource(dataSource)
                        console.log('dataSource change', dataSource)
                    }
                    }
                    recordCreatorProps={{
                        newRecordType: 'dataSource',
                        record: (index) => ({
                            id: '-',
                            full_name: '',
                        }),
                        creatorButtonText: 'Add Customer',
                        onClick: (e) => {
                            setEditable(true)
                        }
                    }}

                    editable={{
                        type: 'single',
                        editableKeys: editableKeys,
                        actionRender: (row, config, defaultDoms) => {
                            return [defaultDoms.save, defaultDoms.cancel || defaultDoms.delete];
                        },
                        onChange: (editableKeys) => setEditableKeys(editableKeys),
                        onSave: onSaveCustomer,
                        onCancel: async (key, record, originRow, newRow) => {
                            setEditable(false)
                            refCustomersTable?.current?.reload()
                        },
                        deletePopconfirmMessage: 'Delete this line?',
                        onlyOneLineEditorAlertMessage: 'Only one line can be edited at the same time',
                        onlyAddOneLineAlertMessage: 'Only add one line'
                    }}
                    bordered
                    expandable={{
                        expandedRowRender: (record) => (<CustomerUserRelationship customer_id={record.id}/>)
                    }}
                    options={{
                        search: {placeholder: 'Please enter keyword', allowClear: true},
                    }}
                />
            </ConfigProvider>


        </PageContainer>
    );
};

export default Customers;