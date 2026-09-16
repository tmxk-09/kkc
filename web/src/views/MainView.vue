<template>
  <div class="app-layout"
       @dragover.prevent="handleGlobalDragOver"
       @dragleave.prevent="handleGlobalDragLeave"
       @drop.prevent="handleGlobalDrop">
    
    <!-- 全局拖拽遮罩 -->
    <div v-show="isDragging" class="global-dropzone">
      <div class="dropzone-content">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="17 8 12 3 7 8"/>
          <line x1="12" y1="3" x2="12" y2="15"/>
        </svg>
        <p>释放文件以在此处上传</p>
      </div>
    </div>
    <!-- ===== 侧边栏 ===== -->
    <aside class="sidebar">
      <div class="sidebar-brand">
        <div class="brand-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor"
               stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="17 8 12 3 7 8"/>
            <line x1="12" y1="3" x2="12" y2="15"/>
          </svg>
        </div>
        <span class="brand-name">FileFlow</span>
      </div>

      <nav class="sidebar-nav">
        <a class="nav-item" :class="{active: currentTab === 'files'}" href="#" @click.prevent="currentTab = 'files'">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
          <span>文件管理</span>
        </a>
        <a v-if="isAdmin" class="nav-item" :class="{active: currentTab === 'users'}" href="#" @click.prevent="currentTab = 'users'; loadUsers()">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
          <span>用户管理</span>
        </a>
        <a class="nav-item" :class="{active: currentTab === 'tickets'}" href="#" @click.prevent="currentTab = 'tickets'; loadTickets()">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path></svg>
          <span>工单支持</span>
        </a>
        <a class="nav-item" :class="{active: currentTab === 'messages'}" href="#" @click.prevent="currentTab = 'messages'; loadMessages()">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
          <span style="display:flex; align-items:center; flex:1">系统消息
            <span v-if="unreadMsgCount > 0" style="margin-left:auto; background:#ef4444; color:#fff; font-size:10px; padding:2px 6px; border-radius:10px; line-height:1;">{{ unreadMsgCount }}</span>
          </span>
        </a>
      </nav>

      <div class="sidebar-footer">
        <div class="footer-top">
          <div class="user-card">
            <div class="user-avatar">{{ avatarChar }}</div>
            <div class="user-meta">
              <span class="user-name" :title="displayName">{{ displayName }}</span>
              <span class="user-role">{{ isAdmin ? '系统管理员' : '普通用户' }}</span>
            </div>
          </div>
          <button class="logout-btn" @click="handleLogout" title="退出登录">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          </button>
        </div>
        <div class="user-storage-panel" v-if="currentUser">
          <div class="us-text">
            <span>{{ formatFileSize(currentUser.usedSpace) }} / {{ formatFileSize(currentUser.totalSpace) }}</span>
            <span class="us-percent">{{ Math.min(100, Math.round((currentUser.usedSpace / currentUser.totalSpace) * 100)) || 0 }}%</span>
          </div>
          <div class="us-bar-bg">
            <div class="us-bar-fill" :style="{ width: (Math.min(100, Math.round((currentUser.usedSpace / currentUser.totalSpace) * 100)) || 0) + '%' }" :class="{'danger': (currentUser.usedSpace / currentUser.totalSpace) > 0.9}"></div>
          </div>
          <button class="gift-space-btn" @click="showGiftForm = true" title="赠送存储空间给其他用户">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 12 20 22 4 22 4 12"></polyline><rect x="2" y="7" width="20" height="5"></rect><line x1="12" y1="22" x2="12" y2="7"></line><path d="M12 7H7.5a2.5 2.5 0 0 1 0-5C11 2 12 7 12 7z"></path><path d="M12 7h4.5a2.5 2.5 0 0 0 0-5C13 2 12 7 12 7z"></path></svg>
            <span>空间转赠</span>
          </button>
        </div>
      </div>
    </aside>

    <!-- ===== 主内容区 ===== -->
    <main class="main-area">
      <!-- 顶部状态栏 -->
      <header v-show="currentTab === 'files'" class="top-header">
        <div class="header-left">
          <h1 class="page-title">文件管理</h1>
          <p class="page-desc">全屏拖拽即可上传文件，或使用取件码快速提取</p>
        </div>
        <div class="header-actions">
          <!-- 取件码提取 -->
          <div class="pickup-wrap">
            <input v-model="pickupCodeInput" class="pickup-input" placeholder="输入6位取件码..." maxlength="6" />
            <button class="pickup-btn" @click="handlePickup">提取</button>
          </div>

          <!-- 上传操作组合 -->
          <div class="upload-action-group">
            <span class="expire-label" title="上传自动销毁期">过期:</span>
            <select v-model="selectedExpireDays" class="expire-select">
              <option v-for="opt in expireOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
            <input type="file" ref="fileInput" multiple @change="handleFileSelect" style="display:none" />
            <button class="upload-btn-primary" @click="$refs.fileInput.click()">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
              <span>选择文件上传</span>
            </button>
          </div>

          <div class="stats-row">
            <div class="stat-card">
              <span class="stat-value">{{ total }}</span>
              <span class="stat-label">文件总数</span>
            </div>
            <div class="stat-card">
              <span class="stat-value">{{ uploadTasks.length }}</span>
              <span class="stat-label">上传任务</span>
            </div>
          </div>
        </div>
      </header>

      <!-- 上传任务列表 -->
      <section v-if="currentTab === 'files' && uploadTasks.length" class="task-section">
        <h3 class="section-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
          上传任务
        </h3>
        <div class="task-grid">
          <div v-for="task in uploadTasks" :key="task.id" class="task-card"
               :class="'status-' + task.status">
            <div class="task-header">
              <div class="task-file-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
              </div>
              <div class="task-meta">
                <span class="task-name">{{ task.fileName }}</span>
                <span class="task-size">{{ formatFileSize(task.fileSize) }}</span>
              </div>
              <span class="task-badge" :class="'badge-' + task.status">
                {{ statusText(task.status) }}
              </span>
            </div>
            <div class="task-progress-bar">
              <div class="progress-fill" :style="{ width: task.progress + '%' }"></div>
            </div>
            <div class="task-footer">
              <span class="task-percent">{{ task.progress }}%</span>
              <span v-if="task.uploadedBytes && task.status === 'uploading'" class="task-uploaded">
                {{ formatFileSize(task.uploadedBytes) }} / {{ formatFileSize(task.fileSize) }}
              </span>
              <span v-if="task.speed && task.status === 'uploading'" class="task-speed">
                {{ formatSpeed(task.speed) }}
              </span>

              <!-- 实时用时变化与剩余时间 -->
              <span v-if="task.status === 'uploading' || task.status === 'hashing' || task.status === 'merging'" class="task-time">
                ⏱️ {{ formatDuration(task.elapsedSeconds) }}
                <span v-if="task.remainingSeconds && task.status === 'uploading'" class="task-remaining">(余 {{ formatDuration(task.remainingSeconds) }})</span>
              </span>

              <span v-if="task.status === 'queued'" class="task-queued-hint">⏳ 排队等待中 (限制同时5个任务)</span>
              <span v-if="task.status === 'hashing'" class="task-hashing-hint">正在计算文件Hash...</span>
              <span v-if="task.status === 'merging'" class="task-merging-hint">正在合并分片...</span>
              <span v-if="task.status === 'interrupted'" class="task-interrupted-hint" style="color:#f59e0b">已中断，请选文件恢复</span>
              <span v-if="task.status === 'instant'" class="task-instant">✓ 秒传完成 (耗时: {{ formatDuration(task.elapsedSeconds) }})</span>
              <span v-if="task.status === 'done'" class="task-done-hint">✓ 上传完成 (总耗时: {{ formatDuration(task.elapsedSeconds) }})</span>

              <!-- 🔥 删除/取消按钮（已中断、排队或失败的任务显示） -->
              <button v-if="task.status === 'interrupted' || task.status === 'error' || task.status === 'queued'"
                      @click="handleDeleteTask(task)"
                      class="task-delete-btn"
                      :title="task.status === 'queued' ? '取消排队' : '删除此任务'">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  <line x1="10" y1="11" x2="10" y2="17"/>
                  <line x1="14" y1="11" x2="14" y2="17"/>
                </svg>
                {{ task.status === 'queued' ? '取消' : '删除' }}
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 文件列表 -->
      <section v-show="currentTab === 'files'" class="file-section">
        <div class="section-toolbar">
          <h3 class="section-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
            所有文件
          </h3>
          <div class="search-wrap">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
            <input v-model="searchKeyword" placeholder="搜索文件名..." @input="loadFileList" />
          </div>
        </div>

        <div class="file-table-wrap" v-loading="listLoading">
          <table class="file-table" v-if="fileList.length">
            <thead>
              <tr>
                <th style="width: 31%">文件名</th>
                <th style="width: 12%">大小</th>
                <th style="width: 15%">取件码</th>
                <th style="width: 18%">上传时间</th>
                <th style="width: 8%">下载次数</th>
                <th style="width: 16%">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in fileList" :key="row.id" class="file-row">
                <td>
                  <div class="file-name-cell">
                    <div class="file-ext-badge">{{ getExt(row.fileName) }}</div>
                    <span class="file-name-text" :title="row.fileName">{{ row.fileName }}</span>
                  </div>
                </td>
                <td class="text-muted">{{ formatFileSize(row.fileSize) }}</td>
                <td>
                  <span v-if="row.pickupCode" class="pickup-badge" title="点击复制" @click="copyCode(row.pickupCode)">
                    {{ row.pickupCode }}
                  </span>
                  <span v-else class="text-muted">-</span>
                </td>
                <td class="text-muted">{{ formatDate(row.createdAt) }}</td>
                <td class="text-center">{{ row.downloads || 0 }}</td>
                <td>
                  <div class="action-btns">
                    <!-- 🚀 极速多线程并发下载 (推荐) -->
                    <button class="action-btn speed-btn" @click="handleSpeedDownload(row)" title="极速并发下载 (多协程切片加速，榨干带宽)">
                      <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
                    </button>
                    <!-- 原生普通下载 -->
                    <button class="action-btn download" @click="handleNormalDownload(row)" title="普通单流下载">
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                    </button>
                    <!-- 终端命令行一键复制 -->
                    <button class="action-btn cli-btn" @click="copyDownloadCommand(row)" title="复制 Linux/Mac 终端极速下载命令 (curl/aria2)">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>
                    </button>
                    <!-- 删除 -->
                    <button class="action-btn delete" @click="handleDelete(row)" title="删除">
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-else class="empty-state">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                 stroke-width="1.5" opacity="0.2">
              <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
            </svg>
            <p>暂无文件</p>
            <span>上传文件后将显示在此处</span>
          </div>
        </div>

        <div v-if="total > 0" class="pagination-bar">
          <span class="page-info">共 {{ total }} 个文件</span>
          <el-pagination v-model:current-page="pageNum" v-model:page-size="pageSize"
                         :total="total" layout="prev, pager, next"
                         small @change="loadFileList" />
        </div>
      </section>

      <!-- ===== 管理员：用户管理视图 ===== -->
      <section v-if="currentTab === 'users' && isAdmin" class="file-section user-section-card" style="margin-top:20px">
        <div class="section-toolbar">
          <h3 class="section-title">
            <div class="title-icon-wrap"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle></svg></div>
            <span>用户管理与配额</span>
          </h3>
          <div style="display:flex; gap:10px;">
            <button class="upload-btn-primary" @click="handleOpSendMsg('all')" style="background: rgba(16, 185, 129, 0.2); border: 1px solid rgba(16, 185, 129, 0.5); color: #10b981;">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
              <span style="margin-left:6px">全员群发</span>
            </button>
            <button class="upload-btn-primary" @click="loadUsers">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"></polyline><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>
              <span style="margin-left:6px">刷新数据</span>
            </button>
          </div>
        </div>

        <div class="table-card" v-loading="userListLoading">
          <table class="modern-table" v-if="userListData.length">
            <thead>
              <tr>
                <th style="width: 10%">ID</th>
                <th style="width: 18%">通行账号</th>
                <th style="width: 24%">资产占用 (G)</th>
                <th style="width: 18%">分配时间</th>
                <th style="width: 12%">当前状态</th>
                <th style="width: 18%">管控操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in userListData" :key="user.id" class="table-row">
                <td class="text-muted"><span class="id-badge">#{{ user.id }}</span></td>
                <td>
                  <div class="user-account-cell">
                    <div class="user-avatar-mini">{{ user.username.charAt(0).toUpperCase() }}</div>
                    <span class="user-name-text">{{ user.username }}</span>
                  </div>
                </td>
                <td>
                  <div class="space-capsule">
                    <div class="space-text">
                      <span class="used">{{ formatFileSize(user.usedSpace) }}</span>
                      <span class="divider">/</span>
                      <span class="total">{{ formatFileSize(user.totalSpace) }}</span>
                    </div>
                    <div class="space-bar-bg">
                      <div class="space-bar-fill" :style="{ width: Math.min(100, (user.usedSpace / user.totalSpace) * 100) + '%' }"
                           :class="{ 'danger-fill': (user.usedSpace / user.totalSpace) > 0.8 }"></div>
                    </div>
                    <button class="reset-space-btn" @click.prevent="handleEditSpace(user)" title="调整配额">
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
                    </button>
                  </div>
                </td>
                <td class="text-muted">{{ formatDate(user.createdAt) }}</td>
                <td>
                  <span v-if="user.status === 1" class="status-badge success">● 正常运作</span>
                  <span v-else class="status-badge error">● 账户封禁</span>
                </td>
                <td>
                  <div class="action-btns" v-if="user.username !== 'admin'">
                    <button class="action-btn-modern space-btn" @click="handleEditSpace(user)" title="设置用户空间配额">
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
                    </button>
                    <button class="action-btn-modern" @click="handleOpSendMsg(user.id)" title="发私信通知">
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
                    </button>
                    <button class="action-btn-modern" @click="handleResetPassword(user)" title="重置密码">
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"></path></svg>
                    </button>
                    <button class="action-btn-modern" :class="user.status === 1 ? 'freeze-btn' : 'unfreeze-btn'" @click="handleToggleStatus(user)" :title="user.status === 1 ? '一键冻结' : '解除封解'">
                      <svg v-if="user.status === 1" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
                      <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <div v-else class="empty-state">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.3"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle></svg>
            <p>尚未存在子账户</p>
          </div>
        </div>
      </section>
      <!-- ===== 工单视图 ===== -->
      <section v-if="currentTab === 'tickets'" class="file-section user-section-card" style="margin-top:20px">
        <div class="section-toolbar">
          <h3 class="section-title">
            <div class="title-icon-wrap"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path></svg></div>
            <span>工单支持 ({{ username === 'admin' ? '处理中心' : '我的反馈' }})</span>
          </h3>
          <div style="display:flex;gap:10px;">
            <button v-if="username !== 'admin'" class="upload-btn-primary" @click="showTicketForm = true">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
              <span style="margin-left:6px">新建工单</span>
            </button>
            <button class="upload-btn-primary" @click="loadTickets" style="background: rgba(255,255,255,0.1)">
              <span style="margin-left:6px; color:#fff">刷新页面</span>
            </button>
          </div>
        </div>

        <div class="table-card" v-loading="ticketLoading">
          <table class="modern-table" v-if="ticketList.length">
            <thead>
              <tr>
                <th style="width: 8%">ID</th>
                <th v-if="username === 'admin'" style="width: 12%">提问用户</th>
                <th style="width: 20%">标题</th>
                <th style="width: 25%">内容</th>
                <th style="width: 20%">回复</th>
                <th style="width: 10%">状态</th>
                <th style="width: 10%" v-if="username === 'admin'">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="ticket in ticketList" :key="ticket.id" class="table-row">
                <td class="text-muted"><span class="id-badge">#{{ ticket.id }}</span></td>
                <td v-if="username === 'admin'">
                   <span class="user-name-text">{{ ticket.username }}</span>
                </td>
                <td style="font-weight:600">{{ ticket.title }}</td>
                <td style="white-space:pre-wrap; font-size:12px; color:rgba(255,255,255,0.7)">{{ ticket.content }}</td>
                <td style="white-space:pre-wrap; font-size:12px; color:#10b981">{{ ticket.replyContent || '-' }}</td>
                <td>
                  <span v-if="ticket.status === 1" class="status-badge success">已回复</span>
                  <span v-else class="status-badge error">待回复</span>
                </td>
                <td v-if="username === 'admin'">
                  <button class="action-btn-modern" @click="handleOpReply(ticket)" :title="ticket.status === 0 ? '回复' : '追加回复'">
                    <svg v-if="ticket.status === 0" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 17 4 12 9 7"/><path d="M20 18v-2a4 4 0 0 0-4-4H4"/></svg>
                    <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
          <div v-else class="empty-state">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.3"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path></svg>
            <p>暂无相关工单</p>
          </div>
        </div>
      </section>

      <!-- 原生自制模态弹窗 -->
      <div v-if="showTicketForm" style="position:fixed; top:0; left:0; width:100vw; height:100vh; z-index: 1000; background:rgba(0,0,0,0.8); display:flex; align-items:center; justify-content:center;">
        <div style="background:rgba(15,23,42,1); border:1px solid rgba(255,255,255,0.1); border-radius:12px; padding:24px; width:400px; max-width:90%">
          <h3 style="color:#fff; margin-top:0; margin-bottom:16px;">提交工单</h3>
          <div style="margin-bottom:12px;">
            <label style="display:block; color:rgba(255,255,255,0.6); font-size:12px; margin-bottom:4px;">标题</label>
            <input v-model="newTicket.title" placeholder="简述您的问题..." style="width:100%; box-sizing:border-box; background:rgba(255,255,255,0.05); border:1px solid rgba(255,255,255,0.1); color:#fff; border-radius:6px; padding:8px 12px; font-size:14px; outline:none;" />
          </div>
          <div style="margin-bottom:20px;">
            <label style="display:block; color:rgba(255,255,255,0.6); font-size:12px; margin-bottom:4px;">详细内容</label>
            <textarea v-model="newTicket.content" rows="4" placeholder="请详细描述您遇到的情况..." style="width:100%; box-sizing:border-box; background:rgba(255,255,255,0.05); border:1px solid rgba(255,255,255,0.1); color:#fff; border-radius:6px; padding:8px 12px; font-size:14px; outline:none; resize:vertical;"></textarea>
          </div>
          <div style="display:flex; justify-content:flex-end; gap:12px;">
            <button @click="showTicketForm = false" style="background:transparent; border:1px solid rgba(255,255,255,0.2); color:#fff; padding:6px 16px; border-radius:6px; cursor:pointer;">取消</button>
            <button @click="handleSubmitTicket" style="background:linear-gradient(135deg, #2563EB, #7c3aed); border:none; color:#fff; padding:6px 16px; border-radius:6px; cursor:pointer; font-weight:600;" :disabled="submittingTicket">{{ submittingTicket ? '提交中...' : '提交' }}</button>
          </div>
        </div>
      </div>
      <!-- ===== 消息通知视图 ===== -->
      <section v-if="currentTab === 'messages'" class="file-section user-section-card" style="margin-top:20px">
        <div class="section-toolbar">
          <h3 class="section-title">
            <div class="title-icon-wrap"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg></div>
            <span>系统通知</span>
          </h3>
          <div style="display:flex;gap:10px;">
            <button class="upload-btn-primary" @click="handleReadAllMsg" v-if="unreadMsgCount > 0" style="background:#3b82f6; border-color:#3b82f6;">
              <span style="margin-left:6px; color:#fff">全部标记已读</span>
            </button>
            <button class="upload-btn-primary" @click="loadMessages" style="background: rgba(255,255,255,0.1)">
              <span style="margin-left:6px; color:#fff">刷新消息</span>
            </button>
          </div>
        </div>

        <div class="table-card" v-loading="msgLoading">
          <div v-if="msgList.length" style="display:flex;flex-direction:column;gap:12px;">
            <div v-for="msg in msgList" :key="msg.id" style="background:rgba(255,255,255,0.02); padding:16px; border-radius:10px; border-left:4px solid transparent; transition:all 0.2s; position:relative; cursor:pointer;" :style="msg.isRead === 0 ? 'border-left-color:#3b82f6; background:rgba(59,130,246,0.05);' : ''" @click="handleReadMsg(msg)">
               <h4 style="margin:0 0 8px 0; color:#fff; font-size:15px; display:flex; align-items:center;">
                 <span v-if="msg.isRead===0" style="width:8px;height:8px;background:#3b82f6;border-radius:50%;margin-right:8px;"></span>
                 {{ msg.title }}
               </h4>
               <p style="margin:0 0 10px 0; font-size:13px; color:rgba(255,255,255,0.6); white-space:pre-wrap; line-height:1.5;">{{ msg.content }}</p>
               <div style="font-size:12px; color:rgba(255,255,255,0.3);">{{ formatDate(msg.createdAt) }}</div>
            </div>
          </div>
          <div v-else class="empty-state">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.3"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
            <p>暂无通知消息</p>
          </div>
        </div>
      </section>

      <!-- 发消息原生弹窗 -->
      <div v-if="showMsgForm" style="position:fixed; top:0; left:0; width:100vw; height:100vh; z-index: 1000; background:rgba(0,0,0,0.8); display:flex; align-items:center; justify-content:center;">
        <div style="background:rgba(15,23,42,1); border:1px solid rgba(255,255,255,0.1); border-radius:12px; padding:24px; width:400px; max-width:90%">
          <h3 style="color:#fff; margin-top:0; margin-bottom:16px;">{{ msgTarget === 'all' ? '发起全员系统广播' : '发送私信通知' }}</h3>
          <div style="margin-bottom:12px;">
            <label style="display:block; color:rgba(255,255,255,0.6); font-size:12px; margin-bottom:4px;">标题</label>
            <input v-model="newMsg.title" placeholder="输入通知标题..." style="width:100%; box-sizing:border-box; background:rgba(255,255,255,0.05); border:1px solid rgba(255,255,255,0.1); color:#fff; border-radius:6px; padding:8px 12px; font-size:14px; outline:none;" />
          </div>
          <div style="margin-bottom:20px;">
            <label style="display:block; color:rgba(255,255,255,0.6); font-size:12px; margin-bottom:4px;">详细内容</label>
            <textarea v-model="newMsg.content" rows="4" placeholder="消息正文..." style="width:100%; box-sizing:border-box; background:rgba(255,255,255,0.05); border:1px solid rgba(255,255,255,0.1); color:#fff; border-radius:6px; padding:8px 12px; font-size:14px; outline:none; resize:vertical;"></textarea>
          </div>
          <div style="display:flex; justify-content:flex-end; gap:12px;">
            <button @click="showMsgForm = false" style="background:transparent; border:1px solid rgba(255,255,255,0.2); color:#fff; padding:6px 16px; border-radius:6px; cursor:pointer;">取消</button>
            <button @click="handleSubmitMsg" style="background:linear-gradient(135deg, #10b981, #059669); border:none; color:#fff; padding:6px 16px; border-radius:6px; cursor:pointer; font-weight:600;" :disabled="submittingMsg">{{ submittingMsg ? '发送中...' : '确定发送' }}</button>
          </div>
        </div>
      </div>
      <!-- 空间赠送原生弹窗 -->
      <div v-if="showGiftForm" style="position:fixed; top:0; left:0; width:100vw; height:100vh; z-index: 1000; background:rgba(0,0,0,0.8); display:flex; align-items:center; justify-content:center;">
        <div style="background:rgba(15,23,42,1); border:1px solid rgba(255,255,255,0.1); border-radius:12px; padding:24px; width:420px; max-width:90%">
          <h3 style="color:#fff; margin-top:0; margin-bottom:16px; display:flex; align-items:center; gap:8px;">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2"><polyline points="20 12 20 22 4 22 4 12"></polyline><rect x="2" y="7" width="20" height="5"></rect><line x1="12" y1="22" x2="12" y2="7"></line><path d="M12 7H7.5a2.5 2.5 0 0 1 0-5C11 2 12 7 12 7z"></path><path d="M12 7h4.5a2.5 2.5 0 0 0 0-5C13 2 12 7 12 7z"></path></svg>
            <span>赠送上传存储空间</span>
          </h3>
          <div style="margin-bottom:12px;">
            <label style="display:block; color:rgba(255,255,255,0.6); font-size:12px; margin-bottom:4px;">接收方用户名</label>
            <input v-model="giftForm.toUsername" placeholder="输入对方用户名..." style="width:100%; box-sizing:border-box; background:rgba(255,255,255,0.05); border:1px solid rgba(255,255,255,0.1); color:#fff; border-radius:6px; padding:8px 12px; font-size:14px; outline:none;" />
          </div>
          <div style="margin-bottom:12px;">
            <label style="display:block; color:rgba(255,255,255,0.6); font-size:12px; margin-bottom:4px;">赠送容量大小 (MB)</label>
            <input v-model.number="giftForm.giftSizeMb" type="number" min="1" placeholder="如 1024 表示 1GB" style="width:100%; box-sizing:border-box; background:rgba(255,255,255,0.05); border:1px solid rgba(255,255,255,0.1); color:#fff; border-radius:6px; padding:8px 12px; font-size:14px; outline:none;" />
            <div v-if="currentUser" style="font-size:11px; color:rgba(255,255,255,0.4); margin-top:4px;">
              当前剩余可用额度：{{ formatFileSize(Math.max(0, currentUser.totalSpace - currentUser.usedSpace)) }}
            </div>
          </div>
          <div style="margin-bottom:20px;">
            <label style="display:block; color:rgba(255,255,255,0.6); font-size:12px; margin-bottom:4px;">登录密码确认（安全确权）</label>
            <input v-model="giftForm.verifyPassword" type="password" placeholder="请输入您的登录密码确认..." style="width:100%; box-sizing:border-box; background:rgba(255,255,255,0.05); border:1px solid rgba(255,255,255,0.1); color:#fff; border-radius:6px; padding:8px 12px; font-size:14px; outline:none;" />
          </div>
          <div style="display:flex; justify-content:flex-end; gap:12px;">
            <button @click="showGiftForm = false" style="background:transparent; border:1px solid rgba(255,255,255,0.2); color:#fff; padding:6px 16px; border-radius:6px; cursor:pointer;">取消</button>
            <button @click="handleGiftSpace" style="background:linear-gradient(135deg, #2563EB, #7c3aed); border:none; color:#fff; padding:6px 16px; border-radius:6px; cursor:pointer; font-weight:600;" :disabled="submittingGift">{{ submittingGift ? '处理中...' : '确认赠送' }}</button>
          </div>
        </div>
      </div>
      <!-- 🚀 极速多线程切片下载浮窗（文叔叔式并发加速） -->
      <transition name="el-zoom-in-bottom">
        <div v-if="speedDownloadTask.visible" class="speed-download-panel">
          <div class="sdp-header">
            <div class="sdp-title-box">
              <span class="sdp-badge">⚡ 极速并发引擎</span>
              <span class="sdp-filename" :title="speedDownloadTask.fileName">{{ speedDownloadTask.fileName }}</span>
            </div>
            <button class="sdp-close-btn" @click="speedDownloadTask.visible = false" title="关闭面板（后台继续下载）">✕</button>
          </div>

          <div class="sdp-body">
            <div class="sdp-stats">
              <div class="sdp-stat-main">
                <span class="sdp-speed-value">{{ formatSpeed(speedDownloadTask.speed) }}</span>
                <span class="sdp-status-tag" :class="speedDownloadTask.status">{{ speedDownloadStatusText }}</span>
              </div>
              <div class="sdp-stat-sub">
                <span>{{ formatFileSize(speedDownloadTask.downloadedBytes) }} / {{ formatFileSize(speedDownloadTask.fileSize) }}</span>
                <span class="sdp-percent">{{ speedDownloadTask.progress }}%</span>
              </div>
            </div>

            <!-- 主进度条 -->
            <div class="sdp-progress-track">
              <div class="sdp-progress-bar" :style="{ width: speedDownloadTask.progress + '%' }"></div>
            </div>

            <!-- 6 条多协程并发通道进度跑马灯 -->
            <div class="sdp-channels" v-if="speedDownloadTask.chunks && speedDownloadTask.chunks.length">
              <div v-for="ch in speedDownloadTask.chunks" :key="ch.index" class="sdp-channel-item" :title="'并发线程 ' + (ch.index + 1) + ': ' + ch.progress + '%'">
                <div class="sdp-ch-fill" :style="{ width: ch.progress + '%' }"></div>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </main>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, triggerRef } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { checkHash, instantUpload, initUpload, uploadChunkWithProgress, completeUpload, getFileList, deleteFile, getDownloadUrl, getPickupUrl, getUserList, updateUserStatus, updateUserSpace, getUserInfo, updateUserPassword, giftSpace, createTicket, getTicketList, replyTicket, getMessageList, readMessage, readAllMessage, sendMessage, getUploadTasks, getProgress, deleteUploadTask } from '../api'
