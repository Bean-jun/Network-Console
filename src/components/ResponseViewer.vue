<script setup lang="ts">
import { ref, computed } from 'vue'
import type { ResponseData } from '../types'

const props = defineProps<{ response: ResponseData }>()
const activeTab = ref<'preview' | 'body' | 'headers' | 'cookies'>('preview')

const formattedBody = computed(() => {
    try { return typeof props.response.data === 'object' ? JSON.stringify(props.response.data, null, 2) : String(props.response.data) }
    catch { return String(props.response.data) }
})

const headersList = computed(() => Object.entries(props.response.headers).map(([key, value]) => ({ key, value })))
</script>

<template>
    <div class="viewer">
        <div class="status-bar">
            <span class="code" :class="'s-' + response.status">{{ response.status }}</span>
            <span class="text">{{ response.statusText }}</span>
            <span>|</span>
            <span>{{ response.time }}ms</span>
            <span>|</span>
            <span>{{ (response.size / 1024).toFixed(2) }} KB</span>
        </div>
        <div class="tabs">
            <div v-for="t in ['preview', 'body', 'headers', 'cookies']" :key="t" class="tab"
                :class="{ active: activeTab === t }" @click="activeTab = t as any">{{ t }}</div>
        </div>
        <div class="content">
            <pre v-if="activeTab === 'preview' || activeTab === 'body'" class="body">{{ formattedBody }}</pre>
            <table v-if="activeTab === 'headers'">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="h in headersList" :key="h.key">
                        <td>{{ h.key }}</td>
                        <td>{{ h.value }}</td>
                    </tr>
                </tbody>
            </table>
            <div v-if="activeTab === 'cookies'" class="empty">No cookies</div>
        </div>
    </div>
</template>

<style scoped>
.viewer {
    display: flex;
    flex-direction: column;
    height: 100%;
}

.status-bar {
    padding: 10px 12px;
    background: #252526;
    border-bottom: 1px solid #333;
    display: flex;
    gap: 12px;
    align-items: center;
}

.code {
    padding: 2px 8px;
    border-radius: 3px;
    font-weight: 600;
}

.s-200,
.s-201 {
    background: #0b5c0b;
    color: #86e086;
}

.s-301,
.s-302 {
    background: #5c4a0b;
    color: #e0d086;
}

.s-400,
.s-404 {
    background: #5c2a0b;
    color: #e0a086;
}

.s-500,
.s-0 {
    background: #5c0b0b;
    color: #e08686;
}

.tabs {
    display: flex;
    background: #252526;
    border-bottom: 1px solid #333;
}

.tab {
    padding: 8px 16px;
    cursor: pointer;
    color: #999;
    border-bottom: 2px solid transparent;
}

.tab:hover {
    color: #ccc;
    background: #2d2d2d;
}

.tab.active {
    color: #fff;
    border-bottom-color: #0078d4;
}

.content {
    flex: 1;
    overflow: auto;
    padding: 12px;
}

.body {
    margin: 0;
    padding: 12px;
    background: #1e1e1e;
    border: 1px solid #333;
    font-family: Consolas, monospace;
    font-size: 12px;
    line-height: 1.6;
    white-space: pre-wrap;
    word-break: break-all;
    min-height: 100%;
}

table {
    width: 100%;
    border: 1px solid #333;
    border-collapse: collapse;
}

th {
    background: #2d2d2d;
    padding: 8px 12px;
    text-align: left;
    color: #999;
    border-bottom: 1px solid #333;
}

td {
    padding: 8px 12px;
    border-bottom: 1px solid #333;
    font-family: Consolas, monospace;
}

.empty {
    padding: 40px;
    text-align: center;
    color: #666;
}
</style>