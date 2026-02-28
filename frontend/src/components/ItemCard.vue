<template>
  <div class="item-card" :class="{ 'is-processing': isProcessing }">
  <div class="item-main">
    <span v-if="type === 'formula' && item.status !== 'none_tool'" 
          class="status-indicator" 
          :class="item.status === 'started' ? 'online' : 'offline'">
    </span>
    <div v-if="type === 'cask'" class="app-icon-container">
        <img v-if="item.iconBase64" 
             :src="'data:image/png;base64,' + item.iconBase64" 
             class="app-icon" />
        <span v-else class="app-placeholder">📦</span>
      </div>
    <div class="info-meta">
      <span class="name">{{ item.name }}</span>
      <span class="version">{{ item.version }}</span>
    </div>
  </div>

  <div class="item-actions" v-if="type === 'formula' && item.status !== 'none_tool'">
    <button @click="$emit('action', item)" 
            class="mac-btn"
            :class="item.status === 'started' ? 'btn-stop' : 'btn-start'"
            :disabled="isProcessing">
      <template v-if="isProcessing">
        <span class="mini-loader"></span>
      </template>
      <template v-else>
        {{ item.status === 'started' ? '停止' : '启动' }}
      </template>
    </button>
  </div>
  </div>
</template>

<script setup>
/**
 * 定义 Props
 * @param item - 软件对象数据
 * @param type - 类型：'formula' 或 'cask'
 * @param isProcessing - 当前项是否正在处理中（由外部 processingMap 决定）
 */
defineProps({
  item: { type: Object, required: true },
  type: { type: String, default: 'formula' },
  isProcessing: { type: Boolean, default: false }
})

/**
 * 定义 Emit
 * 当点击按钮时，通知父组件执行具体的 Start/Stop 逻辑
 */
defineEmits(['action'])
</script>

<style scoped>
/* 这里放置 item-card 相关的专用 CSS */
/* 动画和基础样式已经在 main.css 中，这里可以放一些细微的间距调整 */
.info-meta {
  display: flex;
  flex-direction: column; /* 强制垂直排列 */
  align-items: flex-start; /* 左对齐 */
  gap: 2px; /* 名字和版本号之间留一点小缝隙 */
}

/* 6. 状态指示灯 */
.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.online { background: #34C759; box-shadow: 0 0 8px rgba(52, 199, 89, 0.4); }
.offline { background: #FF3B30; opacity: 0.5; }
.name { font-size: 14px; font-weight: 500; color: #fff; }
.version { font-size: 11px; color: #888; font-family: 'SF Mono', Menlo, monospace; }
.item-main {
  display: flex;
  align-items: center;
  gap: 12px;
}
/* 5. 列表项卡片设计 */
.item-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  margin-bottom: 2px;
  border-radius: 8px;
  transition: background 0.2s ease;
}

.item-card:hover {
  background: rgba(255, 255, 255, 0.06);
}
/* 7. 操作按钮 (macOS 风格) */
.mac-btn {
  padding: 5px 12px;
  border-radius: 5px;
  font-size: 12px;
  font-weight: 500;
  border: 0.5px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.1);
  color: #eee;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-start:hover { background: #34C759; border-color: #34C759; color: white; }
.btn-stop:hover { background: #FF3B30; border-color: #FF3B30; color: white; }
/* 处理中状态 */
.is-processing { opacity: 0.5; pointer-events: none; }
.mini-loader {
  width: 12px;
  height: 12px;
  border: 2px solid rgba(255,255,255,0.3);
  border-radius: 50%;
  border-top-color: #fff;
  animation: spin 1s linear infinite;
  display: inline-block;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.app-icon-container {
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.app-icon {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>