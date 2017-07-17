export interface SettingsMetadata {
    appName: string;
    currentTags:Tagset;
    tags: TagMetadata[] ;
}

export interface TagMetadata {
    name: string;
    allowNewValues: boolean;
    initialFilter: boolean;
}

export type Tagset =  {[id: string]: string};