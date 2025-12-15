<script setup>

import { NButton, NCard, NConfigProvider, NDivider, NForm, NFormItem, NInput, NSelect, NSpace, NTag, darkTheme, NInputNumber } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { herocfg as hcfg } from '../cfg'
import { exportLineup, lineupList } from '../api'

const herocfg = JSON.parse(hcfg)
const router = useRouter()

const lineups = ref([])
const nextid = ref(0)
const total = ref(0)
const loading = ref(false)
const hasMore = ref(true)

const playername = ref('')
const unionname = ref('')
const lineupKey = ref('')
const role = ref(null)
const minLevel = ref(null)

const roleOptions = [
  { label: '全部', value: null },
  { label: '进攻', value: 'attack' },
  { label: '防守', value: 'defend' }
]

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

const fetchLineups = (clear = false) => {
  if (loading.value) return
  if (clear) {
    nextid.value = 0
    hasMore.value = true
    lineups.value = []
  }
  loading.value = true
  lineupList({
    nextid: nextid.value,
    playername: playername.value,
    unionname: unionname.value,
    lineup: lineupKey.value,
    role: role.value ?? '',
    minlevel: minLevel.value ?? ''
  }).then((resp) => {
    const data = resp.data.data
    total.value = data.total
    if (data.list.length > 0) {
      nextid.value = data.list[data.list.length - 1].id
      lineups.value = [...lineups.value, ...data.list]
    } else {
      hasMore.value = false
    }
  }).finally(() => {
    loading.value = false
  })
}

const reset = () => {
  playername.value = ''
  unionname.value = ''
  lineupKey.value = ''
  role.value = null
  minLevel.value = null
  fetchLineups(true)
}

const downloadCsv = () => {
  exportLineup({
    playername: playername.value,
    unionname: unionname.value,
    lineup: lineupKey.value,
    role: role.value ?? '',
    minlevel: minLevel.value ?? ''
  }).then((resp) => {
    const blob = new Blob([resp.data], { type: 'text/csv;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'lineup.csv'
    a.click()
    URL.revokeObjectURL(url)
  })
}

onMounted(() => {
  fetchLineups(true)
})
</script>

<template>
  <n-config-provider :theme="darkTheme">
    <div style="text-align: center; padding: 16px;">
      <n-space justify="center" style="margin-bottom: 8px;">
        <n-button @click="() => router.push('/')">返回战报列表</n-button>
        <n-button @click="() => router.push('/team')">队伍查询</n-button>
        <n-button @click="() => router.push('/battle-info')">战报信息</n-button>
      </n-space>
      <n-divider>同盟战报阵容统计</n-divider>
      <n-form inline :label-width="80" style="justify-content: center;">
        <n-form-item label="名字">
          <n-input v-model:value="playername" placeholder="玩家名字" />
        </n-form-item>
        <n-form-item label="同盟">
          <n-input v-model:value="unionname" placeholder="同盟名称" />
        </n-form-item>
        <n-form-item label="阵容">
          <n-input v-model:value="lineupKey" placeholder="例如 关羽|张飞|赵云" />
        </n-form-item>
        <n-form-item label="角色">
          <n-select v-model:value="role" :options="roleOptions" style="width: 120px" placeholder="选择角色" />
        </n-form-item>
        <n-form-item label="最低等级">
          <n-input-number v-model:value="minLevel" :min="1" style="width: 140px" placeholder="不限制" />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" @click="() => fetchLineups(true)">筛选</n-button>
        </n-form-item>
        <n-form-item>
          <n-button @click="reset">重置</n-button>
        </n-form-item>
        <n-form-item>
          <n-button type="success" @click="downloadCsv">导出CSV</n-button>
        </n-form-item>
      </n-form>
      <div style="margin-top: 12px; color: #9ca3af;">已收集阵容：{{ total }}</div>
      <div class="lineup-list">
        <n-card v-for="item in lineups" :key="item.id" style="margin: 12px auto; max-width: 960px; text-align: left;">
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
        <n-button @click="fetchLineups()" :disabled="!hasMore || loading">{{ hasMore ? '加载更多' : '没有更多了' }}</n-button>
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
