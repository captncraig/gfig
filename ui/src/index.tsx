import * as React from 'react';
import * as ReactDom from 'react-dom';
import { Spinner } from './components/spinner'
import {Settings} from './components/settings'
import {Filters} from './components/filters'

import {AppData, Tagset} from './models'

interface AppState {
    loaded: boolean
    error: string
    metadata: AppData
    filters: Tagset
}

class App extends React.Component<{}, AppState>{
    constructor() {
        super();
        this.state = { loaded: false, error: "", metadata: {} as AppData, filters: {} };
    }
    filtersChanged = (f: Tagset) => {
        this.setState({filters: f})
    }
    public render() {
        if (!this.state.loaded) {
            return <div className="container">
                <Spinner />
            </div>
        }
        return <div className="container">
            <Filters Tags={this.state.metadata.tags} OnFilterChange={this.filtersChanged} CurrentTags={this.state.metadata.currentTags} />
            <Settings Settings={this.state.metadata.settings} />
        </div>
    }
    public async componentDidMount() {
        try {
            let resp = await fetch("/api/metadata");
            let data = await resp.json() as AppData;
            this.setState({ loaded: true, metadata: data })
        } catch (e) {
            console.log("EEE", e)
            this.setState({ error: e })
        }
    }
}
ReactDom.render(
    <App />,
    document.getElementById('app'),
);