import { calculateFileMD5, formatFileSize, formatSpeed } from '../utils'

const emit = defineEmits(['logout'])

// 极速多线程并发切片下载状态（文叔叔式并发加速）
const speedDownloadTask = reactive({
  visible: false,
  fileName: '',
  fileSize: 0,
  downloadedBytes: 0,
  progress: 0,
  speed: 0,
  status: 'idle', // 'idle' | 'downloading' | 'merging' | 'completed' | 'error'
  chunks: []
})

const speedDownloadStatusText = computed(() => {
  switch (speedDownloadTask.status) {
    case 'downloading': return '高速多协程下载中'
    case 'merging': return '合并写入磁盘中'
    case 'completed': return '下载完成'
    case 'error': return '下载出错'
    default: return '准备中'
  }
})

// 获取本地缓存并过滤非法值（'undefined', 'null' 等）
const getStoredUsername = () => {
  const stored = localStorage.getItem('username')
  if (!stored || stored === 'undefined' || stored === 'null') return ''
  return stored
}

const username = ref(getStoredUsername() || 'admin')
const currentUser = ref(null)

// 动态展示名称与权限
const displayName = computed(() => {
  return currentUser.value?.username || (username.value && username.value !== 'undefined' ? username.value : 'admin')
})

const avatarChar = computed(() => {
  const name = displayName.value || 'A'
  return name.charAt(0).toUpperCase()
})

