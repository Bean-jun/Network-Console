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

const request = computed({ get: () => props.modelValue, set: v => emit('update:modelValue', v) })

const computedUrl = computed(() => {
    const params = request.value.queryParams.filter(p => p.enabled && p.key)
    if (!params.length) return request.value.url
    const sep = request.value.url.includes('?') ? '&' : '?'
    const qs = params.map(p => `${encodeURIComponent(p.key)}=${encodeURIComponent(p.value)}`).join('&')
    return request.value.url + sep + qs
})

const addPair = (arr: KeyValuePair[]) => arr.push({ id: Date.now().toString(), key: '', value: '', enabled: true })
const removePair = (arr: KeyValuePair[], id: string) => {
    const i = arr.findIndex(x => x.id === id)
    if (i > -1 && arr.length > 1) arr.splice(i, 1)
}

const sendRequest = async () => {
    emit('loading', true)
    try {
        const start = Date.now()
        const headers = request.value.headers.filter(h => h.enabled && h.key).reduce((a, h) => (a[h.key] = h.value, a), {} as Record<string, string>)
        const res = await axios({ method: request.value.method, url: computedUrl.value, headers, data: request.value.method !== 'GET' ? request.value.body : undefined, validateStatus: () => true })
        emit('sent', { status: res.status, statusText: res.statusText, headers: res.headers, data: res.data, time: Date.now() - start, size: JSON.stringify(res.data).length, contentType: res.headers['content-type'] || '' })
    } catch (e: any) {
        emit('sent', { status: 0, statusText: e.message, headers: {}, data: e.toString(), time: 0, size: 0, contentType: 'text/plain' })
    } finally { emit('loading', false) }
}
</script>

<template>
    <div class="builder">
        <div class="request-line">
            <select v-model="request.method" class="method">
                <option v-for="m in httpMethods" :key="m" :value="m">{{ m }}</option>
            </select>
            <input v-model="request.url" class="url" placeholder="https://api.example.com">
            <button class="send" @click="sendRequest">Send</button>
        </div>
        <div v-if="computedUrl" class="computed">Computed: {{ computedUrl }}</div>
        <div class="sub-tabs">
            <div v-for="t in ['query', 'headers', 'body']" :key="t" class="sub-tab"
                :class="{ active: activeSubTab === t }" @click="activeSubTab = t as any">{{ t === 'query' ? 'Query
                Params' : t === 'headers' ? 'Headers' : 'Body' }}</div>
        </div>
        <div v-if="activeSubTab === 'query' || activeSubTab === 'headers'" class="table-wrap">
            <table>
                <thead>
                    <tr>
                        <th></th>
                        <th>Key</th>
                        <th>Value</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="item in activeSubTab === 'query' ? request.queryParams : request.headers" :key="item.id">
                        <td><input type="checkbox" v-model="item.enabled"></td>
                        <td><input type="text" v-model="item.key"></td>
                        <td><input type="text" v-model="item.value"></td>
                        <td><button class="del"
                                @click="removePair(activeSubTab === 'query' ? request.queryParams : request.headers, item.id)">✕</button>
                        </td>
                    </tr>
                </tbody>
            </table>
            <button class="add" @click="addPair(activeSubTab === 'query' ? request.queryParams : request.headers)">+
                Add</button>
        </div>
        <div v-if="activeSubTab === 'body'" class="body-wrap">
            <div class="body-types"><label v-for="t in ['json', 'form', 'text', 'xml']" :key="t"><input type="radio"
                        v-model="request.bodyType" :value="t"> {{ t.toUpperCase() }}</label></div>
            <textarea v-model="request.body" class="body-input" placeholder='{"key": "value"}'></textarea>
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
    color: #4ec9b0;
    word-break: break-all;
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
}

table {
    width: 100%;
    border-collapse: collapse;
}

th {
    background: #2d2d2d;
    padding: 8px;
    text-align: left;
    color: #999;
    font-weight: 500;
    border-bottom: 1px solid #333;
}

td {
    padding: 4px;
    border-bottom: 1px solid #333;
}

td input[type="text"] {
    width: 100%;
    border: none;
    background: transparent;
}

.del {
    background: transparent;
    color: #f48771;
    padding: 4px 8px;
}

.add {
    margin-top: 8px;
    padding: 6px 12px;
    background: #3c3c3c;
    color: #ccc;
    border: 1px solid #555;
}

.body-wrap {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.body-types {
    display: flex;
    gap: 16px;
}

.body-input {
    flex: 1;
    min-height: 200px;
    padding: 10px;
    font-family: Consolas, monospace;
    font-size: 13px;
}
</style>