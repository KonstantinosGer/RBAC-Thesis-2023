import React, {useContext, useState} from "react";
import {ConfigProvider, notification, Popconfirm} from 'antd';
import type {ProColumns} from '@ant-design/pro-table';
import {EditableProTable} from '@ant-design/pro-table';
import enUSIntl from "antd/lib/locale/en_US";
import EmployeeUserRelationship from "./EmployeeUserRelationship";
import {PageContainer} from "@ant-design/pro-components";
import {GlobalStateContext} from "../context/GlobalContext";
import axiosApiInstance from "../api/axiosClient";

type Props = {};

export type Employee = {
    id: string
    full_name: string
};

const Employees = (props: Props) => {
    // From GlobalStateContext
    const {
        refEmployeesTable,
    } = useContext(GlobalStateContext);

    //
    // Initialize State
    //
    const [editableKeys,setEditableKeys] = useState<React.Key[]>()
    const [dataSource, setDataSource] = useState<Employee[]>();


    const onSaveEmployee = async (key: any, row: Employee & { index?: number | undefined },
                            newLineConfig: Employee & { index?: number | undefined }) => {
        if (row.id == '-') {
            // create new employee
            try {
                await axiosApiInstance.post<Employee>('/api/employees/', {
                    full_name: row.full_name
                });
                notification.success({message: 'Success'});
                refEmployeesTable?.current?.reload();
            } catch (e: any) {
                notification.error({message: e.response.data.message});
                refEmployeesTable?.current?.reload();
            }

        } else {
            // update existing employee
            try {
                await axiosApiInstance.put<Employee>('/api/employees/', {
                    ...row
                })
                notification.success({message: 'Success'});
                refEmployeesTable?.current?.reload();
            } catch (e: any) {
                notification.error({message: e.response.data.message});
                refEmployeesTable?.current?.reload();
            }
        }


    }


    const onDeleteEmployee = async (employee: Employee) => {
        try {
            await axiosApiInstance.delete('/api/employees/' + employee.id)
            notification.success({message: 'Success'})
            refEmployeesTable?.current?.reload()
        } catch (e: any) {
            notification.error({message: e.response.data.message})
        }
    }

    const columns: ProColumns <Employee>[] = [
        {title: 'Id', dataIndex: 'id', align: "center", editable: false},
        {title: 'Full name', dataIndex: 'full_name', align: "center"},

        {
            title: 'Action',
            valueType: 'option',
            width: 200,
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
                <Popconfirm title={'Delete this line?'} onConfirm={() => onDeleteEmployee(record)}>
                    <a key="delete">Delete</a>
                </Popconfirm>
            ],
        },
    ];


    return (
        <PageContainer>

            <ConfigProvider locale={enUSIntl}>
                <EditableProTable<Employee>
                    request={async (params, sort, filter) => {
                        try {
                            const res = await axiosApiInstance.get<Employee[]>('/api/employees/', {
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

                    actionRef={refEmployeesTable}
                    columns={columns}
                    rowKey="id"
                    controlled={true}
                    value={dataSource}
                    onChange={(dataSource) => setDataSource(dataSource as any)}
                    recordCreatorProps={{
                        newRecordType: 'dataSource',
                        record: (index) => ({
                            id: '-', //parseInt((Math.random() * 1000000).toFixed(0))
                            full_name: '',
                        }),
                        creatorButtonText: 'Add Employee',
                        onClick: (e) => {
                            // console.log('click', e)
                        }
                    }}
                    pagination={{pageSize: 8, hideOnSinglePage: false, showQuickJumper: true}}
                    editable={{
                        type: 'single',
                        editableKeys: editableKeys,
                        actionRender: (row, config, defaultDoms) => {
                            return [defaultDoms.save, defaultDoms.cancel || defaultDoms.delete];
                        },
                        onChange: (editableKeys) => setEditableKeys(editableKeys),
                        onSave: onSaveEmployee,
                        onCancel: async (key, record, originRow, newRow) => refEmployeesTable?.current?.reload(),
                        deletePopconfirmMessage: 'Delete this line?',
                        onlyOneLineEditorAlertMessage: 'Only one line can be edited at the same time',
                        onlyAddOneLineAlertMessage: 'Only add one line'
                    }}
                    bordered
                    expandable={{
                        expandedRowRender: (record) => (<EmployeeUserRelationship employee_id={record.id}/>)
                    }}
                    options={{
                        search: {placeholder: 'Please enter keyword', allowClear: true},
                    }}
                />
            </ConfigProvider>


        </PageContainer>
    );


};

export default Employees;