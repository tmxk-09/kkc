<template>
  <div class="login-container">
    <!-- 动态背景 -->
    <div class="bg-layer">
      <div class="bg-orb orb-1"></div>
      <div class="bg-orb orb-2"></div>
      <div class="bg-orb orb-3"></div>
      <div class="bg-grid"></div>
    </div>

    <!-- ===== 落地页头部 ===== -->
    <header v-if="currentMode === 'landing'" class="landing-header">
      <div class="landing-logo">
        <div class="logo-inner-ring">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="17 8 12 3 7 8"/>
            <line x1="12" y1="3" x2="12" y2="15"/>
          </svg>
        </div>
        <span class="landing-brand">FileFlow</span>
      </div>
      <div class="landing-nav">
        <button class="nav-btn-pickup" @click="showPickupModal = true">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 8v13H3V8"/><path d="M1 3h22v5H1z"/><path d="M10 12h4"/></svg>
          提取文件
        </button>
        <button class="nav-btn-login" @click="toggleToLogin">
          管理登录
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
        </button>
      </div>
    </header>

    <!-- ===== 落地页 Hero 区 ===== -->
    <main v-if="currentMode === 'landing'" class="landing-hero" :class="{ 'hero-visible': heroVisible }">
      <div class="hero-badge">Next-Gen Transfer</div>
      <h1 class="hero-title">重新定义你的<br/><span class="text-gradient">文件传输体验</span></h1>
      <p class="hero-desc">极速、安全、无死角的文件分享中枢。采用自研多线程切片传输引擎，搭配高并发架构完美保障大文件的高效流转，同时支持匿名极速提件闭环。</p>
      
      <div class="hero-features">
        <div class="feat-card">
          <div class="feat-icon blue"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/><polyline points="16 16 12 12 8 16"/></svg></div>
          <h3>极速分片并发</h3>
          <p>榨干千兆带宽极限，支持100GB+超大文件原生上传</p>
        </div>
        <div class="feat-card">
          <div class="feat-icon purple"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 2v6h-6"/><path d="M3 12a9 9 0 0 1 15-6.7L21 8"/><path d="M3 22v-6h6"/><path d="M21 12a9 9 0 0 1-15 6.7L3 16"/></svg></div>
          <h3>完美断点续传</h3>
          <p>无论网络抖动还是页面刷新，上传进度依旧坚若磐石</p>
        </div>
        <div class="feat-card">
          <div class="feat-icon cyan"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg></div>
          <h3>匿名限时提件</h3>
          <p>通过6位安全提货码获取核心资料，到期自动彻底销毁</p>
        </div>
      </div>
    </main>

    <!-- ===== 取件码模态框 ===== -->
    <div class="pickup-modal-overlay" :class="{ 'modal-active': showPickupModal }" @click.self="showPickupModal = false">
      <div class="pickup-modal-card">
        <button class="close-btn" @click="showPickupModal = false"><svg viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" stroke-width="2" fill="none"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
        <div class="pm-icon"><svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 8v13H3V8"/><path d="M1 3h22v5H1z"/><path d="M10 12h4"/></svg></div>
        <h2>输入取件码</h2>
        <p>请输入完整的 6 位字符以提取被共享的文件</p>
        <div class="pm-input-wrapper">
          <input type="text" v-model="pickupCode" placeholder="示例: A1B2C3" maxlength="6" @keyup.enter="handlePickupDownload" />
        </div>
        <button class="pm-submit-btn" :disabled="pickupCode.trim().length !== 6" @click="handlePickupDownload">立即提取</button>
      </div>
    </div>

    <!-- ===== 登录卡片 ===== -->
    <div v-if="currentMode === 'login'" class="login-card" :class="{ 'card-visible': cardVisible }">
      <button class="back-btn" @click="toggleToLanding" title="返回首页">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="19" y1="12" x2="5" y2="12"/><polyline points="12 19 5 12 12 5"/></svg>
      </button>
      
      <!-- Logo 区 -->
      <div class="logo-area">
        <div class="logo-ring">
          <div class="logo-inner">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                 stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="17 8 12 3 7 8"/>
              <line x1="12" y1="3" x2="12" y2="15"/>
            </svg>
          </div>
        </div>
        <h1 class="app-title">FileFlow</h1>
        <p class="app-desc">高性能文件传输系统接入模块</p>
      </div>

      <!-- 分隔线 -->
      <div class="divider">
        <div class="divider-line"></div>
        <span class="divider-text">{{ isRegister ? '通行证注册' : '系统登录' }}</span>
        <div class="divider-line"></div>
      </div>

      <!-- 表单 -->
      <form class="login-form" @submit.prevent="handleSubmit">
        <div class="input-group" :class="{ focused: userFocused, filled: form.username }">
          <label class="input-label">用户名</label>
          <div class="input-wrap">
            <svg class="input-icon" width="18" height="18" viewBox="0 0 24 24" fill="none"
                 stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
            <input type="text" v-model="form.username" placeholder="请输入系统管理员用户"
                   @focus="userFocused = true" @blur="userFocused = false" />
          </div>
        </div>

        <div class="input-group" :class="{ focused: passFocused, filled: form.password }">
          <label class="input-label">验证密钥</label>
          <div class="input-wrap">
            <svg class="input-icon" width="18" height="18" viewBox="0 0 24 24" fill="none"
                 stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
            <input :type="showPass ? 'text' : 'password'" v-model="form.password"
                   placeholder="请输入核准密码"
                   @focus="passFocused = true" @blur="passFocused = false"
                   @keyup.enter="handleLogin" />
            <button type="button" class="toggle-pass" @click="showPass = !showPass"
                    tabindex="-1">
              <svg v-if="!showPass" width="18" height="18" viewBox="0 0 24 24" fill="none"
                   stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
              <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none"
                   stroke="currentColor" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>
            </button>
          </div>
        </div>

        <div class="input-group" v-if="isRegister" :class="{ focused: phoneFocused, filled: form.phone }">
          <label class="input-label">手机号码</label>
          <div class="input-wrap">
            <svg class="input-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
            </svg>
            <input type="text" v-model="form.phone" placeholder="需要接收预警信息"
                   @focus="phoneFocused = true" @blur="phoneFocused = false" />
          </div>
        </div>

        <div class="input-group" v-if="isRegister" :class="{ focused: sqFocused, filled: form.securityQuestion }">
          <label class="input-label">密保问题</label>
          <div class="input-wrap">
            <svg class="input-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path><line x1="12" y1="17" x2="12.01" y2="17"></line>
            </svg>
            <input type="text" v-model="form.securityQuestion" placeholder="例如：您的第一所小学名字"
                   @focus="sqFocused = true" @blur="sqFocused = false" />
          </div>
        </div>

        <div class="input-group" v-if="isRegister" :class="{ focused: saFocused, filled: form.securityAnswer }">
          <label class="input-label">密保答案</label>
          <div class="input-wrap">
            <svg class="input-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"></path>
            </svg>
            <input type="text" v-model="form.securityAnswer" placeholder="请务必牢记以防忘记秘钥"
                   @focus="saFocused = true" @blur="saFocused = false" />
          </div>
        </div>

        <button type="submit" class="login-btn" :class="{ loading }" :disabled="loading">
          <span v-if="!loading" class="btn-content">
            <span>{{ isRegister ? '立即申请' : '安全认证' }}</span>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                 stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="5" y1="12" x2="19" y2="12"/><polyline points="12 5 19 12 12 19"/>
            </svg>
          </span>
          <span v-else class="btn-loading">
            <svg class="spinner" width="20" height="20" viewBox="0 0 24 24"><circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" stroke-width="3" stroke-dasharray="32" stroke-linecap="round"/></svg>
            <span>认证中...</span>
          </span>
        </button>
        
        <!-- 切换登录/注册 -->
        <div class="switch-mode" @click="isRegister = !isRegister">
          {{ isRegister ? '已有通行证？立即登录' : '没有账号？申请注册' }}
        </div>
      </form>

      <div class="footer-info">
        <p>Powered by <strong>FileFlow</strong></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { login, register, getPickupUrl, checkPickupCode } from '../api'
