import * as React from 'react'
import { Setting } from '../models'

import * as _ from 'lodash'

interface SettingsProps {
    Settings: Setting[];
}

export class Settings extends React.Component<SettingsProps, {}>{

    public render() {
        var groups = _(this.props.Settings).groupBy<string>((s) => s.name.split('.')[0])
            .map((k: Setting[], v: string) => {
                return <h1 key={v}>{v}</h1>
            }).value();
        return <div>{groups}</div>
    }
}
