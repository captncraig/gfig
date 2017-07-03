import * as React from 'react';
import {ScaleLoader as Loader} from 'halogen'

export default class Spinner extends React.Component<{},{}>{
    public render(){
        var s = {textAlign: "center"}
        var s2 = {margin: "0 auto", width:"60px"}
        return <div >
            <h1 style={s}>LOADING</h1>
            <div style={s2}>
            <Loader color="#26A65B" size="16px" margin="4px"/>
            </div>
        </div>
    }
}