import { ElMessage } from 'element-plus'

const emit = defineEmits(['login-success'])

// 页面模式: landing | login
const currentMode = ref('landing')
const showPickupModal = ref(false)
const pickupCode = ref('')
const heroVisible = ref(false)
const isRegister = ref(false)

const form = ref({ username: '', password: '', phone: '', securityQuestion: '', securityAnswer: '' })
const loading = ref(false)
const showPass = ref(false)
const userFocused = ref(false)
const passFocused = ref(false)
const phoneFocused = ref(false)
const sqFocused = ref(false)
const saFocused = ref(false)
const cardVisible = ref(false)

onMounted(() => {
  setTimeout(() => { heroVisible.value = true }, 100)
})

function toggleToLogin() {
  currentMode.value = 'login'
  heroVisible.value = false
  setTimeout(() => { cardVisible.value = true }, 50)
}

function toggleToLanding() {
  cardVisible.value = false
  setTimeout(() => { 
    currentMode.value = 'landing'
    setTimeout(() => { heroVisible.value = true }, 50)
  }, 300)
}

async function handlePickupDownload() {
  const code = pickupCode.value.trim()
  if (code.length !== 6) {
    ElMessage.warning('请输入完整的 6 位字符以提取被共享的文件')
    return
  }
  
  try {
    const res = await checkPickupCode(code)
    if (res.code === 200) {
      const a = document.createElement('a')
      a.href = getPickupUrl(code)
      a.click()
      showPickupModal.value = false
      pickupCode.value = ''
    } else {
      ElMessage.error(res.msg || '文件不存在或取件码已失效')
    }
  } catch (e) {
    ElMessage.error('网络系统异常，提取失败')
  }
}

