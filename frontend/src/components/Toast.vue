<template>
  <transition name="toast-slide"> 
    <div class="toast-notification" :class="type">
      <div class="toast-content">
        <span class="toast-icon">
          {{ type == 'success' ? '✅' : '❌' }}
        </span>
        <span class="toast-msg">{{ msg }}</span>
      </div>
    </div>
  </transition>
</template>

<script setup>
/**
 * props 接收三个参数：
 * @param msg  - 提示文字
 * @param type - 样式类型 ('success' | 'error')
 */
defineProps({
  msg: { type: String, default: '' },
  type: { type: String, default: 'success' }
})

</script>

<style scoped>
.toast-notification {
  position: fixed;
  top: 40px;
  right: 20px;
  min-width: 240px;
  background: rgba(40, 40, 40, 0.8); /* 半透明暗色 */
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  padding: 14px 18px;
  z-index: 9999;
  /* 确保 Wails 窗口下不会干扰拖拽 */
  --wails-draggable: no-drag !important;
  pointer-events: none; /* 提示框通常不需要点击，防止遮挡下方操作 */
}

/* 成功样式 */
.success {
  border-left: 4px solid #34C759;
}

/* 错误样式 */
.error {
  border-left: 4px solid #FF3B30;
}

.toast-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toast-icon {
  font-size: 16px;
}

.toast-msg {
  font-size: 13px;
  font-weight: 500;
  color: #fff;
  letter-spacing: 0.2px;
}

/* --- 进场和出场动画 --- */

/* 初始状态：透明且向右偏移 */
.toast-slide-enter-from,
.toast-slide-leave-to {
  opacity: 0;
  transform: translateX(30px) scale(0.95);
}

/* 动画过程：使用 macOS 惯用的弹性贝塞尔曲线 */
.toast-slide-enter-active,
.toast-slide-leave-active {
  transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
}

/* 结束状态 */
.toast-slide-enter-to,
.toast-slide-leave-from {
  opacity: 1;
  transform: translateX(0) scale(1);
}

</style>