const isAdmin = computed(() => {
  return displayName.value === 'admin' || currentUser.value?.username === 'admin' || username.value === 'admin'
})

const currentTab = ref('files')
const userListData = ref([])
const userListLoading = ref(false)

// 空间赠送状态
const showGiftForm = ref(false)
const submittingGift = ref(false)
const giftForm = reactive({
  toUsername: '',
  giftSizeMb: 1024,
  verifyPassword: ''
})

async function handleGiftSpace() {
  if (!giftForm.toUsername.trim()) {
    return ElMessage.warning('请输入接收方用户名')
  }
  if (!giftForm.giftSizeMb || giftForm.giftSizeMb <= 0) {
    return ElMessage.warning('赠送容量必须大于 0 MB')
  }
  if (!giftForm.verifyPassword) {
    return ElMessage.warning('请输入当前登录密码进行安全验证')
  }

  submittingGift.value = true
  try {
    const res = await giftSpace({
      toUsername: giftForm.toUsername.trim(),
      giftSizeMb: Number(giftForm.giftSizeMb),
      verifyPassword: giftForm.verifyPassword
    })
    if (res.code === 200) {
      ElMessage.success(`成功赠送 ${giftForm.giftSizeMb}MB 空间给 ${giftForm.toUsername}！`)
      showGiftForm.value = false
      giftForm.toUsername = ''
      giftForm.verifyPassword = ''
      loadUserInfo()
    } else {
      ElMessage.error(res.msg || '空间赠送失败')
    }
  } catch (e) {
    ElMessage.error('赠送请求失败，请检查网络')
  } finally {
    submittingGift.value = false
  }
}

