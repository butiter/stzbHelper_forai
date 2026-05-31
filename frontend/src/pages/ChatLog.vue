<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NCard, NDataTable, NForm, NFormItem, NInput, NInputNumber, NPagination, NSpace, NTabPane, NTabs, NTag } from 'naive-ui'
import { GetChatMessages, GetWorldAnnouncements } from '../../wailsjs/go/main/App'

const pageSize = 50
const chatPage = ref(1)
const chatTotal = ref(0)
const chatLoading = ref(false)
const channel = ref(null)
const keyword = ref('')
const chatList = ref([])

const annPage = ref(1)
const annTotal = ref(0)
const annLoading = ref(false)
const announcementList = ref([])

function parseResponse(text) {
    const resp = JSON.parse(text)
    return resp.data || { list: [], total: 0 }
}

function formatTime(value) {
    if (!value) return ''
    const date = new Date(value * 1000)
    return date.toLocaleString('zh-CN', { hour12: false })
}

function loadChat(reset = false) {
    if (reset) chatPage.value = 1
    chatLoading.value = true
    GetChatMessages(chatPage.value, pageSize, channel.value || 0, keyword.value).then(v => {
        const data = parseResponse(v)
        chatList.value = data.list || []
        chatTotal.value = data.total || 0
    }).finally(() => {
        chatLoading.value = false
    })
}

function loadAnnouncements(reset = false) {
    if (reset) annPage.value = 1
    annLoading.value = true
    GetWorldAnnouncements(annPage.value, pageSize).then(v => {
        const data = parseResponse(v)
        announcementList.value = data.list || []
        annTotal.value = data.total || 0
    }).finally(() => {
        annLoading.value = false
    })
}

const chatColumns = [
    { title: '时间', key: 'time', width: 170, render: row => formatTime(row.time) },
    { title: '频道', key: 'channel', width: 80, render: row => h(NTag, { size: 'small' }, { default: () => row.channel }) },
    { title: '玩家', key: 'player', width: 150 },
    { title: '同盟', key: 'union_name', width: 150 },
    { title: '内容', key: 'content' },
    { title: '语音', key: 'voice', width: 90, render: row => row.voice && row.voice !== 'null' ? '有' : '' }
]

const announcementColumns = [
    { title: '时间', key: 'time', width: 170, render: row => formatTime(row.time) },
    { title: '玩家', key: 'player' },
    { title: '事件', key: 'event_type', width: 120, render: () => '首次占领' },
    { title: '土地等级', key: 'land_level', width: 120 },
    { title: '原始数据', key: 'raw' }
]

onMounted(() => {
    loadChat()
    loadAnnouncements()
})
</script>

<template>
    <div class="page">
        <h2 class="page-title">聊天公告</h2>
        <NTabs type="line" animated>
            <NTabPane name="chat" tab="聊天记录">
                <NCard>
                    <NForm inline label-placement="left">
                        <NFormItem label="频道">
                            <NInputNumber v-model:value="channel" clearable :show-button="false" />
                        </NFormItem>
                        <NFormItem label="关键词">
                            <NInput v-model:value="keyword" clearable placeholder="玩家 / 同盟 / 内容" />
                        </NFormItem>
                        <NFormItem>
                            <NButton type="primary" @click="loadChat(true)">查询</NButton>
                        </NFormItem>
                    </NForm>
                    <NDataTable :columns="chatColumns" :data="chatList" :loading="chatLoading" :bordered="false" />
                    <NSpace justify="end" class="pager">
                        <NPagination v-model:page="chatPage" :page-size="pageSize" :item-count="chatTotal" @update:page="loadChat" />
                    </NSpace>
                </NCard>
            </NTabPane>
            <NTabPane name="announcement" tab="首占公告">
                <NCard>
                    <NSpace justify="end" class="toolbar">
                        <NButton @click="loadAnnouncements(true)">刷新</NButton>
                    </NSpace>
                    <NDataTable :columns="announcementColumns" :data="announcementList" :loading="annLoading" :bordered="false" />
                    <NSpace justify="end" class="pager">
                        <NPagination v-model:page="annPage" :page-size="pageSize" :item-count="annTotal" @update:page="loadAnnouncements" />
                    </NSpace>
                </NCard>
            </NTabPane>
        </NTabs>
    </div>
</template>

<style scoped>
.page {
    padding: 16px;
}

.page-title {
    margin: 0 0 16px;
    font-size: 20px;
    font-weight: 600;
}

.toolbar {
    margin-bottom: 12px;
}

.pager {
    margin-top: 16px;
}
</style>
