import * as React from 'react';
import * as ReactDom from 'react-dom';
import Spinner from './components/spinner'

interface AppState {
    loaded: boolean
    error: string
}
class App extends React.Component<{}, AppState>{
    public render() {
        
        return <div className="container">
            <Spinner/>
        </div>
    }
    public componentDidMount() {
        console.log("CDM")
    }
}
ReactDom.render(
    <App />,
    document.getElementById('app'),
);


