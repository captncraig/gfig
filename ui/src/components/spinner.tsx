import * as React from 'react';
import { GridLoader as Loader } from 'halogen'

export class Spinner extends React.Component<{}, {}>{
    public render() {
        const s = { textAlign: "center" };
        const s2 = { margin: "0 auto", width: "60px" };
        return <div >
            <h1 style={s}>LOADING</h1>
            <div style={s2}>
                <Loader color="#26A65B" size="16px" />
            </div>
        </div>
    }
}