async function loadUsers() {
  if (!isAdmin.value) return
  userListLoading.value = true
  try {
    const res = await getUserList()
    if (res.code === 200) userListData.value = res.data || []
  } catch (e) {
    ElMessage.error('获取用户列表失败')
  } finally {
    userListLoading.value = false
  }
}

async function handleToggleStatus(user) {
  const newStatus = user.status === 1 ? 0 : 1
  try {
    const res = await updateUserStatus({ id: user.id, status: newStatus })
    if (res.code === 200) {
      user.status = newStatus; ElMessage.success(newStatus === 1 ? '已解封该账户' : '已冻结该账户')
    } else ElMessage.error(res.msg || '操作失败')
  } catch(e) {}
}

async function handleEditSpace(user) {
  ElMessageBox.prompt(`请输入为用户【${user.username}】设置的总空间配额（单位：GB）：`, '设置空间配额', {
    confirmButtonText: '确定下发', cancelButtonText: '取消',
    inputValue: (user.totalSpace / (1024*1024*1024)).toFixed(1),
    inputPattern: /^\d+(\.\d+)?$/, inputErrorMessage: '容量需为非负数字'
  }).then(async ({ value }) => {
    const totalSpace = Math.floor(parseFloat(value) * 1024 * 1024 * 1024)
    try {
      const res = await updateUserSpace({ id: user.id, totalSpace })
      if (res.code === 200) {
        user.totalSpace = totalSpace
        ElMessage.success(`用户【${user.username}】空间配额已成功调整为 ${value} GB`)
        loadUsers()
        loadUserInfo()
      } else {
        ElMessage.error(res.msg || '空间设置失败')
      }
    } catch(e) {
      ElMessage.error('网络异常，空间配额设置失败')
    }
  }).catch(() => {})
}

async function handleResetPassword(user) {
  ElMessageBox.prompt(`请输入为用户 ${user.username} 设置的新密码：`, '重置密码', {
    confirmButtonText: '确定重置', cancelButtonText: '取消',
    inputType: 'password',
    inputPattern: /^\S+$/, inputErrorMessage: '密码不能为空或包含空格'
  }).then(async ({ value }) => {
    try {
      const res = await updateUserPassword({ id: user.id, password: value })
      if (res.code === 200) {
        ElMessage.success('密码重置成功')
      } else {
        ElMessage.error(res.msg || '密码重置失败')
      }
    } catch(e) {
      ElMessage.error('网络请求失败')
    }
  }).catch(() => {})
}

// ==== 工单支持 ====
const ticketList = ref([])
const ticketLoading = ref(false)
const showTicketForm = ref(false)
const submittingTicket = ref(false)
const newTicket = reactive({ title: '', content: '' })

// ==== 系统通知管理 ====
const unreadMsgCount = ref(0)
const msgList = ref([])
const msgLoading = ref(false)
const showMsgForm = ref(false)
const msgTarget = ref('all')
const submittingMsg = ref(false)
const newMsg = reactive({ title: '', content: '' })

async function loadMessages() {
  msgLoading.value = true
  try {
    const res = await getMessageList()
    if (res.code === 200) {
      msgList.value = res.data.list || []
      unreadMsgCount.value = res.data.unreadCount || 0
    }
  } catch(e) {
    // console.error(e)
  } finally {
    msgLoading.value = false
  }
}

async function handleReadMsg(msg) {
  if (msg.isRead === 1) return
  msg.isRead = 1
  unreadMsgCount.value = Math.max(0, unreadMsgCount.value - 1)
  try {
    await readMessage(msg.id)
  } catch(e) {}
}

async function handleReadAllMsg() {
  try {
    await readAllMessage()
    msgList.value.forEach(m => m.isRead = 1)
    unreadMsgCount.value = 0
    ElMessage.success('已清空未读状态')
  } catch(e) {}
}

function handleOpSendMsg(targetId) {
  msgTarget.value = targetId
  newMsg.title = ''
  newMsg.content = ''
  showMsgForm.value = true
}

async function handleSubmitMsg() {
  if (!newMsg.title.trim() || !newMsg.content.trim()) {
    return ElMessage.warning('标题和细节必须填写')
  }
  submittingMsg.value = true
  try {
    const res = await sendMessage({ targetId: msgTarget.value, title: newMsg.title, content: newMsg.content })
    if (res.code === 200) {
      ElMessage.success('发布通知成功')
      showMsgForm.value = false
      if (msgTarget.value !== 'all' || username.value !== 'admin') {
         // 可选刷新
      }
    } else {
      ElMessage.error(res.msg || '发送通信失败')
    }
  } catch(e) {
    ElMessage.error('无法送达网络')
  } finally {
    submittingMsg.value = false
  }
}

async function loadTickets() {
  ticketLoading.value = true
  try {
    const res = await getTicketList()
    if (res.code === 200) ticketList.value = res.data || []
  } catch (e) {
    ElMessage.error('获取工单列表失败')
  } finally {
    ticketLoading.value = false
  }
}

async function handleSubmitTicket() {
  if (!newTicket.title.trim() || !newTicket.content.trim()) {
    return ElMessage.warning('标题和内容均不能为空')
  }
  submittingTicket.value = true
  try {
    const res = await createTicket({ title: newTicket.title, content: newTicket.content })
    if (res.code === 200) {
      ElMessage.success('工单提交成功')
      showTicketForm.value = false
      newTicket.title = ''
      newTicket.content = ''
      loadTickets()
    } else {
      ElMessage.error(res.msg || '提交通信失败')
    }
  } catch (e) {
    ElMessage.error('网络失败，请重试')
  } finally {
    submittingTicket.value = false
  }
}

async function handleOpReply(ticket) {
  ElMessageBox.prompt(`请输入对待处理工单 #${ticket.id} (${ticket.title}) 的回复：`, '管理员处理', {
    confirmButtonText: '确定回复',
    cancelButtonText: '取消',
    inputType: 'textarea',
    inputPattern: /\S+/,
    inputErrorMessage: '回复内容不能为空'
  }).then(async ({ value }) => {
    try {
      const res = await replyTicket({ id: ticket.id, replyContent: value })
      if (res.code === 200) {
        ElMessage.success('回复成功')
        loadTickets()
      } else {
        ElMessage.error(res.msg || '通信失败')
      }
    } catch(e) {
      ElMessage.error('网络错误')
    }
  }).catch(() => {})
}

const isDragging = ref(false)
let dragCounter = 0
const pickupCodeInput = ref('')
const selectedExpireDays = ref(1)
const expireOptions = [
  { label: '1天后', value: 1 },
  { label: '5天后', value: 5 },
  { label: '7天后', value: 7 },
  { label: '30天后', value: 30 },
  { label: '永久有效', value: -1 }
]
const fileInput = ref(null)
const uploadTasks = ref([])
let taskIdCounter = 0
const fileList = ref([])
const listLoading = ref(false)
const searchKeyword = ref('')
const pageNum = ref(1)
const pageSize = ref(20)
const total = ref(0)

async function loadUserInfo() {
  try {
    const res = await getUserInfo()
    if (res.code === 200 && res.data) {
      currentUser.value = res.data
      if (res.data.username) {
        username.value = res.data.username
        localStorage.setItem('username', res.data.username)
      }
    }
  } catch (e) {
    console.warn('获取用户信息失败', e)
  }
}

async function loadUploadTasks() {
  try {
    const res = await getUploadTasks()
    if (res.code === 200 && res.data) {
      res.data.forEach(item => {
        // 如果当前列表里没有这个 uploadId
        if (!uploadTasks.value.find(t => t.uploadId === item.uploadId)) {
          uploadTasks.value.push(reactive({
            id: ++taskIdCounter,
            uploadId: item.uploadId,
            fileName: item.fileName,
            fileSize: item.fileSize,
            fileHash: item.fileHash,
            status: item.status, // 后端返回 interrupted 或 merging
            progress: item.progress,
            uploadedBytes: 0,
            file: null, // 物理文件句柄已丢，需要重新选择
            startTime: null,
            elapsedSeconds: 0,
            remainingSeconds: null,
            timer: null
          }))
        }
      })
    } else if (res.code !== 200) {
      console.error('获取上传任务失败:', res.msg)
    }
  } catch (e) {
    console.error('加载上传任务网络异常', e)
  }
}

