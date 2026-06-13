<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import axios from 'axios'
import type { RequestData, HttpMethod, KeyValuePair } from '../types'

const props = defineProps<{ modelValue: RequestData }>()
const emit = defineEmits(['update:modelValue', 'sent', 'loading'])

interface FormField extends KeyValuePair {
    type: 'text' | 'file'
    file?: File | null
}

const httpMethods: HttpMethod[] = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS']
const activeSubTab = ref<'query' | 'headers' | 'body'>('query')
const bodyMode = ref<'form-data' | 'x-www-form-urlencoded' | 'raw'>('raw')
const rawType = ref<'json' | 'text' | 'xml'>('json')
const formData = ref<FormField[]>([{ id: '1', key: '', value: '', enabled: true, type: 'text' }])
const jsonError = ref<string>('')
const xmlError = ref<string>('')
const useProxy = ref(true)
const proxyUrl = window.location.origin + '/proxy'

const request = computed({
    get: () => props.modelValue,
    set: v => emit('update:modelValue', v)
})

const computedUrl = computed(() => {
    const params = request.value.queryParams.filter(p => p.enabled && p.key)
    if (!params.length) return request.value.url
    const sep = request.value.url.includes('?') ? '&' : '?'
    const qs = params.map(p => `${encodeURIComponent(p.key)}=${encodeURIComponent(p.value)}`).join('&')
    return request.value.url + sep + qs
})

const bodyPlaceholder = computed(() => {
    return rawType.value === 'json' ? '{"key": "value"}' : rawType.value === 'xml' ? '<root></root>' : 'Enter request body...'
})

watch(() => [rawType.value, request.value.body], () => {
    jsonError.value = ''
    xmlError.value = ''
    if (rawType.value === 'json' && request.value.body.trim()) {
        try { JSON.parse(request.value.body) } catch (e: any) { jsonError.value = e.message }
    }
    if (rawType.value === 'xml' && request.value.body.trim()) {
        try {
            const parser = new DOMParser()
            const doc = parser.parseFromString(request.value.body, 'text/xml')
            const parseError = doc.querySelector('parsererror')
            if (parseError) xmlError.value = parseError.textContent || 'Invalid XML'
        } catch (e: any) { xmlError.value = e.message }
    }
}, { immediate: true })

const addFormField = () => formData.value.push({ id: Date.now().toString(), key: '', value: '', enabled: true, type: 'text', file: null })
const removeFormField = (id: string) => { const i = formData.value.findIndex(x => x.id === id); if (i > -1 && formData.value.length > 1) formData.value.splice(i, 1) }
const addPair = (arr: KeyValuePair[]) => arr.push({ id: Date.now().toString(), key: '', value: '', enabled: true })
const removePair = (arr: KeyValuePair[], id: string) => { const i = arr.findIndex(x => x.id === id); if (i > -1 && arr.length > 1) arr.splice(i, 1) }
const handleFileChange = (field: FormField, event: Event) => {
    const input = event.target as HTMLInputElement
    if (input.files && input.files[0]) { field.file = input.files[0]; field.value = input.files[0].name }
}
const setActiveTab = (tab: string) => { activeSubTab.value = tab as 'query' | 'headers' | 'body' }

const getBodyContent = () => {
    if (bodyMode.value === 'raw') return request.value.body
    const data = formData.value.filter(x => x.enabled && x.key)
    if (bodyMode.value === 'x-www-form-urlencoded') return data.map(x => `${encodeURIComponent(x.key)}=${encodeURIComponent(x.value)}`).join('&')
    const fd = new FormData()
    data.forEach(x => { if (x.type === 'file' && x.file) fd.append(x.key, x.file); else fd.append(x.key, x.value) })
    return fd
}

const sendRequest = async () => {
    emit('loading', true)
    try {
        const start = Date.now()
        const headers = request.value.headers.filter(h => h.enabled && h.key).reduce((a, h) => (a[h.key] = h.value, a), {} as Record<string, string>)
        if (bodyMode.value === 'x-www-form-urlencoded') headers['Content-Type'] = 'application/x-www-form-urlencoded'

        let finalUrl = computedUrl.value
        if (useProxy.value) finalUrl = `${proxyUrl}?url=${encodeURIComponent(computedUrl.value)}`

        const res = await axios({
            method: request.value.method,
            url: finalUrl,
            headers,
            data: request.value.method !== 'GET' ? getBodyContent() : undefined,
            validateStatus: () => true
        })
        emit('sent', { status: res.status, statusText: res.statusText, headers: res.headers, data: res.data, time: Date.now() - start, size: JSON.stringify(res.data).length, contentType: res.headers['content-type'] || '' })
    } catch (e: any) {
        emit('sent', { status: 0, statusText: e.message, headers: {}, data: e.toString(), time: 0, size: 0, contentType: 'text/plain' })
    } finally { emit('loading', false) }
}

const formatJson = () => {
    try { const parsed = JSON.parse(request.value.body); request.value.body = JSON.stringify(parsed, null, 2); jsonError.value = '' }
    catch (e: any) { jsonError.value = e.message }
}
</script>

