import * as React from 'react'
import { TagMetadata, Tagset } from '../models'
import * as _ from 'lodash'

interface FiltersProps {
    Tags: TagMetadata[];
    CurrentTags: Tagset
    OnFilterChange(f: Tagset): void;
}

interface FilterState {
    SelectedTags: Tagset
    MaxSelected: string
}

export class Filters extends React.Component<FiltersProps, FilterState>{
    public render() {
        var active = true;
        var breadcrumbs = this.props.Tags.map((tag) => {
            var cname = ""
            if (!active) {
                cname = "unselected"
            } else if (tag.name == this.state.MaxSelected) {
                active = false;
            }
            return <li className={cname} key={tag.name} onClick={this.setMax.bind(this, tag.name)}><a><strong>{tag.name}</strong> {this.state.SelectedTags[tag.name]}</a></li>
        })
        return <div className='filters'>
            <ul className="breadcrumbs"> {breadcrumbs} </ul>
        </div>
    }
    setMax = (n: string) => {
        this.setState({MaxSelected: n})
        this.updateFilters();
    }

    updateFilters = () => {
        console.log(this,"@@@")
        var filter:  Tagset = {}
        for(var t of this.props.Tags){
            filter[t.name] = this.state.SelectedTags[t.name];
            if (t.name == this.state.MaxSelected){
                break
            } 
        }
        this.props.OnFilterChange(filter);
    }
    constructor(props: FiltersProps) {
        super(props)
        var name = ""
        var max = _(this.props.Tags).findLast(t => t.initialFilter);
        if (max) {
            name = max.name
        }
        this.state = { SelectedTags: props.CurrentTags, MaxSelected: name }
    }
    public componentDidMount(){
        this.updateFilters();
    }
}