onMounted(async () => { 
  loadFileList()
  await loadUserInfo() // 核心：确保 userId 已准备好（拦截器可能需要）
  loadMessages()
  loadUploadTasks() 
})

onUnmounted(() => {
  uploadTasks.value.forEach(t => {
    if (t.timer) clearInterval(t.timer)
  })
})

function handleFileSelect(e) {
  const files = e.target.files
  if (files) { for (const f of files) queueUpload(f) }
  e.target.value = ''
}

function handleGlobalDragOver(e) {
  e.dataTransfer.dropEffect = 'copy'
  if (!isDragging.value) {
    dragCounter++
    isDragging.value = true
  }
}

function handleGlobalDragLeave(e) {
  dragCounter--
  if (dragCounter <= 0) {
    dragCounter = 0
    isDragging.value = false
  }
}

function handleGlobalDrop(e) {
  dragCounter = 0
  isDragging.value = false
  const files = e.dataTransfer.files
  if (files) { for (const f of files) queueUpload(f) }
}

async function copyCode(code) {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(code)
    } else {
      const input = document.createElement('input');
      input.value = code;
      document.body.appendChild(input);
      input.select();
      document.execCommand('copy');
      document.body.removeChild(input);
    }
    ElMessage.success('取件码已复制: ' + code)
  } catch (e) {
    ElMessage.error('复制失败: ' + (e.message || ''))
  }
}

function handlePickup() {
  const code = pickupCodeInput.value.trim()
  if (!code) return ElMessage.warning('请输入取件码')
  if (code.length !== 6) return ElMessage.warning('取件码应该是6位字符')
  
  // 对于提取功能，不需要调用额外 API 查询，可以直接拉取流。为了优雅处理错误可以在打开前校验
  // 这里做一个简单的 iframe 触发下载，避免离开当前页
  const a = document.createElement('a')
  a.href = getPickupUrl(code)
  a.click()
  pickupCodeInput.value = ''
}

// 强制触发视图更新
function forceUpdate() {
  triggerRef(uploadTasks)
}

// 并发限制与排队配置（每个用户最多同时执行 5 个上传任务）
const MAX_CONCURRENT_TASKS = 5

function formatDuration(sec) {
  if (!sec || sec < 0) return '0秒'
  const s = Math.floor(sec)
  if (s < 60) return `${s}秒`
  const m = Math.floor(s / 60)
  const remSec = s % 60
  if (m < 60) return `${m}分${remSec > 0 ? remSec + '秒' : ''}`
  const h = Math.floor(m / 60)
  const remMin = m % 60
  return `${h}小时${remMin > 0 ? remMin + '分' : ''}`
}

function queueUpload(file) {
  // 如果当前列表有已被中断或正在合并的该任务（通过文件名+大小判断），则复用以防重复
  let task = uploadTasks.value.find(t => t.fileName === file.name && Math.abs(t.fileSize - file.size) < 2 && (t.status === 'interrupted' || t.status === 'merging'))
  
  if (task) {
    task.file = file // 必须补回物理文件句柄
    task.status = 'queued'
    task.elapsedSeconds = 0
    task.remainingSeconds = null
    task.speed = 0
  } else {
    task = reactive({
      id: ++taskIdCounter,
      fileName: file.name,
      fileSize: file.size,
      file,
      status: 'queued',
      progress: 0,
      speed: 0,
      uploadId: null,
      uploadedBytes: 0,
      startTime: null,
      elapsedSeconds: 0,
      remainingSeconds: null,
      timer: null
    })
    uploadTasks.value.unshift(task)
  }

  forceUpdate()
  processUploadQueue()
}

function processUploadQueue() {
  // 当前正在活动的任务数 (hashing, uploading, merging)
  const activeTasks = uploadTasks.value.filter(t => ['hashing', 'uploading', 'merging'].includes(t.status))
  const activeCount = activeTasks.length

  if (activeCount >= MAX_CONCURRENT_TASKS) {
    return
  }

  const availableSlots = MAX_CONCURRENT_TASKS - activeCount
  // 找出所有排队中的任务，按添加顺序优先执行 (FIFO)
  const queuedTasks = uploadTasks.value.filter(t => t.status === 'queued')
  queuedTasks.sort((a, b) => a.id - b.id)

  const tasksToStart = queuedTasks.slice(0, availableSlots)
  for (const t of tasksToStart) {
    executeUploadTask(t)
  }
}

async function executeUploadTask(task) {
  if (task.status !== 'queued') return

  task.status = 'hashing'
  task.startTime = Date.now()
  task.elapsedSeconds = 0
  task.remainingSeconds = null

  // 启动动态用时定时器 (每秒刷新一次用时及剩余时间)
  if (task.timer) clearInterval(task.timer)
  task.timer = setInterval(() => {
    if (['hashing', 'uploading', 'merging'].includes(task.status)) {
      task.elapsedSeconds = Math.max(1, Math.floor((Date.now() - task.startTime) / 1000))
      // 动态预估剩余时间
      if (task.status === 'uploading' && task.speed > 0 && task.fileSize > task.uploadedBytes) {
        task.remainingSeconds = Math.max(0, Math.round((task.fileSize - task.uploadedBytes) / task.speed))
      } else {
        task.remainingSeconds = null
      }
      forceUpdate()
    }
  }, 1000)

  const stopTaskTimer = () => {
    if (task.timer) {
      clearInterval(task.timer)
      task.timer = null
    }
    if (task.startTime) {
      task.elapsedSeconds = Math.max(1, Math.floor((Date.now() - task.startTime) / 1000))
    }
  }

  try {
    forceUpdate()

    // 1. Hash 计算阶段（hash阶段占进度 0-5%）
    const fileHash = await calculateFileMD5(task.file, (p) => {
      if (!task.uploadId) {
        task.progress = Math.round(p * 0.05)
      }
      task.status = 'hashing'
      forceUpdate()
    })

    // 2. 秒传检查
    const checkRes = await checkHash(fileHash)
    if (checkRes.code === 200 && checkRes.data.exists) {
      task.status = 'instant'
      task.progress = 50
      forceUpdate()

      try {
        const instantRes = await instantUpload({
          fileHash,
          fileName: task.file.name,
          fileSize: task.file.size,
          expireDays: selectedExpireDays.value
        })

        if (instantRes.code !== 200) {
          throw new Error(instantRes.msg || '秒传失败')
        }

        task.progress = 100
        stopTaskTimer()
        forceUpdate()
        ElMessage.success(`${task.fileName} 秒传成功! (无需上传任何数据)`)
        pageNum.value = 1
        setTimeout(() => { loadFileList(); loadUserInfo(); }, 300)
        setTimeout(() => {
          uploadTasks.value = uploadTasks.value.filter(t => t.id !== task.id)
          forceUpdate()
          processUploadQueue()
        }, 3000)
        processUploadQueue()
        return
      } catch (instantErr) {
        console.warn('秒传失败，降级为普通上传:', instantErr)
        ElMessage.warning(`秒传失败: ${instantErr.message}，切换为普通上传...`)
        task.status = 'uploading'
        task.progress = 5
        forceUpdate()
      }
    }

    // 3. 初始化分片上传 (后端如果超过5个并发会返回友好提示，在此双重保障)
    const initRes = await initUpload({ fileName: task.file.name, fileSize: task.file.size, fileHash })
    if (initRes.code !== 200) throw new Error(initRes.msg)
    const { uploadId, chunkSize, totalChunk } = initRes.data
    task.uploadId = uploadId
    task.status = 'uploading'
    forceUpdate()

    // 4. 断点续传核心获取真实进度
    let uploadedChunks = []
    try {
      const progRes = await getProgress(uploadId)
      if (progRes.code === 200 && progRes.data.uploadedChunks) {
        uploadedChunks = progRes.data.uploadedChunks
      }
    } catch(e) { console.warn('获取断点进度失败', e) }

    // 5. 并行分片上传
    const concurrency = 6
    let completedChunks = uploadedChunks.length
    let initialBytes = Math.min(completedChunks * chunkSize, task.file.size)
    task.uploadedBytes = initialBytes
    task.progress = Math.round((initialBytes / task.file.size) * 100)
    forceUpdate()

    if (completedChunks >= totalChunk) {
      console.log('所有分片已存在，准备开始合并')
    } else {
      const chunkProgressMap = {}
      uploadedChunks.forEach(ci => { 
        const start = ci * chunkSize
        const end = Math.min(start + chunkSize, task.file.size)
        chunkProgressMap[ci] = end - start 
      })

      const uploadQueue = []
      for (let i = 0; i < totalChunk; i++) {
        if (!uploadedChunks.includes(i)) uploadQueue.push(i)
      }

      const workers = []
      for (let w = 0; w < concurrency; w++) {
        workers.push((async () => {
          while (uploadQueue.length > 0) {
            const ci = uploadQueue.shift()
            if (ci === undefined) break
            const start = ci * chunkSize
            const end = Math.min(start + chunkSize, task.file.size)
            const chunk = task.file.slice(start, end)

            chunkProgressMap[ci] = 0

            const res = await uploadChunkWithProgress(uploadId, ci, chunk, (loaded) => {
              chunkProgressMap[ci] = loaded
              const currentUploaded = Object.values(chunkProgressMap).reduce((a, b) => a + b, 0)
              const elapsed = (Date.now() - task.startTime) / 1000
              task.uploadedBytes = currentUploaded
              task.speed = elapsed > 0 ? (currentUploaded - initialBytes) / elapsed : 0
              task.progress = Math.min(95, Math.round((currentUploaded / task.file.size) * 100))
              if (task.speed > 0) {
                task.remainingSeconds = Math.max(0, Math.round((task.fileSize - currentUploaded) / task.speed))
              }
              forceUpdate()
            })

            if (res.code !== 200) throw new Error(`分片 ${ci} 失败`)
            completedChunks++
            chunkProgressMap[ci] = (end - start)
          }
        })())
      }
      await Promise.all(workers)
    }

    // 6. 合并阶段
    task.status = 'merging'
    task.progress = 96
    task.speed = 0
    task.remainingSeconds = null
    forceUpdate()

    const completeRes = await completeUpload(uploadId, selectedExpireDays.value)
    if (completeRes.code !== 200) throw new Error(completeRes.msg)

    // 7. 完成阶段
    task.status = 'done'
    task.progress = 100
    task.speed = 0
    stopTaskTimer()
    forceUpdate()
    ElMessage.success(`${task.fileName} 上传完成! 总耗时: ${formatDuration(task.elapsedSeconds)}`)
    pageNum.value = 1
    setTimeout(() => { loadFileList(); loadUserInfo(); }, 300)
    setTimeout(() => {
      uploadTasks.value = uploadTasks.value.filter(t => t.id !== task.id)
      forceUpdate()
      processUploadQueue()
    }, 3000)
  } catch (e) {
    stopTaskTimer()
    task.status = 'error'
    task.speed = 0
    forceUpdate()
    ElMessage.error(`${task.fileName} 失败: ${e.message}`)
  } finally {
    processUploadQueue()
  }
}

