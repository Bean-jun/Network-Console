export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD' | 'OPTIONS'

export interface KeyValuePair {
    id: string
    key: string
    value: string
    enabled: boolean
}

export interface RequestData {
    method: HttpMethod
    url: string
    queryParams: KeyValuePair[]
    headers: KeyValuePair[]
    body: string
    bodyType: 'json' | 'form' | 'text' | 'xml'
}

export interface ResponseData {
    status: number
    statusText: string
    headers: Record<string, string>
    data: any
    time: number
    size: number
    contentType: string
}