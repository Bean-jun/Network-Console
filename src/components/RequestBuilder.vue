<script setup lang="ts">
import { ref, computed } from 'vue'
import axios from 'axios'
import type { RequestData, ResponseData, HttpMethod, KeyValuePair } from '../types'

const props = defineProps<{ modelValue: RequestData }>()
const emit = defineEmits<{
    'update:modelValue': [value: RequestData]
    'sent': [response: ResponseData]
    'loading': [value: boolean]
}>()

const httpMethods: HttpMethod[] = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS']
const activeSubTab = ref<'query' | 'headers' | 'body'>('query')
const bodyMode = ref<'form-data' | 'x-www-form-urlencoded' | 'raw'>('raw')
const rawType = ref<'json' | 'text' | 'xml'>('json')
const formData = ref<KeyValuePair[]>([{ id: '1', key: '', value: '', enabled: true }])

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
    return rawType.value === 'json' ? '{"key": "value"}' : 'Enter request body...'
})

const addPair = (arr: KeyValuePair[]) => arr.push({ id: Date.now().toString(), key: '', value: '', enabled: true })
const removePair = (arr: KeyValuePair[], id: string) => {
    const i = arr.findIndex(x => x.id === id)
    if (i > -1 && arr.length > 1) arr.splice(i, 1)
}

const setActiveTab = (tab: string) => {
    activeSubTab.value = tab as 'query' | 'headers' | 'body'
}

const getBodyContent = () => {
    if (bodyMode.value === 'raw') return request.value.body
    const data = formData.value.filter(x => x.enabled && x.key)
    if (bodyMode.value === 'x-www-form-urlencoded') {
        return data.map(x => `${encodeURIComponent(x.key)}=${encodeURIComponent(x.value)}`).join('&')
    }
    const fd = new FormData()
    data.forEach(x => fd.append(x.key, x.value))
    return fd
}

const sendRequest = async () => {
    emit('loading', true)
    try {
        const start = Date.now()
        const headers = request.value.headers.filter(h => h.enabled && h.key).reduce((a, h) => (a[h.key] = h.value, a), {} as Record<string, string>)

        if (bodyMode.value === 'x-www-form-urlencoded') {
            headers['Content-Type'] = 'application/x-www-form-urlencoded'
        }

        const res = await axios({
            method: request.value.method,
            url: computedUrl.value,
            headers,
            data: request.value.method !== 'GET' ? getBodyContent() : undefined,
            validateStatus: () => true
        })
        emit('sent', {
            status: res.status,
            statusText: res.statusText,
            headers: res.headers,
            data: res.data,
            time: Date.now() - start,
            size: JSON.stringify(res.data).length,
            contentType: res.headers['content-type'] || ''
        })
    } catch (e: any) {
        emit('sent', { status: 0, statusText: e.message, headers: {}, data: e.toString(), time: 0, size: 0, contentType: 'text/plain' })
    } finally { emit('loading', false) }
}

const formatJson = () => {
    try {
        const parsed = JSON.parse(request.value.body)
        request.value.body = JSON.stringify(parsed, null, 2)
    } catch { }
}
</script>

<template>
    <div class="builder">
        <!-- 请求行 -->
        <div class="request-line">
            <select v-model="request.method" class="method">
                <option v-for="m in httpMethods" :key="m" :value="m">{{ m }}</option>
            </select>
            <input v-model="request.url" class="url" placeholder="https://api.example.com/data">
            <button class="send" @click="sendRequest">Send</button>
        </div>

        <!-- Computed URL 预览 -->
        <div v-if="computedUrl" class="computed">
            <span class="label">Computed URL:</span>
            <span class="value">{{ computedUrl }}</span>
        </div>

        <!-- 子选项卡 -->
        <div class="sub-tabs">
            <div v-for="t in ['query', 'headers', 'body']" :key="t" class="sub-tab"
                :class="{ active: activeSubTab === t }" @click="setActiveTab(t)">
                {{ t === 'query' ? 'Query Params' : t === 'headers' ? 'Headers' : 'Body' }}
            </div>
        </div>

        <!-- Query Params -->
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

        <!-- Headers -->
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

        <!-- Body - 统一表格风格 -->
        <div v-if="activeSubTab === 'body'" class="body-wrap">
            <!-- Body 类型选择 -->
            <div class="body-modes">
                <label v-for="mode in ['form-data', 'x-www-form-urlencoded', 'raw']" :key="mode" class="mode-label">
                    <input type="radio" v-model="bodyMode" :value="mode">
                    {{ mode }}
                </label>
            </div>

            <!-- form-data / urlencoded 用表格编辑 -->
            <div v-if="bodyMode !== 'raw'" class="table-wrap">
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
                            <td><button class="del-btn" @click="removePair(formData, item.id)">✕</button></td>
                        </tr>
                    </tbody>
                </table>
                <button class="add-btn" @click="addPair(formData)">+ Add Field</button>
            </div>

            <!-- raw 类型用文本域 -->
            <div v-if="bodyMode === 'raw'" class="raw-wrap">
                <div class="raw-types">
                    <label v-for="t in ['json', 'text', 'xml']" :key="t" class="type-label">
                        <input type="radio" v-model="rawType" :value="t"> {{ t.toUpperCase() }}
                    </label>
                    <button v-if="rawType === 'json'" class="format-btn" @click="formatJson">Format</button>
                </div>
                <textarea v-model="request.body" class="body-textarea" :placeholder="bodyPlaceholder"></textarea>
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
</style>