async function loadFileList() {
  listLoading.value = true
  try {
    const res = await getFileList({ keyword: searchKeyword.value, pageNum: pageNum.value, pageSize: pageSize.value })
    if (res.code === 200) { fileList.value = res.data.list || []; total.value = res.data.total || 0 }
  } catch (e) { console.error(e) } finally { listLoading.value = false }
}

function handleNormalDownload(row) {
  const a = document.createElement('a')
  a.href = getDownloadUrl(row.id)
  a.download = row.fileName
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

function handleDownload(row) {
  handleNormalDownload(row)
}

function copyDownloadCommand(row) {
  const downloadUrl = window.location.origin + getDownloadUrl(row.id)
  const safeName = (row.fileName || 'downloaded_file').replace(/"/g, '\\"')
  const aria2Cmd = `aria2c -s 16 -x 16 -k 1M -o "${safeName}" "${downloadUrl}"`
  const curlCmd = `curl -L -C - -o "${safeName}" "${downloadUrl}"`
  const text = `# 极速 16 协程分片并发下载 (需要 aria2):\n${aria2Cmd}\n\n# 或断点续传普通下载 (curl):\n${curlCmd}`

  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success('已复制终端极速下载命令（含 aria2 16协程 与 curl）')
    }).catch(() => {
      promptCopyCommand(text)
    })
  } else {
    promptCopyCommand(text)
  }
}

function promptCopyCommand(text) {
  ElMessageBox.alert(
    `<textarea style="width:100%;height:140px;font-family:monospace;font-size:12px;background:#0f172a;color:#38bdf8;padding:10px;border-radius:8px;border:1px solid rgba(255,255,255,0.1);outline:none;resize:none;" readonly>${text}</textarea>`,
    '终端极速下载命令',
    { dangerouslyUseHTMLString: true, confirmButtonText: '确定' }
  )
}

// 🚀 文叔叔式并发加速：多协程 Range 切片下载引擎
async function handleSpeedDownload(row) {
  const downloadUrl = getDownloadUrl(row.id)
  let totalSize = row.fileSize || 0

  // 1. 如果无大小，先用 HEAD 探测 Content-Length
  if (!totalSize) {
    try {
      const headRes = await fetch(downloadUrl, { method: 'HEAD' })
      const cl = headRes.headers.get('content-length')
      if (cl) totalSize = parseInt(cl, 10)
    } catch (e) {
      console.warn('HEAD probe failed:', e)
    }
  }

  // 小于 2MB 或无法探测大小，直接走普通下载
  if (!totalSize || totalSize < 2 * 1024 * 1024) {
    ElMessage.info('文件较小，正在通过高速通道直接单流下载')
    handleNormalDownload(row)
    return
  }

  // 2. 尝试唤起现代浏览器 File System Access API 直写落盘（零内存开销，文叔叔同款黑科技）
  let fileHandle = null
  let writableStream = null
  let useFileSystem = false

  if (typeof window.showSaveFilePicker === 'function') {
    try {
      fileHandle = await window.showSaveFilePicker({
        suggestedName: row.fileName
      })
      writableStream = await fileHandle.createWritable()
      useFileSystem = true
    } catch (err) {
      if (err.name === 'AbortError') {
        // 用户主动取消
        return
      }
      console.warn('showSaveFilePicker unsupported or cancelled, fallback to memory aggregation:', err)
    }
  }

  // 3. 规划 6 个并发切片通道
  const CONCURRENCY = 6
  const chunkSize = Math.ceil(totalSize / CONCURRENCY)
  const chunks = []
  for (let i = 0; i < CONCURRENCY; i++) {
    const start = i * chunkSize
    if (start >= totalSize) break
    const end = Math.min(totalSize - 1, (i + 1) * chunkSize - 1)
    chunks.push({
      index: i,
      start,
      end,
      total: end - start + 1,
      downloaded: 0,
      progress: 0
    })
  }

  // 初始化下载看板
  speedDownloadTask.visible = true
  speedDownloadTask.fileName = row.fileName
  speedDownloadTask.fileSize = totalSize
  speedDownloadTask.downloadedBytes = 0
  speedDownloadTask.progress = 0
  speedDownloadTask.speed = 0
  speedDownloadTask.status = 'downloading'
  speedDownloadTask.chunks = chunks

  // 4. 实时速率测速器 (400ms 刷新一次)
  let lastLoaded = 0
  let lastTime = performance.now()
  const speedTimer = setInterval(() => {
    const now = performance.now()
    const duration = (now - lastTime) / 1000
    if (duration >= 0.3) {
      const delta = speedDownloadTask.downloadedBytes - lastLoaded
      speedDownloadTask.speed = Math.max(0, Math.round(delta / duration))
      lastLoaded = speedDownloadTask.downloadedBytes
      lastTime = now
    }
  }, 400)

  // 5. 并发切片下载工作线程
  const chunkBuffers = useFileSystem ? null : new Array(chunks.length)

  try {
    const downloadWorkers = chunks.map(async (chunk) => {
      const resp = await fetch(downloadUrl, {
        headers: {
          'Range': `bytes=${chunk.start}-${chunk.end}`
        }
      })

      if (resp.status !== 206 && resp.status !== 200) {
        throw new Error(`切片 ${chunk.index} HTTP 状态码异常: ${resp.status}`)
      }

      const reader = resp.body.getReader()
      const parts = []

      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        parts.push(value)
        chunk.downloaded += value.length
        speedDownloadTask.downloadedBytes += value.length
        chunk.progress = Math.min(100, Math.round((chunk.downloaded / chunk.total) * 100))
        speedDownloadTask.progress = Math.min(100, Math.round((speedDownloadTask.downloadedBytes / speedDownloadTask.fileSize) * 100))
      }

      if (useFileSystem && writableStream) {
        // 直接 seek 到切片起始偏移写入磁盘
        const chunkBlob = new Blob(parts)
        await writableStream.seek(chunk.start)
        await writableStream.write(chunkBlob)
      } else {
        chunkBuffers[chunk.index] = parts
      }
    })

    // 并发启动 6 条并发连接
    await Promise.all(downloadWorkers)

    clearInterval(speedTimer)
    speedDownloadTask.speed = 0
    speedDownloadTask.progress = 100

    // 6. 完成组装或关闭流
    if (useFileSystem && writableStream) {
      speedDownloadTask.status = 'merging'
      await writableStream.close()
      speedDownloadTask.status = 'completed'
      ElMessage.success('🎉 极速并发下载完成，已直写保存至磁盘！')
    } else {
      speedDownloadTask.status = 'merging'
      const flatParts = []
      for (let i = 0; i < chunkBuffers.length; i++) {
        if (chunkBuffers[i]) {
          flatParts.push(...chunkBuffers[i])
        }
      }
      const fullBlob = new Blob(flatParts, { type: 'application/octet-stream' })
      const blobUrl = URL.createObjectURL(fullBlob)
      const a = document.createElement('a')
      a.href = blobUrl
      a.download = row.fileName
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(blobUrl)
      speedDownloadTask.status = 'completed'
      ElMessage.success('🎉 极速并发下载完成，正在保存文件！')
    }

    // 4秒后自动收起看板
    setTimeout(() => {
      if (speedDownloadTask.status === 'completed') {
        speedDownloadTask.visible = false
      }
    }, 4000)

  } catch (err) {
    clearInterval(speedTimer)
    console.error('Speed download error:', err)
    speedDownloadTask.status = 'error'
    ElMessage.error(`极速下载遇到异常: ${err.message || '网络连接中断'}，为您自动降级至普通单流下载`)
    handleNormalDownload(row)
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(`确定删除 "${row.fileName}" 吗？`, '删除确认', { type: 'warning' })
    const res = await deleteFile(row.id)
    if (res.code === 200) { ElMessage.success('删除成功'); loadFileList(); loadUserInfo(); } else { ElMessage.error(res.msg) }
  } catch (e) { /* cancelled */ }
}

// 🔥 删除已中断或失败、取消排队的上传任务
async function handleDeleteTask(task) {
  try {
    if (task.timer) {
      clearInterval(task.timer)
      task.timer = null
    }

    if (task.status === 'queued') {
      uploadTasks.value = uploadTasks.value.filter(t => t.id !== task.id)
      forceUpdate()
      processUploadQueue()
      ElMessage.success('已取消排队')
      return
    }

    await ElMessageBox.confirm(
      `确定删除上传任务 "${task.fileName}" 吗？\n已上传的分片将被清理。`,
      '删除确认',
      { type: 'warning' }
    )

    if (task.uploadId) {
      const res = await deleteUploadTask(task.uploadId)
      if (res.code !== 200) {
        ElMessage.error(res.msg || '删除失败')
        return
      }
    }
    ElMessage.success('任务已删除')
    uploadTasks.value = uploadTasks.value.filter(t => t.id !== task.id)
    forceUpdate()
    processUploadQueue()
  } catch (e) {
    if (e !== 'cancel') {
      console.error('删除任务失败:', e)
      ElMessage.error('删除失败')
    }
  }
}

