<script setup>
import { NButton, NCard, NConfigProvider, NDivider, NForm, NFormItem, NInputNumber, NSpace, NTag, darkTheme } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { herocfg as hcfg } from '../cfg'
import { exportLineup, lineupList } from '../api'

const herocfg = JSON.parse(hcfg)
const router = useRouter()

const records = ref([])
const nextid = ref(0)
const total = ref(0)
const loading = ref(false)
const hasMore = ref(true)
const minLevel = ref(null)

const formatHero = (id, star, level, name) => {
  const display = name || herocfg[id]?.name || id || '-'
  return `${display} ${star}红 ${level}级`
}

const formatTime = (timestamp) => {
  const date = new Date(timestamp * 1000)
  const y = date.getFullYear()
  const m = `${date.getMonth() + 1}`.padStart(2, '0')
  const d = `${date.getDate()}`.padStart(2, '0')
  const hh = `${date.getHours()}`.padStart(2, '0')
  const mm = `${date.getMinutes()}`.padStart(2, '0')
  const ss = `${date.getSeconds()}`.padStart(2, '0')
  return `${y}-${m}-${d} ${hh}:${mm}:${ss}`
}

const fetchRecords = (clear = false) => {
  if (loading.value) return
  if (clear) {
    nextid.value = 0
    hasMore.value = true
    records.value = []
  }
  loading.value = true
  lineupList({
    nextid: nextid.value,
    minlevel: minLevel.value ?? ''
  }).then((resp) => {
    const data = resp.data.data
    total.value = data.total
    if (data.list.length > 0) {
      nextid.value = data.list[data.list.length - 1].id
      records.value = [...records.value, ...data.list]
    } else {
      hasMore.value = false
    }
  }).finally(() => {
    loading.value = false
  })
}

const reset = () => {
  minLevel.value = null
  fetchRecords(true)
}

const downloadCsv = () => {
  exportLineup({
    minlevel: minLevel.value ?? ''
  }).then((resp) => {
    const blob = new Blob([resp.data], { type: 'text/csv;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'battle_lineup.csv'
    a.click()
    URL.revokeObjectURL(url)
  })
}

onMounted(() => {
  fetchRecords(true)
})
</script>

<template>
  <n-config-provider :theme="darkTheme">
    <div style="text-align: center; padding: 16px;">
      <n-space justify="center" style="margin-bottom: 8px;">
        <n-button @click="() => router.push('/')">返回战报列表</n-button>
        <n-button @click="() => router.push('/team')">队伍查询</n-button>
        <n-button @click="() => router.push('/lineup')">阵容统计</n-button>
      </n-space>
      <n-divider>战报信息</n-divider>
      <div style="color: #9ca3af; margin-bottom: 8px;">过滤掉最低等级低于阈值的战报，并导出整理后的CSV。</div>
      <n-form inline :label-width="80" style="justify-content: center;">
        <n-form-item label="最低等级">
          <n-input-number v-model:value="minLevel" :min="1" style="width: 160px" placeholder="不限制" />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" @click="() => fetchRecords(true)">筛选</n-button>
        </n-form-item>
        <n-form-item>
          <n-button @click="reset">重置</n-button>
        </n-form-item>
        <n-form-item>
          <n-button type="success" @click="downloadCsv">下载CSV</n-button>
        </n-form-item>
      </n-form>
      <div style="margin-top: 12px; color: #9ca3af;">记录总数：{{ total }}</div>
      <div class="record-list">
        <n-card v-for="item in records" :key="item.id" style="margin: 12px auto; max-width: 960px; text-align: left;">
          <template #header>
            <n-space justify="space-between" align="center" style="width: 100%;">
              <div>
                <div style="font-weight: 600;">{{ item.player_name }}</div>
                <div style="font-size: 12px; color: #94a3b8;">{{ item.union_name || '未记录同盟' }}</div>
              </div>
              <n-tag type="info">{{ item.player_role === 'defend' ? '防守' : '进攻' }}</n-tag>
            </n-space>
          </template>
          <div style="margin-bottom: 8px; font-size: 12px; color: #94a3b8;">
            战报ID：{{ item.battle_id }} ｜ 记录时间：{{ formatTime(item.record_time) }}
          </div>
          <n-space vertical>
            <div class="lineup-row">
              <span class="lineup-label">大营</span>
              <span>{{ formatHero(item.hero1_id, item.hero1_star, item.hero1_level, item.hero1_name) }}</span>
            </div>
            <div class="lineup-row">
              <span class="lineup-label">中军</span>
              <span>{{ formatHero(item.hero2_id, item.hero2_star, item.hero2_level, item.hero2_name) }}</span>
            </div>
            <div class="lineup-row">
              <span class="lineup-label">前锋</span>
              <span>{{ formatHero(item.hero3_id, item.hero3_star, item.hero3_level, item.hero3_name) }}</span>
            </div>
            <div style="font-size: 12px; color: #94a3b8;">阵容标识：{{ item.lineup_key }}</div>
          </n-space>
        </n-card>
      </div>
      <n-space justify="center" style="margin-top: 12px;">
        <n-button @click="fetchRecords()" :disabled="!hasMore || loading">{{ hasMore ? '加载更多' : '没有更多了' }}</n-button>
      </n-space>
    </div>
  </n-config-provider>
</template>

<style scoped>
.lineup-row {
  display: flex;
  gap: 12px;
  align-items: center;
  font-size: 14px;
}

.lineup-label {
  display: inline-block;
  width: 40px;
  color: #cbd5e1;
}
</style>
