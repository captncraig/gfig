export interface SettingsMetadata {
    appName: string
    currentTags: {[id: string]: string}
    tags: TagMetadata[] 
}
export interface TagMetadata {
    name: string
    allowNewValues: boolean
    initialFilter: boolean
}