function handleLogout() {
  localStorage.removeItem('token'); localStorage.removeItem('username')
  emit('logout')
}

function formatDate(str) {
  if (!str) return '-'; const d = new Date(str); const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function statusText(s) {
  return { 
    queued: '排队中',
    hashing: '计算Hash', 
    uploading: '上传中', 
    merging: '合并中', 
    done: '完成', 
    error: '失败', 
    instant: '秒传', 
    interrupted: '已中断' 
  }[s] || s
}

function getExt(name) {
  const idx = name.lastIndexOf('.')
  return idx >= 0 ? name.substring(idx + 1).toUpperCase().slice(0, 4) : 'FILE'
}
</script>

<style scoped>
/* ===== 布局 ===== */
.app-layout {
  display: flex;
  min-height: 100vh;
  background: #060b1a;
  font-family: 'Open Sans', sans-serif;
}

/* ===== 侧边栏 ===== */
.sidebar {
  width: 240px;
  background: rgba(10, 15, 30, 0.9);
  border-right: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  flex-direction: column;
  padding: 24px 16px;
  flex-shrink: 0;
  backdrop-filter: blur(20px);
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 8px;
  margin-bottom: 36px;
}

.brand-icon {
  width: 40px; height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #2563EB, #7c3aed);
  border-radius: 12px;
  color: white;
  flex-shrink: 0;
}

.brand-name {
  font-family: 'Poppins', sans-serif;
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.3px;
}

.sidebar-nav {
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.45);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease-out;
  cursor: pointer;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.7);
}

.nav-item.active {
  background: rgba(37, 99, 235, 0.15);
  color: #60a5fa;
}

.sidebar-footer {
  display: flex;
  flex-direction: column;
  padding: 16px 8px 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  margin-top: auto;
}
.footer-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}
.user-storage-panel {
  margin-top: 16px;
  width: 100%;
  padding: 0 4px;
}
.us-text {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 8px;
}
.us-percent {
  font-weight: 600;
  color: #60a5fa;
}
.us-bar-bg {
  height: 4px;
  width: 100%;
  background: rgba(255,255,255,0.08);
  border-radius: 4px;
  overflow: hidden;
}
.us-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #2563EB, #7c3aed);
  border-radius: 4px;
  transition: width 0.4s ease;
}
.us-bar-fill.danger {
  background: #ef4444;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 34px; height: 34px;
  border-radius: 10px;
  background: linear-gradient(135deg, #2563EB, #7c3aed);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 14px;
}

.user-meta {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 13px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
}

.user-role {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.3);
}

.logout-btn {
  background: none; border: none; padding: 8px;
  border-radius: 10px; color: rgba(255, 255, 255, 0.3);
  cursor: pointer; transition: all 0.2s ease-out;
  display: flex; align-items: center;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

/* ===== 主内容区 ===== */
.main-area {
  flex: 1;
  padding: 28px 36px;
  overflow-y: auto;
  max-height: 100vh;
}

/* 顶部 */
.top-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 28px;
}

.page-title {
  font-family: 'Poppins', sans-serif;
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 4px;
}

.page-desc {
  color: rgba(255, 255, 255, 0.35);
  font-size: 14px;
  margin: 0;
}

.header-left {
  display: flex;
  flex-direction: column;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.pickup-wrap {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 20px;
  overflow: hidden;
  height: 40px;
}

.pickup-input {
  background: transparent;
  border: none;
  padding: 0 16px;
  color: #fff;
  font-size: 14px;
  outline: none;
  width: 130px;
}

.pickup-input::placeholder {
  color: rgba(255, 255, 255, 0.3);
}

.pickup-btn {
  background: rgba(37, 99, 235, 0.1);
  color: #60a5fa;
  border: none;
  cursor: pointer;
  padding: 0 16px;
  height: 100%;
  font-weight: 600;
  transition: all 0.2s;
}

.pickup-btn:hover {
  background: #2563eb;
  color: #fff;
}

.upload-action-group {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 20px;
  padding: 4px;
  padding-left: 12px;
  height: 40px;
}

.expire-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  margin-right: 6px;
}

.expire-select {
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.8);
  font-size: 13px;
  outline: none;
  cursor: pointer;
  margin-right: 12px;
}

.expire-select option {
  background: #0a0f1e;
  color: #fff;
}

.upload-btn-primary {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: 16px;
  padding: 0 16px;
  height: 32px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-btn-primary:hover {
  background: #3b82f6;
  transform: translateY(-1px);
}

.stats-row {
  display: flex;
  gap: 12px;
}

.stat-card {
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 90px;
}

.stat-value {
  font-family: 'Poppins', sans-serif;
  font-size: 22px;
  font-weight: 700;
  color: #60a5fa;
}

.stat-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
  margin-top: 2px;
}

/* ===== 全局拖拽遮罩 ===== */
.global-dropzone {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(6, 11, 26, 0.85);
  backdrop-filter: blur(8px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4px dashed rgba(59, 130, 246, 0.5);
  border-radius: 12px;
  pointer-events: none; /* 让事件穿透，保持对父容器 drag 监听 */
}

.dropzone-content {
  text-align: center;
  color: #3b82f6;
  animation: pulse 2s infinite;
}

.dropzone-content p {
  color: white;
  font-size: 24px;
  font-weight: 600;
  margin-top: 16px;
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 0.8; }
  50% { transform: scale(1.05); opacity: 1; }
  100% { transform: scale(1); opacity: 0.8; }
}

/* ===== 提取码标记 ===== */
.pickup-badge {
  display: inline-block;
  padding: 4px 8px;
  background: rgba(37, 99, 235, 0.15);
  color: #60a5fa;
  border-radius: 6px;
  font-family: monospace;
  font-size: 13px;
  letter-spacing: 1px;
  cursor: pointer;
  border: 1px solid rgba(59, 130, 246, 0.2);
  transition: all 0.2s;
}

.pickup-badge:hover {
  background: rgba(37, 99, 235, 0.3);
  border-color: rgba(59, 130, 246, 0.5);
}

/* ===== 任务区 ===== */
.task-section, .file-section {
  margin-bottom: 28px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'Poppins', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 16px;
}

.section-title svg {
  color: #60a5fa;
}