async function handleSubmit() {
  if (!form.value.username || !form.value.password) {
    ElMessage.warning('请提供有效的系统凭证')
    return
  }
  if (isRegister.value && (!form.value.phone || !form.value.securityQuestion || !form.value.securityAnswer)) {
    ElMessage.warning('注册需要完整填写手机和安保档案')
    return
  }
  loading.value = true
  try {
    if (isRegister.value) {
      const res = await register(form.value)
      if (res.code === 200) {
        ElMessage.success('注册成功，请重新登录！')
        isRegister.value = false // 帮用户切回登录
        form.value.password = '' // 自动清空密码让重输
      } else {
        ElMessage.error(res.msg || '注册失败')
      }
    } else {
      const res = await login(form.value)
      if (res.code === 200) {
        localStorage.setItem('token', res.data.token)
        const uname = res.data.username || (res.data.user && res.data.user.username) || form.value.username || 'admin'
        localStorage.setItem('username', uname)
        ElMessage.success('验证成功，正在接入控制台')
        emit('login-success')
      } else {
        ElMessage.error(res.msg || '验证失败：非法的凭证')
      }
    }
  } catch (e) {
    ElMessage.error('服务不可达，请检查网关接入')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* ===== 核心变量与前卫动态背景 ===== */
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: #00030a; /* 极深蓝底色 */
  font-family: 'Inter', 'Poppins', sans-serif;
}

.bg-layer {
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.7;
  pointer-events: none;
  mix-blend-mode: screen;
}

.orb-1 {
  width: 700px; height: 700px;
  background: radial-gradient(circle, #00f0ff 0%, transparent 60%);
  top: -20%; left: -10%;
  animation: orbFloat1 22s ease-in-out infinite alternate;
}

.orb-2 {
  width: 650px; height: 650px;
  background: radial-gradient(circle, #ff0055 0%, #7c3aed 50%, transparent 70%);
  bottom: -20%; right: -5%;
  animation: orbFloat2 28s ease-in-out infinite alternate-reverse;
}

.orb-3 {
  width: 400px; height: 400px;
  background: radial-gradient(circle, #3b82f6 0%, transparent 60%);
  top: 40%; left: 40%;
  animation: orbFloat3 15s ease-in-out infinite alternate;
}

.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255,255,255,0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255,255,255,0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  mask-image: radial-gradient(circle at center, black 20%, transparent 80%);
  -webkit-mask-image: radial-gradient(circle at center, black 20%, transparent 80%);
}

@keyframes orbFloat1 {
  0% { transform: translate(0, 0) scale(1) rotate(0); }
  100% { transform: translate(120px, 80px) scale(1.1) rotate(45deg); }
}
@keyframes orbFloat2 {
  0% { transform: translate(0, 0) scale(1); }
  100% { transform: translate(-100px, -100px) scale(1.2); }
}
@keyframes orbFloat3 {
  0% { transform: translate(-50px, 0); opacity: 0.5; }
  100% { transform: translate(80px, 60px); opacity: 0.8; }
}

/* ===== 落地页头 ===== */
.landing-header {
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 5%;
  z-index: 10;
  backdrop-filter: blur(12px);
  background: linear-gradient(180deg, rgba(0,0,0,0.4) 0%, transparent 100%);
}

.landing-logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-inner-ring {
  width: 40px; height: 40px;
  background: linear-gradient(135deg, #00f0ff, #7c3aed);
  border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  color: white;
  box-shadow: 0 0 15px rgba(0,240,255,0.3);
}

.landing-brand {
  font-family: 'Outfit', 'Poppins', sans-serif;
  font-size: 24px;
  font-weight: 800;
  color: #fff;
  letter-spacing: 0.5px;
}

.landing-nav {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-btn-pickup {
  display: flex; align-items: center; gap: 8px;
  background: transparent;
  color: #cbd5e1;
  border: none;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: color 0.3s;
}

.nav-btn-pickup:hover {
  color: #00f0ff;
}

.nav-btn-login {
  display: flex; align-items: center; gap: 6px;
  background: rgba(255,255,255,0.1);
  backdrop-filter: blur(10px);
  color: #fff;
  border: 1px solid rgba(255,255,255,0.15);
  padding: 10px 24px;
  border-radius: 100px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 15px rgba(0,0,0,0.2);
}

.nav-btn-login:hover {
  background: rgba(255,255,255,0.15);
  transform: translateY(-2px);
  border-color: rgba(255,255,255,0.3);
  box-shadow: 0 8px 20px rgba(0,0,0,0.4);
}

/* ===== Hero 区 ===== */
.landing-hero {
  position: relative;
  z-index: 1;
  text-align: center;
  max-width: 900px;
  opacity: 0;
  transform: translateY(30px) scale(0.95);
  transition: all 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.landing-hero.hero-visible {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.hero-badge {
  display: inline-block;
  padding: 6px 16px;
  border-radius: 100px;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.1);
  color: #94a3b8;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 2px;
  text-transform: uppercase;
  margin-bottom: 28px;
  box-shadow: inset 0 1px 1px rgba(255,255,255,0.1);
}

.hero-title {
  font-family: 'Outfit', 'Poppins', sans-serif;
  font-size: 72px;
  font-weight: 800;
  line-height: 1.1;
  color: #fff;
  margin-bottom: 24px;
  letter-spacing: -1.5px;
}

.text-gradient {
  background: linear-gradient(to right, #00f0ff, #ff0055, #7c3aed);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-size: 200% auto;
  animation: gradientShine 5s linear infinite;
}
@keyframes gradientShine {
  0% { background-position: 0% center; }
  100% { background-position: 200% center; }
}

.hero-desc {
  font-size: 18px;
  color: #94a3b8;
  line-height: 1.6;
  max-width: 660px;
  margin: 0 auto 56px;
}

.hero-features {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-top: 24px;
}

.feat-card {
  background: rgba(10, 15, 30, 0.5);
  backdrop-filter: blur(24px) saturate(1.2);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 20px;
  padding: 32px 28px;
  text-align: left;
  transition: all 0.4s ease;
  position: relative;
  overflow: hidden;
  box-shadow: 0 4px 30px rgba(0,0,0,0.1);
}

.feat-card::before {
  content: "";
  position: absolute;
  top: 0; left: 0; width: 100%; height: 100%;
  background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, transparent 100%);
  opacity: 0;
  transition: opacity 0.4s;
}

.feat-card:hover {
  transform: translateY(-8px);
  border-color: rgba(255,255,255,0.15);
  box-shadow: 0 20px 40px rgba(0,0,0,0.3);
}

.feat-card:hover::before {
  opacity: 1;
}

.feat-icon {
  width: 52px; height: 52px;
  border-radius: 16px;
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 24px;
  position: relative;
  z-index: 1;
}
.feat-icon.blue { background: linear-gradient(135deg, rgba(0,240,255,0.2), transparent); color: #00f0ff; border: 1px solid rgba(0,240,255,0.3); }
.feat-icon.purple { background: linear-gradient(135deg, rgba(124,58,237,0.2), transparent); color: #a78bfa; border: 1px solid rgba(124,58,237,0.3); }
.feat-icon.cyan { background: linear-gradient(135deg, rgba(255,0,85,0.2), transparent); color: #ff0055; border: 1px solid rgba(255,0,85,0.3); }

.feat-card h3 {
  font-size: 19px; color: #f8fafc; margin-bottom: 12px; font-weight: 700; position: relative; z-index: 1;
}
.feat-card p {
  font-size: 14px; color: #94a3b8; line-height: 1.6; position: relative; z-index: 1;
}

/* ===== 提取码 Modal ===== */
.pickup-modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0, 5, 15, 0.7);
  backdrop-filter: blur(16px);
  z-index: 100;
  display: flex; align-items: center; justify-content: center;
  opacity: 0; pointer-events: none;
  transition: all 0.4s;
}

.pickup-modal-overlay.modal-active {
  opacity: 1; pointer-events: auto;
}

.pickup-modal-card {
  width: 440px;
  background: rgba(20, 25, 45, 0.7);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 28px;
  padding: 48px 40px;
  text-align: center;
  position: relative;
  transform: scale(0.9);
  transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
  box-shadow: 0 40px 80px -20px rgba(0,0,0,0.8), inset 0 1px 1px rgba(255,255,255,0.1);
}

.modal-active .pickup-modal-card {
  transform: scale(1);
}

.close-btn {
  position: absolute; top: 20px; right: 20px;
  background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1);
  color: rgba(255,255,255,0.6);
  cursor: pointer; padding: 6px; border-radius: 50%;
  transition: all 0.2s;
  display: flex; align-items: center; justify-content: center;
}
.close-btn:hover { background: rgba(255,255,255,0.15); color: #fff; transform: rotate(90deg); }

.pm-icon {
  width: 68px; height: 68px; border-radius: 20px;
  background: linear-gradient(135deg, rgba(0,240,255,0.2), rgba(124,58,237,0.2));
  color: #00f0ff;
  display: flex; align-items: center; justify-content: center;
  margin: 0 auto 24px;
  box-shadow: 0 0 20px rgba(0,240,255,0.2);
}

.pickup-modal-card h2 {
  font-family: 'Outfit', sans-serif;
  font-size: 24px; font-weight: 700; color: #fff; margin-bottom: 10px;
}
.pickup-modal-card p {
  font-size: 14px; color: #94a3b8; margin-bottom: 28px;
}

.pm-input-wrapper input {
  width: 100%; height: 64px;
  background: rgba(0,0,0,0.4);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 16px;
  text-align: center;
  font-family: monospace; font-size: 28px; letter-spacing: 12px;
  color: #fff; text-transform: uppercase;
  outline: none; transition: all 0.3s;
  margin-bottom: 32px;
  box-shadow: inset 0 2px 10px rgba(0,0,0,0.5);
}
.pm-input-wrapper input:focus {
  border-color: #00f0ff; background: rgba(0,240,255,0.05);
  box-shadow: 0 0 0 4px rgba(0,240,255,0.15), inset 0 2px 10px rgba(0,0,0,0.5);
}

.pm-submit-btn {
  width: 100%; height: 56px; border-radius: 16px;
  background: linear-gradient(135deg, #00f0ff, #7c3aed);
  color: #fff; font-size: 17px; font-weight: 700; border: none;
  cursor: pointer; transition: all 0.3s;
  position: relative; overflow: hidden;
}
.pm-submit-btn::after {
  content: ''; position: absolute; top: -50%; left: -50%;
  width: 200%; height: 200%;
  background: rgba(255,255,255,0.2); transform: rotate(45deg) translateY(-100%);
  transition: transform 0.6s ease;
}
.pm-submit-btn:not(:disabled):hover {
  transform: translateY(-2px); box-shadow: 0 12px 24px rgba(124,58,237,0.4);
}
.pm-submit-btn:not(:disabled):hover::after {
  transform: rotate(45deg) translateY(100%);
}
.pm-submit-btn:disabled {
  background: rgba(255,255,255,0.05); color: rgba(255,255,255,0.3); border: 1px solid rgba(255,255,255,0.1); cursor: not-allowed;
}

/* ===== 核心：登录卡片 Glassmorphism ===== */
.login-card {
  position: relative;
  z-index: 1;
  width: 420px;
  padding: 48px 44px 40px;
  background: rgba(12, 17, 30, 0.55);
  backdrop-filter: blur(48px) saturate(1.8);
  border-radius: 32px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  box-shadow:
    0 30px 60px -12px rgba(0, 0, 0, 0.8),
    inset 0 1px 1px rgba(255, 255, 255, 0.1);
  opacity: 0;
  transform: translateY(30px) scale(0.95);
  transition: all 0.7s cubic-bezier(0.16, 1, 0.3, 1);
}

.login-card::before {
  content: ""; position: absolute; inset: -1px; border-radius: 33px;
  background: linear-gradient(135deg, rgba(255,255,255,0.2) 0%, transparent 100%);
  z-index: -1; pointer-events: none; -webkit-mask-image: linear-gradient(white, black);
}

.login-card.card-visible {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.back-btn {
  position: absolute; top: 20px; left: 20px;
  background: transparent; border: 1px solid transparent; color: #94a3b8;
  cursor: pointer; padding: 8px; border-radius: 12px;
  transition: all 0.2s;
  display: flex; align-items: center;
}
.back-btn:hover { background: rgba(255,255,255,0.05); color: #fff; border-color: rgba(255,255,255,0.1); transform: translateX(-4px); }

.logo-area { text-align: center; margin-bottom: 32px; }
.logo-ring {
  display: inline-flex; align-items: center; justify-content: center;
  width: 64px; height: 64px; border-radius: 20px;
  background: linear-gradient(135deg, #00f0ff, #ff0055); position: relative; margin-bottom: 16px;
  box-shadow: 0 10px 25px -4px rgba(255, 0, 85, 0.4);
}
.logo-ring::before {
  content: ''; position: absolute; inset: -4px; border-radius: 24px;
  background: inherit;
  z-index: -1; filter: blur(12px); opacity: 0.6; animation: loginLogoGlow 3s ease-in-out infinite alternate;
}
@keyframes loginLogoGlow {
  0% { transform: scale(0.95); opacity: 0.4; }
  100% { transform: scale(1.05); opacity: 0.7; }
}
.logo-inner { color: white; display: flex; align-items: center; justify-content: center; }

.app-title { font-family: 'Outfit', sans-serif; font-size: 32px; font-weight: 800; color: #fff; margin: 0 0 6px; letter-spacing: 0.5px; }
.app-desc { color: #94a3b8; font-size: 14px; margin: 0; font-weight: 400; }

.divider { display: flex; align-items: center; gap: 16px; margin-bottom: 28px; }
.divider-line { flex: 1; height: 1px; background: linear-gradient(90deg, transparent, rgba(255,255,255,0.1), transparent); }
.divider-text { font-size: 11px; color: #64748b; font-weight: 700; text-transform: uppercase; letter-spacing: 3px; }

.input-group { margin-bottom: 22px; position: relative; }
.input-label { display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 8px; text-transform: uppercase; letter-spacing: 1px; transition: color 0.3s; }
.input-group.focused .input-label { color: #00f0ff; }

.input-wrap {
  position: relative; display: flex; align-items: center;
  background: rgba(0, 0, 0, 0.3); border: 1px solid rgba(255, 255, 255, 0.08); 
  border-radius: 16px; transition: all 0.3s ease; 
  overflow: hidden;
  box-shadow: inset 0 2px 4px rgba(0,0,0,0.3);
}

.input-wrap::after {
  content: ''; position: absolute; bottom: 0; left: 0; right: 0; height: 2px;
  background: linear-gradient(90deg, #00f0ff, #7c3aed);
  transform: scaleX(0); transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1);
  transform-origin: center;
}
.input-group.focused .input-wrap { 
  border-color: rgba(255, 255, 255, 0.2); 
  background: rgba(0, 0, 0, 0.5);
  box-shadow: inset 0 2px 4px rgba(0,0,0,0.5), 0 8px 16px rgba(0,0,0,0.1); 
}
.input-group.focused .input-wrap::after { transform: scaleX(1); }

.input-icon { margin-left: 18px; color: #64748b; flex-shrink: 0; transition: color 0.3s, transform 0.3s; }
.input-group.focused .input-icon { color: #00f0ff; transform: scale(1.1); }
.input-wrap input { flex: 1; padding: 16px; background: transparent; border: none; outline: none; font-size: 15px; color: #f8fafc; font-weight: 500;}
.input-wrap input:-webkit-autofill,
.input-wrap input:-webkit-autofill:hover,
.input-wrap input:-webkit-autofill:focus,
.input-wrap input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 50px transparent inset !important;
  -webkit-text-fill-color: #f8fafc !important;
  transition: background-color 5000s ease-in-out 0s;
}
.input-wrap input::placeholder { color: #475569; font-weight: 400; }
.toggle-pass { background: none; border: none; padding: 0 16px; color: #64748b; cursor: pointer; display: flex; align-items: center; transition: color 0.2s; }
.toggle-pass:hover { color: #cbd5e1; }

.login-btn { 
  width: 100%; padding: 16px 24px; border: none; border-radius: 16px; 
  font-size: 16px; font-weight: 700; color: white; cursor: pointer; 
  background: linear-gradient(135deg, #00f0ff, #7c3aed); 
  position: relative; overflow: hidden; transition: all 0.3s ease; 
  margin-top: 10px; font-family: 'Outfit', sans-serif; letter-spacing: 1px;
  box-shadow: 0 4px 15px rgba(124,58,237,0.3);
}
.login-btn::after {
  content: ''; position: absolute; top: -50%; left: -50%;
  width: 200%; height: 200%;
  background: linear-gradient(to right, transparent, rgba(255,255,255,0.4), transparent);
  transform: rotate(45deg) translateY(-100%);
  transition: transform 0.8s ease;
}
.login-btn:hover:not(:disabled) { 
  transform: translateY(-2px); 
  box-shadow: 0 15px 30px rgba(124,58,237,0.5); 
}
.login-btn:hover:not(:disabled)::after { transform: rotate(45deg) translateY(100%); }
.login-btn:active:not(:disabled) { transform: translateY(1px); }
.login-btn:disabled { cursor: not-allowed; opacity: 0.6; filter: grayscale(50%); }

.btn-content, .btn-loading { position: relative; z-index: 1; display: flex; align-items: center; justify-content: center; gap: 10px; }
.spinner { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.switch-mode { 
  margin-top: 24px; text-align: center; color: #94a3b8; 
  font-size: 14px; font-weight: 600; cursor: pointer; transition: color 0.3s; 
}
.switch-mode:hover { color: #00f0ff; text-decoration: underline; text-underline-offset: 4px; }

.footer-info { margin-top: 32px; text-align: center; }
.footer-info p { font-size: 12px; color: #475569; margin: 0; }
.footer-info strong { color: #64748b; font-family: 'Outfit', sans-serif;}

@media (prefers-reduced-motion: reduce) {
  .bg-orb, .logo-ring::before, .text-gradient { animation: none; }
  .login-card { transition: opacity 0.3s; transform: none; }
  .login-card.card-visible { transform: none; }
}
</style>