import * as React from 'react';

type Props = {
    accessible: boolean
    fallback: JSX.Element
    children: JSX.Element
};

//Helper wrapper component used for authorization purposes
export const Access = (props: Props): JSX.Element => {

    if (!props.accessible) {
        return props.fallback
    }

    return props.children
};