.task-grid {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.task-card {
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  transition: border-color 0.2s;
}

.task-card.status-uploading { border-left: 3px solid #3b82f6; }
.task-card.status-done, .task-card.status-instant { border-left: 3px solid #22c55e; }
.task-card.status-error { border-left: 3px solid #ef4444; }
.task-card.status-hashing { border-left: 3px solid #f59e0b; }
.task-card.status-merging { border-left: 3px solid #a855f7; }

.task-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.task-file-icon {
  width: 36px; height: 36px;
  border-radius: 10px;
  background: rgba(59, 130, 246, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #60a5fa;
  flex-shrink: 0;
}

.task-meta {
  flex: 1;
  min-width: 0;
}

.task-name {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-size {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}

.task-badge {
  padding: 4px 10px;
  border-radius: 100px;
  font-size: 11px;
  font-weight: 600;
  white-space: nowrap;
  letter-spacing: 0.3px;
}

.badge-queued { background: rgba(148, 163, 184, 0.15); color: #94a3b8; }
.badge-hashing { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }
.badge-uploading { background: rgba(59, 130, 246, 0.15); color: #60a5fa; }
.badge-merging { background: rgba(168, 85, 247, 0.15); color: #a855f7; }
.badge-done, .badge-instant { background: rgba(34, 197, 94, 0.15); color: #22c55e; }
.badge-error { background: rgba(239, 68, 68, 0.15); color: #ef4444; }
.badge-interrupted { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }

.task-progress-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 100px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #2563EB, #3b82f6);
  border-radius: 100px;
  transition: width 0.3s ease-out;
}

.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}

/* 🔥 删除按钮样式 */
.task-delete-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
  padding: 4px 8px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 4px;
  color: #ef4444;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.task-delete-btn:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: #ef4444;
  transform: translateY(-1px);
}

.task-delete-btn svg {
  flex-shrink: 0;
}

.task-instant {
  color: #22c55e;
  font-weight: 600;
}

.task-done-hint {
  color: #22c55e;
  font-weight: 600;
}

.task-hashing-hint {
  color: #f59e0b;
  font-size: 12px;
}

.task-merging-hint {
  color: #a855f7;
  font-size: 12px;
}

.task-uploaded {
  color: rgba(255, 255, 255, 0.45);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}

.task-speed {
  color: #60a5fa;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.task-time {
  color: #38bdf8;
  font-variant-numeric: tabular-nums;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.task-remaining {
  color: rgba(56, 189, 248, 0.75);
  font-size: 11px;
}

.task-queued-hint {
  color: #94a3b8;
  font-style: italic;
}

/* ===== 文件列表 ===== */
.section-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.search-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  transition: all 0.2s ease-out;
}

.search-wrap:focus-within {
  border-color: rgba(37, 99, 235, 0.5);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.search-wrap svg {
  color: rgba(255, 255, 255, 0.25);
  flex-shrink: 0;
}

.search-wrap input {
  background: none; border: none; outline: none;
  color: #f1f5f9; font-size: 14px; width: 200px;
  font-family: 'Open Sans', sans-serif;
}

.search-wrap input::placeholder {
  color: rgba(255, 255, 255, 0.2);
}

.file-table-wrap {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 16px;
  overflow: hidden;
}

.file-table {
  width: 100%;
  border-collapse: collapse;
}

.file-table th {
  padding: 14px 18px;
  text-align: left;
  font-size: 11px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.35);
  text-transform: uppercase;
  letter-spacing: 1px;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.file-row {
  transition: background 0.15s ease-out;
}

.file-row:hover {
  background: rgba(37, 99, 235, 0.06);
}

.file-row td {
  padding: 14px 18px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-ext-badge {
  width: 40px; height: 40px;
  border-radius: 10px;
  background: rgba(37, 99, 235, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 700;
  color: #60a5fa;
  letter-spacing: 0.5px;
  flex-shrink: 0;
}

.file-name-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.85);
}

.text-muted {
  color: rgba(255, 255, 255, 0.4) !important;
  font-size: 13px !important;
}

.text-center {
  text-align: center;
}

.action-btns {
  display: flex;
  gap: 6px;
}

.action-btn {
  width: 34px; height: 34px;
  border: none;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease-out;
}

.action-btn.download {
  background: rgba(37, 99, 235, 0.1);
  color: #60a5fa;
}

.action-btn.download:hover {
  background: rgba(37, 99, 235, 0.25);
  transform: translateY(-1px);
}

.action-btn.delete {
  background: rgba(239, 68, 68, 0.08);
  color: rgba(239, 68, 68, 0.6);
}

.action-btn.delete:hover {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  transform: translateY(-1px);
}

.action-btn.speed-btn {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.15), rgba(234, 88, 12, 0.25));
  color: #f59e0b;
  border: 1px solid rgba(245, 158, 11, 0.35);
}

.action-btn.speed-btn:hover {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.3), rgba(234, 88, 12, 0.4));
  color: #fbbf24;
  box-shadow: 0 0 12px rgba(245, 158, 11, 0.45);
  transform: translateY(-1px);
}

.action-btn.cli-btn {
  background: rgba(100, 116, 139, 0.15);
  color: #94a3b8;
  border: 1px solid rgba(148, 163, 184, 0.25);
}

.action-btn.cli-btn:hover {
  background: rgba(56, 189, 248, 0.18);
  color: #38bdf8;
  border-color: rgba(56, 189, 248, 0.5);
  box-shadow: 0 0 10px rgba(56, 189, 248, 0.35);
  transform: translateY(-1px);
}

/* 空状态 */
.empty-state {
  padding: 64px 24px;
  text-align: center;
  color: rgba(255, 255, 255, 0.2);
}

.empty-state p {
  font-size: 16px;
  margin: 12px 0 4px;
  color: rgba(255, 255, 255, 0.35);
}

.empty-state span {
  font-size: 13px;
}

/* 分页 */
.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
}

.page-info {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.35);
}

:deep(.el-pagination) {
  --el-pagination-bg-color: transparent;
  --el-pagination-text-color: rgba(255,255,255,0.5);
  --el-pagination-button-bg-color: rgba(255,255,255,0.04);
}

/* ===== 新版子系统卡片与 Modern 表格 ===== */
.user-section-card {
  background: rgba(10, 15, 30, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 20px;
  padding: 24px;
  backdrop-filter: blur(20px);
  box-shadow: 0 10px 40px -10px rgba(0,0,0,0.5);
}

.title-icon-wrap {
  width: 32px; height: 32px;
  background: linear-gradient(135deg, rgba(0, 240, 255, 0.2), rgba(124, 58, 237, 0.2));
  color: #00f0ff;
  border-radius: 10px;
  display: flex; align-items: center; justify-content: center;
  border: 1px solid rgba(0, 240, 255, 0.2);
}

.table-card {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  overflow: hidden;
}

.modern-table {
  width: 100%;
  border-collapse: collapse;
}

.modern-table th {
  text-align: left;
  padding: 18px 20px;
  font-size: 12px;
  font-family: 'Outfit', sans-serif;
  color: #94a3b8;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.table-row {
  transition: all 0.2s ease;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}

.table-row:hover {
  background: rgba(255, 255, 255, 0.04);
}

.modern-table td {
  padding: 16px 20px;
  font-size: 14px;
  color: #f1f5f9;
}

.id-badge {
  background: rgba(255, 255, 255, 0.05);
  padding: 4px 8px;
  border-radius: 6px;
  font-family: monospace;
  font-size: 13px;
  color: #64748b;
  border: 1px solid rgba(255,255,255,0.1);
}

.user-account-cell {
  display: flex; align-items: center; gap: 12px;
}
.user-avatar-mini {
  width: 32px; height: 32px;
  border-radius: 8px;
  background: linear-gradient(135deg, #00f0ff, #7c3aed);
  display: flex; align-items: center; justify-content: center;
  font-weight: 700; color: #fff; font-size: 14px;
  box-shadow: 0 2px 10px rgba(0, 240, 255, 0.3);
}
.user-name-text {
  font-weight: 600; font-size: 14px;
}

.space-capsule {
  display: flex; flex-direction: column; gap: 8px;
  max-width: 200px;
  position: relative;
}
.space-text {
  display: flex; align-items: baseline; gap: 4px; font-size: 13px;
}
.space-text .used { font-weight: 700; color: #f8fafc; }
.space-text .divider { color: #475569; margin: 0 4px; }
.space-text .total { color: #64748b; font-size: 12px; }

.space-bar-bg {
  width: 100%; height: 6px; background: rgba(0,0,0,0.5);
  border-radius: 3px; overflow: hidden;
  box-shadow: inset 0 1px 2px rgba(0,0,0,0.5);
}
.space-bar-fill {
  height: 100%; background: linear-gradient(90deg, #00f0ff, #3b82f6);
  border-radius: 3px; transition: width 0.4s ease;
}
.space-bar-fill.danger-fill {
  background: linear-gradient(90deg, #f59e0b, #ef4444);
}

.reset-space-btn {
  position: absolute; right: 0; top: 0;
  background: transparent; border: none; color: #00f0ff;
  cursor: pointer; opacity: 0; transition: all 0.2s;
  padding: 2px;
}
.space-capsule:hover .reset-space-btn {
  opacity: 1; transform: translateY(-2px);
}
.reset-space-btn:hover { color: #fff; }

.status-badge {
  display: inline-block; padding: 6px 10px; border-radius: 20px;
  font-size: 12px; font-weight: 600; letter-spacing: 0.5px;
}
.status-badge.success { background: rgba(16, 185, 129, 0.15); color: #34d399; border: 1px solid rgba(16, 185, 129, 0.3); }
.status-badge.error { background: rgba(239, 68, 68, 0.15); color: #f87171; border: 1px solid rgba(239, 68, 68, 0.3); }

.action-btn-modern {
  background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1);
  color: #cbd5e1; border-radius: 10px; width: 34px; height: 34px;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; transition: all 0.2s;
}
.action-btn-modern.freeze-btn:hover { background: rgba(239, 68, 68, 0.15); color: #ef4444; border-color: rgba(239, 68, 68, 0.3); }
.action-btn-modern.unfreeze-btn:hover { background: rgba(16, 185, 129, 0.15); color: #10b981; border-color: rgba(16, 185, 129, 0.3); }
.action-btn-modern.space-btn:hover { background: rgba(147, 51, 234, 0.2); color: #c084fc; border-color: rgba(192, 132, 252, 0.4); transform: translateY(-1px); }

.gift-space-btn {
  margin-top: 8px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 5px 10px;
  background: rgba(37, 99, 235, 0.15);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 6px;
  color: #60a5fa;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.gift-space-btn:hover {
  background: rgba(37, 99, 235, 0.3);
  border-color: rgba(59, 130, 246, 0.5);
  color: #93c5fd;
}

/* 🚀 极速多线程并发下载看板（文叔叔式科技极客风） */
.speed-download-panel {
  position: fixed;
  right: 28px;
  bottom: 28px;
  width: 380px;
  background: rgba(15, 23, 42, 0.9);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.35);
  border-radius: 16px;
  box-shadow: 0 20px 45px rgba(0, 0, 0, 0.6), 0 0 20px rgba(56, 189, 248, 0.15);
  padding: 16px 18px;
  z-index: 9999;
  font-family: inherit;
  color: #f1f5f9;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.sdp-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.sdp-title-box {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow: hidden;
}

.sdp-badge {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.25), rgba(234, 88, 12, 0.25));
  border: 1px solid rgba(245, 158, 11, 0.5);
  color: #fbbf24;
  white-space: nowrap;
}

.sdp-filename {
  font-size: 13px;
  font-weight: 600;
  color: #e2e8f0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 200px;
}

.sdp-close-btn {
  background: transparent;
  border: none;
  color: #64748b;
  font-size: 14px;
  cursor: pointer;
  padding: 4px;
  line-height: 1;
  transition: color 0.2s;
}

.sdp-close-btn:hover {
  color: #f43f5e;
}

.sdp-stats {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 8px;
}

.sdp-stat-main {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.sdp-speed-value {
  font-size: 22px;
  font-weight: 800;
  color: #38bdf8;
  letter-spacing: -0.5px;
  text-shadow: 0 0 12px rgba(56, 189, 248, 0.4);
}

.sdp-status-tag {
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 4px;
  background: rgba(56, 189, 248, 0.15);
  color: #7dd3fc;
}

.sdp-status-tag.completed {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
}

.sdp-status-tag.error {
  background: rgba(239, 68, 68, 0.2);
  color: #f87171;
}

.sdp-stat-sub {
  font-size: 12px;
  color: #94a3b8;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
}

.sdp-percent {
  font-weight: 700;
  color: #f8fafc;
}

.sdp-progress-track {
  width: 100%;
  height: 6px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 12px;
  position: relative;
}

.sdp-progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #0ea5e9, #6366f1, #a855f7);
  background-size: 200% 100%;
  border-radius: 3px;
  transition: width 0.2s ease-out;
  animation: sdpGradientMove 2s linear infinite;
  box-shadow: 0 0 8px rgba(99, 102, 241, 0.6);
}

@keyframes sdpGradientMove {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* 6条并发管道通道进度展示 */
.sdp-channels {
  display: flex;
  gap: 4px;
  padding-top: 4px;
  border-top: 1px dashed rgba(255, 255, 255, 0.08);
}

.sdp-channel-item {
  flex: 1;
  height: 4px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 2px;
  overflow: hidden;
  position: relative;
}

.sdp-ch-fill {
  height: 100%;
  background: #38bdf8;
  border-radius: 2px;
  transition: width 0.15s ease-out;
  box-shadow: 0 0 4px rgba(56, 189, 248, 0.8);
}

/* ===== 无障碍 ===== */
@media (prefers-reduced-motion: reduce) {
  .upload-zone:hover { transform: none; }
  .upload-icon-ring.pulse { animation: none; }
  .action-btn:hover, .feat-card:hover, .pm-submit-btn:hover { transform: none; }
  .sdp-progress-bar { animation: none; }
}
</style>