<template>
    <div class="builder">
        <div class="request-line">
            <select v-model="request.method" class="method">
                <option v-for="m in httpMethods" :key="m" :value="m">{{ m }}</option>
            </select>
            <input v-model="request.url" class="url" placeholder="https://api.example.com/data">
            <label class="proxy-switch"><input type="checkbox" v-model="useProxy"><span
                    class="proxy-text">代理</span></label>
            <button class="send" @click="sendRequest">Send</button>
        </div>

        <div v-if="computedUrl" class="computed">
            <span class="label">{{ useProxy ? '代理 URL:' : 'Computed URL:' }}</span>
            <span class="value">{{ useProxy ? proxyUrl + '?url=...' : computedUrl }}</span>
        </div>

        <div class="sub-tabs">
            <div v-for="t in ['query', 'headers', 'body']" :key="t" class="sub-tab"
                :class="{ active: activeSubTab === t }" @click="setActiveTab(t)">
                {{ t === 'query' ? 'Query Params' : t === 'headers' ? 'Headers' : 'Body' }}
            </div>
        </div>

        <div v-if="activeSubTab === 'query'" class="table-wrap">
            <table class="kv-table">
                <thead>
                    <tr>
                        <th style="width: 40px"></th>
                        <th>Key</th>
                        <th>Value</th>
                        <th style="width: 50px"></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="item in request.queryParams" :key="item.id">
                        <td><input type="checkbox" v-model="item.enabled"></td>
                        <td><input type="text" v-model="item.key" placeholder="param_name"></td>
                        <td><input type="text" v-model="item.value" placeholder="value"></td>
                        <td><button class="del-btn" @click="removePair(request.queryParams, item.id)">✕</button></td>
                    </tr>
                </tbody>
            </table>
            <button class="add-btn" @click="addPair(request.queryParams)">+ Add Parameter</button>
        </div>

        <div v-if="activeSubTab === 'headers'" class="table-wrap">
            <table class="kv-table">
                <thead>
                    <tr>
                        <th style="width: 40px"></th>
                        <th>Key</th>
                        <th>Value</th>
                        <th style="width: 50px"></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="item in request.headers" :key="item.id">
                        <td><input type="checkbox" v-model="item.enabled"></td>
                        <td><input type="text" v-model="item.key" placeholder="Content-Type"></td>
                        <td><input type="text" v-model="item.value" placeholder="application/json"></td>
                        <td><button class="del-btn" @click="removePair(request.headers, item.id)">✕</button></td>
                    </tr>
                </tbody>
            </table>
            <button class="add-btn" @click="addPair(request.headers)">+ Add Header</button>
        </div>

        <div v-if="activeSubTab === 'body'" class="body-wrap">
            <div class="body-modes">
                <label v-for="mode in ['form-data', 'x-www-form-urlencoded', 'raw']" :key="mode" class="mode-label">
                    <input type="radio" v-model="bodyMode" :value="mode">{{ mode }}
                </label>
            </div>

            <div v-if="bodyMode === 'form-data'" class="table-wrap">
                <table class="kv-table form-table">
                    <thead>
                        <tr>
                            <th style="width: 40px"></th>
                            <th style="width: 80px">Type</th>
                            <th>Key</th>
                            <th>Value</th>
                            <th style="width: 50px"></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in formData" :key="item.id">
                            <td><input type="checkbox" v-model="item.enabled"></td>
                            <td><select v-model="item.type" class="type-select">
                                    <option value="text">Text</option>
                                    <option value="file">File</option>
                                </select></td>
                            <td><input type="text" v-model="item.key" placeholder="field_name"></td>
                            <td>
                                <input v-if="item.type === 'text'" type="text" v-model="item.value" placeholder="value">
                                <div v-else class="file-input-wrap">
                                    <input type="file" :id="'file-' + item.id" class="file-input"
                                        @change="handleFileChange(item, $event)">
                                    <label :for="'file-' + item.id" class="file-label">{{ item.value || 'Choose file...'
                                        }}</label>
                                </div>
                            </td>
                            <td><button class="del-btn" @click="removeFormField(item.id)">✕</button></td>
                        </tr>
                    </tbody>
                </table>
                <button class="add-btn" @click="addFormField">+ Add Field</button>
            </div>

            <div v-if="bodyMode === 'x-www-form-urlencoded'" class="table-wrap">
                <table class="kv-table">
                    <thead>
                        <tr>
                            <th style="width: 40px"></th>
                            <th>Key</th>
                            <th>Value</th>
                            <th style="width: 50px"></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in formData" :key="item.id">
                            <td><input type="checkbox" v-model="item.enabled"></td>
                            <td><input type="text" v-model="item.key" placeholder="field_name"></td>
                            <td><input type="text" v-model="item.value" placeholder="value"></td>
                            <td><button class="del-btn" @click="removeFormField(item.id)">✕</button></td>
                        </tr>
                    </tbody>
                </table>
                <button class="add-btn" @click="addFormField">+ Add Field</button>
            </div>

            <div v-if="bodyMode === 'raw'" class="raw-wrap">
                <div class="raw-types">
                    <label v-for="t in ['json', 'text', 'xml']" :key="t" class="type-label">
                        <input type="radio" v-model="rawType" :value="t"> {{ t.toUpperCase() }}
                    </label>
                    <button v-if="rawType === 'json'" class="format-btn" @click="formatJson">Format</button>
                </div>
                <div v-if="jsonError" class="error-bar"><span class="error-icon">⚠</span><span class="error-text">JSON
                        Error: {{ jsonError }}</span></div>
                <div v-if="xmlError" class="error-bar"><span class="error-icon">⚠</span><span class="error-text">XML
                        Error: {{ xmlError }}</span></div>
                <textarea v-model="request.body" class="body-textarea" :class="{ 'has-error': jsonError || xmlError }"
                    :placeholder="bodyPlaceholder"></textarea>
            </div>
        </div>
    </div>
