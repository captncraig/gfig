import * as React from 'react'
import { Setting } from '../models'

import * as _ from 'lodash'

interface SettingsProps {
    Settings: Setting[];
}


export class Settings extends React.Component<SettingsProps, {}>{
    public render() {
        var mainStyle ={
            border: "1px solid #ddd",
            borderTop: "none",
        }
        var groupStyle = {
            boxSizing: 'border-box',
        };
        var headStyle = {
            padding: "5px 8px",
            fontWeight: 700 as 700,
            backgroundColor: "#eee",
            borderTop: "1px solid #ddd",
            margin: 0,
        }
        var groups = _(this.props.Settings).groupBy<string>((s) => s.name.split('.')[0])
            .map((k: Setting[], v: string) => {
                return <div style={groupStyle} key={v}><h4 style={headStyle}>{v}</h4>
                    {k.map((s) => <SingleSetting key={s.name} Setting={s} />)}
                </div>
            }).value();
        return <div style={mainStyle}>{groups}</div>
    }
}

interface SingleSettingProps {
    Setting: Setting;
}

class SingleSetting extends React.Component<SingleSettingProps, {}>{

    public render() {
        var s = this.props.Setting;
        var rowStyle = {
            paddingTop: 0,
            paddingBottom: 0,
            display: "flex",
            flexDirection: "row" as "row",
            flexWrap: "wrap" as "wrap",
            alignItems: "stretch"as "stretch",
            borderTop: "1px solid #ddd",
            cursor: "default",
        }
        return <div style={rowStyle}>
            <div><a><strong>{s.name}</strong></a></div>
        </div>
    }
}