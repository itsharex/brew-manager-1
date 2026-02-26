<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useBrew } from './composables/useBrew'
import Toast from './components/Toast.vue'
import ItemCard from './components/ItemCard.vue'

const { 
  data, searchQuery, processingMap, toast,
  filteredFormulae, filteredCasks, 
  updateList, handleService 
} = useBrew()

// 手动刷新按钮逻辑
async function manualRefresh() {
  data.loading = true
  await updateList()
  data.loading = false
  showToast("数据已同步")
}

// 定时器管理
let timer = null
onMounted(() => {
  updateList()
  timer = setInterval(updateList, 10000)
})

onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<template>
  <div class="container">
    <header class="drag-region">
      <div class="header-content">
        <div class="title-group">
          <h2>Brew Manager</h2>
          <span class="sync-tag">Auto Sync ON</span>
        </div>
        
        <div class="toolbar">
          <div class="search-box">
            <span class="search-icon">🔍</span>
            <input v-model="searchQuery" type="text" placeholder="搜索组件..." class="search-input" />
          </div>
          <button @click="manualRefresh" class="btn-refresh" :disabled="data.loading">
            <span v-if="data.loading" class="mini-loader"></span>
            <span v-else>刷新</span>
          </button>
        </div>
      </div>
    </header>

    <div class="main-content">
      <div class="lists-wrapper">
        <section class="list-column">
          <div class="column-header">
            <h3>TERMINAL TOOLS</h3>
            <span class="count-badge">{{ filteredFormulae.length }}</span>
          </div>
          <div class="scroll-area">
            <ItemCard 
              v-for="item in filteredFormulae" 
              :key="item.name"
              type="formula"
              :item="item"
              :is-processing="processingMap.has(item.name)"
              @action="handleService" 
            />
          </div>
        </section>

        <section class="list-column">
          <div class="column-header">
            <h3>GUI APPLICATIONS</h3>
            <span class="count-badge">{{ filteredCasks.length }}</span>
          </div>
          <div class="scroll-area">
            <ItemCard 
              v-for="item in filteredCasks" 
              :key="item.name"
              type="cask"
              :item="item"
            />
          </div>
        </section>
      </div>
    </div>

    <Toast v-if="toast.show" :msg="toast.msg" :type="toast.type" />
  </div>
</template>

<style scoped>
@import "././assets/main.css"
</style>