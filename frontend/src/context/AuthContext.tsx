import {createContext, useContext, useEffect, useState} from "react";
import {
    GoogleAuthProvider,
    onAuthStateChanged,
    signInWithPopup,
    signOut
} from 'firebase/auth';
import {auth} from '../config/firebase';
import {User, UserCredential} from "@firebase/auth/dist/auth-public";
import {GlobalStateContext} from "./GlobalContext";
import axiosApiInstance from "../api/axiosClient";


//User login credentials and functionality

type UserAuthContextType = {
    user: User | undefined | null
    logout: () => Promise<void>
    signInWithGoogle: () => Promise<UserCredential>
    permissions: string[][]
    can: (action: string, object: string) => boolean
}

const UserAuthContext = createContext<UserAuthContextType>({
    user: null,
    logout: () => new Promise<void>(() => false),
    signInWithGoogle: () => new Promise<UserCredential>(() => false),
    permissions: [],
    can: (action, object) => false
})

export const UserAuthContextProvider = ({children}: any) => {
    const [user, setUser] = useState<User | undefined | null>(undefined);
    const [permissions, setPermissions] = useState<string[][]>([]);
    const {authorizing, setAuthorizing} = useContext(GlobalStateContext);

    const signInWithGoogle = () => {
        const provider = new GoogleAuthProvider();
        return signInWithPopup(auth, provider)
    };

    const logout = () => {
        return signOut(auth)
    }

    const fetchUserPermissions = async () => {
        try {
            const res = await axiosApiInstance.post('/api/casbin/permissions')
            const resPermissions = res.data || []
            setPermissions(resPermissions)
            setAuthorizing!(false)
        } catch (e: any) {
            console.log(e)
            setPermissions([])
            setAuthorizing!(false)
        }
    }

    const can = (action: string, object: string): boolean => {
        if (!permissions)
            return false

        for (const permission of permissions) {
            const permRole = permission[0]
            const permObject = permission[1]
            const permAction = permission[2]
            if (permAction == action && permObject == object)
                return true
        }

        return false
    }

    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, (currentUser) => {
            // console.log(currentUser);
            setUser(currentUser);

            if (currentUser) {
                // fetch their permissions
                fetchUserPermissions()
            } else {
                setPermissions([])
                setAuthorizing!(false)
            }

        });
        return () => {
            unsubscribe();
        };
    }, []);


    return (
        <UserAuthContext.Provider
            value={{user, logout, signInWithGoogle, permissions, can}}>
            {children}
        </UserAuthContext.Provider>
    )
}

export const UserAuth = () => {
    return useContext(UserAuthContext)
}