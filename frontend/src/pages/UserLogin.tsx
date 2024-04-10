import {useEffect, useState} from 'react';
import {LoginFormPage} from '@ant-design/pro-components';
import {Divider, notification} from 'antd';
import {ReactComponent as Logo} from '../assets/dm_logo_long.svg';
import {useLocation, useNavigate} from "react-router-dom";
import GoogleButton from 'react-google-button'
import {UserAuth} from "../context/AuthContext";


type Props = {};

type LoginUserCredentials = {
    email: string,
    password: string
}


export function UserLogin(props: Props) {
    const [authenticating, setAuthenticating] = useState<boolean>(false);
    const navigate = useNavigate();
    const location = useLocation();

    const {signInWithGoogle, user} = UserAuth();

    const handleGoogleSignIn = async () => {
        if (authenticating) return

        setAuthenticating(true)
        try {
            await signInWithGoogle();
            const origin = location.state?.from?.pathname || '/';
            navigate(origin);

        } catch (e) {
            //TODO handle error
            console.log(e);
            notification.error({message: 'Could not login.'})
            setAuthenticating(false)
        }
        setAuthenticating(false)
    };

    useEffect(() => {
        if (user != null) {
            // console.log(location.state)
            const origin = location.state?.from?.pathname || '/';
            navigate(origin);
        }
    }, [user]);

    return <div style={{backgroundColor: 'white', height: 'calc(100vh - 48px)', margin: 0}}>
        <LoginFormPage
            backgroundImageUrl="https://gw.alipayobjects.com/zos/rmsportal/FfdJeJRQWjEeGTpqgBKj.png"
            logo={<Logo fill='#006d75'
                        style={{width: 400, height: 200, marginTop: -70, marginLeft: -170}}/>}
            subTitle={<>Role Based Access Control System</>}

            submitter={{submitButtonProps: {style: {display: "none"}}}}


            actions={
                <div>
                    <Divider plain>
                        <span style={{color: '#CCC', fontWeight: 'normal', fontSize: 14}}>Sign in with your business email</span>
                    </Divider>
                    <GoogleButton
                        style={{width: '100%', pointerEvents: authenticating ? 'none' : 'auto', backgroundColor: '#006d75'}}
                        disabled={authenticating}
                        onClick={handleGoogleSignIn}
                    />
                </div>
            }
        >
        </LoginFormPage>
    </div>;


}