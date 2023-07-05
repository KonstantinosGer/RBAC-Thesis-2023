import React, {createContext, Dispatch, SetStateAction, useRef, useState} from 'react';
import {ActionType} from "@ant-design/pro-components";

//create a context, with createContext api
type ContextType = {
    authorizing?: boolean
    setAuthorizing?: Dispatch<SetStateAction<boolean>>
    refUsersTable?: React.MutableRefObject<ActionType | undefined>
    refEmployeesTable?: React.MutableRefObject<ActionType | undefined>
    refCustomersTable?: React.MutableRefObject<ActionType | undefined>
    refEmployeeUserTable?: React.MutableRefObject<ActionType | undefined>
    refCustomerUserTable?: React.MutableRefObject<ActionType | undefined>
}

export const GlobalStateContext = createContext<ContextType>({});


const GlobalStateProvider = ({children}: any) => {
    // this state will be shared with all components
    const [authorizing, setAuthorizing] = useState<boolean>(true);
    const refUsersTable = useRef<ActionType>();
    const refEmployeesTable = useRef<ActionType>();
    const refCustomersTable = useRef<ActionType>();
    const refEmployeeUserTable = useRef<ActionType>();
    const refCustomerUserTable = useRef<ActionType>();

    return (
        // this is the provider providing state
        <GlobalStateContext.Provider value={{
            authorizing,
            setAuthorizing,
            refUsersTable,
            refEmployeesTable,
            refCustomersTable,
            refEmployeeUserTable,
            refCustomerUserTable,
        }}>
            {children}
        </GlobalStateContext.Provider>
    );
};

export default GlobalStateProvider;