</template>

<style scoped>
.builder {
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    height: 100%;
}

.request-line {
    display: flex;
    gap: 8px;
}

.method {
    width: 100px;
    padding: 6px 8px;
    font-weight: 600;
}

.url {
    flex: 1;
    padding: 6px 10px;
}

.proxy-switch {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 0 12px;
    background: #2d2d2d;
    border: 1px solid #555;
    border-radius: 3px;
    cursor: pointer;
    color: #ccc;
    font-size: 12px;
}

.proxy-switch input[type="checkbox"]:checked+.proxy-text {
    color: #4ec9b0;
}

.send {
    padding: 6px 20px;
    background: #0e639c;
    color: white;
    font-weight: 500;
}

.send:hover {
    background: #1177bb;
}

.computed {
    padding: 8px 12px;
    background: #252526;
    border: 1px solid #333;
    font-family: Consolas, monospace;
    font-size: 11px;
    word-break: break-all;
}

.computed .label {
    color: #888;
    margin-right: 8px;
}

.computed .value {
    color: #4ec9b0;
}

.sub-tabs {
    display: flex;
    border-bottom: 1px solid #333;
}

.sub-tab {
    padding: 8px 16px;
    cursor: pointer;
    color: #999;
    border-bottom: 2px solid transparent;
    font-size: 12px;
}

.sub-tab:hover {
    color: #ccc;
    background: #2d2d2d;
}

.sub-tab.active {
    color: #fff;
    border-bottom-color: #0078d4;
}

.table-wrap {
    flex: 1;
    overflow: auto;
    display: flex;
    flex-direction: column;
}

.kv-table {
    width: 100%;
    border-collapse: collapse;
}

.kv-table th {
    background: #2d2d2d;
    padding: 8px;
    text-align: left;
    color: #999;
    font-weight: 500;
    border-bottom: 1px solid #333;
    font-size: 12px;
}

.kv-table td {
    padding: 4px;
    border-bottom: 1px solid #333;
}

.kv-table input[type="text"] {
    width: 100%;
    border: none;
    background: transparent;
    padding: 6px 4px;
}

.kv-table input[type="text"]:focus {
    background: #3c3c3c;
}

.type-select {
    width: 100%;
    padding: 4px;
    border: none;
    background: #3c3c3c;
    font-size: 11px;
}

.file-input-wrap {
    position: relative;
    width: 100%;
}

.file-input {
    position: absolute;
    opacity: 0;
    width: 100%;
    height: 100%;
    cursor: pointer;
}

.file-label {
    display: block;
    padding: 6px 8px;
    background: #3c3c3c;
    border: 1px solid #555;
    border-radius: 3px;
    cursor: pointer;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 11px;
    color: #ccc;
}

.file-label:hover {
    background: #464646;
}

.del-btn {
    background: transparent;
    color: #f48771;
    padding: 4px 8px;
    font-size: 14px;
}

.del-btn:hover {
    background: #3c2020;
}

.add-btn {
    align-self: flex-start;
    margin-top: 8px;
    padding: 6px 12px;
    background: #3c3c3c;
    color: #ccc;
    border: 1px solid #555;
}

.add-btn:hover {
    background: #464646;
}

.body-wrap {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;
    overflow: auto;
}

.body-modes,
.raw-types {
    display: flex;
    gap: 16px;
    align-items: center;
    padding: 4px 0;
}

.mode-label,
.type-label {
    display: flex;
    align-items: center;
    gap: 4px;
    cursor: pointer;
    color: #ccc;
    font-size: 12px;
}

.format-btn {
    margin-left: auto;
    padding: 4px 12px;
    background: #3c3c3c;
    color: #ccc;
    border: 1px solid #555;
    font-size: 11px;
}

.format-btn:hover {
    background: #464646;
}

.error-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background: #3c2020;
    border: 1px solid #6b2020;
    border-radius: 3px;
}

.error-icon {
    color: #f48771;
    font-size: 14px;
}

.error-text {
    color: #f48771;
    font-size: 12px;
    font-family: Consolas, monospace;
}

.raw-wrap {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.body-textarea {
    flex: 1;
    min-height: 200px;
    padding: 12px;
    font-family: Consolas, monospace;
    font-size: 13px;
    line-height: 1.6;
}

.body-textarea.has-error {
    border-color: #f48771;
}
</style>