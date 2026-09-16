<template>
  <LoginView v-if="!isLoggedIn" @login-success="isLoggedIn = true" />
  <MainView v-else @logout="isLoggedIn = false" />
</template>

<script setup>
import { ref } from 'vue'
import LoginView from './views/LoginView.vue'
import MainView from './views/MainView.vue'

const isLoggedIn = ref(!!localStorage.getItem('token'))
</script>

<style>
*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: #050a18;
  color: #e0e0e0;
}

:root {
  --el-color-primary: #2563EB;
  --el-color-primary-light-3: #60a5fa;
  --el-color-primary-light-5: #93bbfd;
  --el-color-primary-dark-2: #1d4ed8;
  --el-bg-color: #0a0f1e;
  --el-text-color-primary: #f1f5f9;
  --el-text-color-regular: #94a3b8;
  --el-border-color: rgba(255, 255, 255, 0.08);
  --el-fill-color-blank: rgba(255, 255, 255, 0.04);
}

/* Element Plus dark overrides */
.el-message-box {
  background: #0f172a !important;
  border: 1px solid rgba(255,255,255,0.08) !important;
}

.el-message-box__title {
  color: #f1f5f9 !important;
}

.el-message-box__content {
  color: #94a3b8 !important;
}

/* Scrollbar */
::-webkit-scrollbar { width: 5px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 10px; }
::-webkit-scrollbar-thumb:hover { background: rgba(255,255,255,0.18); }

/* Selection */
::selection {
  background: rgba(37, 99, 235, 0.3);
  color: #fff;
}
</style>
