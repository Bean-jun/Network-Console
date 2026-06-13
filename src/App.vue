<script setup lang="ts">
import { ref } from 'vue'
import RequestBuilder from './components/RequestBuilder.vue'
import ResponseViewer from './components/ResponseViewer.vue'
import type { RequestData, ResponseData } from './types'

const activeTab = ref<'request' | 'response'>('request')
const currentRequest = ref<RequestData>({
    method: 'GET',
    url: '',
    queryParams: [{ id: '1', key: '', value: '', enabled: true }],
    headers: [{ id: '1', key: '', value: '', enabled: true }],
    body: '',
    bodyType: 'json'
})
const currentResponse = ref<ResponseData | null>(null)
const isLoading = ref(false)

const handleRequestSent = (response: ResponseData) => {
    currentResponse.value = response
    activeTab.value = 'response'
}
</script>

<template>
    <div class="app-container">
        <header class="app-header">
            <h1>Network Console</h1>
            <span v-if="isLoading" class="loading">Sending...</span>
        </header>
        <main class="app-main">
            <div class="tabs">
                <div class="tab" :class="{ active: activeTab === 'request' }" @click="activeTab = 'request'">Request
                </div>
                <div class="tab" :class="{ active: activeTab === 'response' }" @click="activeTab = 'response'">Response
                </div>
            </div>
            <div class="tab-content">
                <RequestBuilder v-if="activeTab === 'request'" v-model="currentRequest" @sent="handleRequestSent"
                    @loading="isLoading = $event" />
                <ResponseViewer v-else-if="currentResponse" :response="currentResponse" />
                <div v-else class="empty">Send a request to see results</div>
            </div>
        </main>
    </div>
</template>

<style scoped>
.app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
}

.app-header {
    padding: 8px 16px;
    background: #252526;
    border-bottom: 1px solid #333;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.app-header h1 {
    font-size: 14px;
    font-weight: 500;
    color: #ccc;
}

.loading {
    color: #0078d4;
    animation: pulse 1.5s infinite;
}

@keyframes pulse {

    0%,
    100% {
        opacity: 1;
    }

    50% {
        opacity: 0.5;
    }
}

.app-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.tab-content {
    flex: 1;
    overflow: auto;
}

.empty {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #666;
}
</style>