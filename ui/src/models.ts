export interface AppData {
    appName: string;
    currentTags:Tagset;
    tags: TagMetadata[];
    settings: Setting[];
}
export interface Setting {
    name: string;
}
export interface TagMetadata {
    name: string;
    allowNewValues: boolean;
    initialFilter: boolean;
}

export type Tagset =  {[id: string]: string};