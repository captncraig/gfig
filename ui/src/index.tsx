import * as React from 'react';
import * as ReactDom from 'react-dom';
import { Spinner } from './components/spinner'
import {Settings} from './components/settings'
import {SettingsMetadata} from './models'

interface AppState {
    loaded: boolean
    error: string
    metadata: SettingsMetadata
}

class App extends React.Component<{}, AppState>{
    constructor() {
        super();
        this.state = { loaded: false, error: "", metadata: {} as SettingsMetadata };
    }
    public render() {
        if (!this.state.loaded) {
            return <div className="container">
                <Spinner />
            </div>
        }
        return <div className="container"><Settings Metadata={this.state.metadata}/></div>
    }
    public async componentDidMount() {
        try {
            let resp = await fetch("/api/metadata");
            let data = await resp.json() as SettingsMetadata;
            this.setState({ loaded: true, metadata: data })
        } catch (e) {
            this.setState({ error: e })
        }
    }
}
ReactDom.render(
    <App />,
    document.getElementById('app'),
);


